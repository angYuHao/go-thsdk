// Package httpserver exposes a THSDK client over HTTP so the SDK can also run as
// a standalone service. It is a plain net/http handler, so it can be mounted in
// an existing application instead of running cmd/thsdk-server.
package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	thsdk "github.com/angYuHao/go-thsdk"
)

// Config configures a Server.
type Config struct {
	// Client is the authenticated SDK client the server forwards to.
	Client *thsdk.Client
	// AllowWrites enables the account-mutation methods. Off by default, because
	// they change the account's cloud watchlist without an undo.
	AllowWrites bool
	// CallTimeout bounds each upstream call when the request carries no deadline.
	// Zero means 30 seconds.
	CallTimeout time.Duration
	// MaxBodyBytes limits a request body. Zero means 1 MiB, the ABI's own limit.
	MaxBodyBytes int64
	// Logger receives one line per request. Nil disables logging.
	Logger func(format string, args ...any)
}

// Server adapts a THSDK client to HTTP.
type Server struct {
	config   Config
	catalog  map[string]methodInfo
	handler  http.Handler
	maxBytes int64
	timeout  time.Duration
}

// methodEntry is one row of the /v1/methods listing.
type methodEntry struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Write    bool   `json:"write"`
}

// New validates the configuration and builds the route table.
func New(config Config) (*Server, error) {
	if config.Client == nil {
		return nil, errors.New("httpserver: Client 不能为空")
	}
	timeout := config.CallTimeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	if timeout < 0 || timeout > 10*time.Minute {
		return nil, errors.New("httpserver: CallTimeout 必须在 0 到 10 分钟之间")
	}
	maxBytes := config.MaxBodyBytes
	if maxBytes == 0 {
		maxBytes = 1 << 20
	}
	if maxBytes < 0 {
		return nil, errors.New("httpserver: MaxBodyBytes 不能为负")
	}

	server := &Server{
		config:   config,
		catalog:  methodCatalog(),
		maxBytes: maxBytes,
		timeout:  timeout,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealth)
	mux.HandleFunc("/v1/methods", server.handleMethods)
	mux.HandleFunc("/v1/call/", server.handleCall)
	server.handler = mux
	return server, nil
}

// Handler returns the http.Handler, ready for ListenAndServe or mounting.
func (s *Server) Handler() http.Handler { return s.handler }

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "只支持 GET")
		return
	}
	enabled := 0
	for _, info := range s.catalog {
		if info.Write && !s.config.AllowWrites {
			continue
		}
		enabled++
	}
	s.writeJSON(w, http.StatusOK, map[string]any{
		"status":        "ok",
		"authenticated": s.config.Client.Authenticated(),
		"methods":       enabled,
	})
}

func (s *Server) handleMethods(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "只支持 GET")
		return
	}
	names := make([]string, 0, len(s.catalog))
	for name := range s.catalog {
		names = append(names, name)
	}
	sort.Strings(names)

	entries := make([]methodEntry, 0, len(names))
	for _, name := range names {
		info := s.catalog[name]
		if info.Write && !s.config.AllowWrites {
			continue
		}
		entries = append(entries, methodEntry{Name: name, Category: info.Category, Write: info.Write})
	}
	s.writeJSON(w, http.StatusOK, map[string]any{
		"count":   len(entries),
		"writes":  s.config.AllowWrites,
		"methods": entries,
	})
}

func (s *Server) handleCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "只支持 POST")
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/v1/call/")
	name = strings.Trim(name, "/")
	if name == "" || strings.Contains(name, "/") {
		s.writeError(w, http.StatusBadRequest, "invalid_method", "路径应为 /v1/call/<method>")
		return
	}
	info, known := s.catalog[name]
	if !known {
		s.writeError(w, http.StatusNotFound, "unknown_method", fmt.Sprintf("未暴露的方法 %q", name))
		return
	}
	if info.Write && !s.config.AllowWrites {
		s.writeError(w, http.StatusForbidden, "writes_disabled",
			"写方法未启用，请用 -allow-writes 显式开启")
		return
	}

	params, err := s.decodeParams(w, r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid_body", err.Error())
		return
	}

	ctx, cancel := s.requestContext(r)
	defer cancel()

	started := time.Now()
	data, callErr := s.config.Client.CallRaw(ctx, name, params)
	elapsed := time.Since(started)
	if callErr != nil {
		status, code := classify(callErr)
		s.log("%s %s -> %d %s (%s)", r.Method, r.URL.Path, status, code, elapsed.Round(time.Millisecond))
		s.writeError(w, status, code, callErr.Error())
		return
	}
	s.log("%s %s -> 200 %d bytes (%s)", r.Method, r.URL.Path, len(data.Bytes()), elapsed.Round(time.Millisecond))

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if data.IsEmpty() {
		_, _ = w.Write([]byte("null"))
		return
	}
	_, _ = w.Write(data.Bytes())
}

// decodeParams reads the body as the method's params object. An empty body means
// "no params".
func (s *Server) decodeParams(w http.ResponseWriter, r *http.Request) (any, error) {
	reader := http.MaxBytesReader(w, r.Body, s.maxBytes)
	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("请求体读取失败或超过大小限制")
	}
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" {
		return nil, nil
	}
	if !json.Valid(raw) {
		return nil, errors.New("请求体不是合法 JSON")
	}
	// Hand the document through untouched so the ABI's own field-name and depth
	// validation reports the precise reason.
	return json.RawMessage(raw), nil
}

func (s *Server) requestContext(r *http.Request) (context.Context, context.CancelFunc) {
	if deadline, ok := r.Context().Deadline(); ok {
		return context.WithDeadline(r.Context(), deadline)
	}
	return context.WithTimeout(r.Context(), s.timeout)
}

// classify maps an SDK error to an HTTP status and a stable error code.
func classify(err error) (int, string) {
	if errors.Is(err, thsdk.ErrNotAuthenticated) {
		return http.StatusServiceUnavailable, "not_authenticated"
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout, "timeout"
	}
	var apiErr *thsdk.Error
	if errors.As(err, &apiErr) {
		switch apiErr.Code {
		case "permission_denied", "no_permission":
			return http.StatusForbidden, apiErr.Code
		case "no_data":
			return http.StatusNotFound, apiErr.Code
		case "rate_limited":
			return http.StatusTooManyRequests, apiErr.Code
		case "timeout", "query_failed":
			return http.StatusGatewayTimeout, apiErr.Code
		case "invalid_request":
			return http.StatusBadRequest, apiErr.Code
		case "session_unavailable":
			return http.StatusServiceUnavailable, apiErr.Code
		default:
			return http.StatusBadGateway, apiErr.Code
		}
	}
	var paramErr *thsdk.ParameterError
	if errors.As(err, &paramErr) {
		return http.StatusBadRequest, "parameter_error"
	}
	var authErr *thsdk.AuthError
	if errors.As(err, &authErr) {
		return http.StatusUnauthorized, "auth_error"
	}
	return http.StatusBadGateway, "upstream_error"
}

func (s *Server) writeError(w http.ResponseWriter, status int, code, message string) {
	s.writeJSON(w, status, map[string]any{
		"error": map[string]any{"code": code, "message": message},
	})
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(payload)
}

func (s *Server) log(format string, args ...any) {
	if s.config.Logger != nil {
		s.config.Logger(format, args...)
	}
}
