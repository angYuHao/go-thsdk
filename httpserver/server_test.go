package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	thsdk "github.com/angYuHao/go-thsdk"
)

// The server needs a *thsdk.Client, which requires the native component. These
// tests cover the routing and validation paths that do not reach the component,
// so they run without hq.dll.

func TestNewValidatesClient(t *testing.T) {
	if _, err := New(Config{}); err == nil {
		t.Fatal("expected an error for a nil client")
	}
}

func TestCatalogCoversReadsAndWrites(t *testing.T) {
	catalog := methodCatalog()
	if len(catalog) != 80 {
		t.Fatalf("catalog size = %d, want 80", len(catalog))
	}
	writes := 0
	for _, info := range catalog {
		if info.Write {
			writes++
		}
	}
	if writes != len(writeCatalog()) {
		t.Fatalf("write entries = %d, want %d", writes, len(writeCatalog()))
	}
}

func TestClassifyMapsCodes(t *testing.T) {
	cases := []struct {
		code   string
		status int
	}{
		{"permission_denied", http.StatusForbidden},
		{"no_data", http.StatusNotFound},
		{"rate_limited", http.StatusTooManyRequests},
		{"timeout", http.StatusGatewayTimeout},
		{"invalid_request", http.StatusBadRequest},
		{"session_unavailable", http.StatusServiceUnavailable},
		{"encode_failed", http.StatusBadGateway},
	}
	for _, item := range cases {
		status, code := classify(&thsdk.Error{Code: item.code})
		if status != item.status || code != item.code {
			t.Errorf("%s -> (%d, %s), want (%d, %s)", item.code, status, code, item.status, item.code)
		}
	}
}

func TestClassifyMapsLocalErrors(t *testing.T) {
	if status, code := classify(thsdk.ErrNotAuthenticated); status != http.StatusServiceUnavailable || code != "not_authenticated" {
		t.Errorf("ErrNotAuthenticated -> (%d, %s)", status, code)
	}
	if status, code := classify(context.DeadlineExceeded); status != http.StatusGatewayTimeout || code != "timeout" {
		t.Errorf("DeadlineExceeded -> (%d, %s)", status, code)
	}
}

func TestMethodsRouteShape(t *testing.T) {
	server := &Server{
		config:  Config{},
		catalog: methodCatalog(),
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/methods", nil)
	server.handleMethods(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	var payload struct {
		Count   int           `json:"count"`
		Writes  bool          `json:"writes"`
		Methods []methodEntry `json:"methods"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Writes {
		t.Error("writes should be off by default")
	}
	if payload.Count != 70 {
		t.Errorf("count = %d, want 70 read-only methods", payload.Count)
	}
	for _, entry := range payload.Methods {
		if entry.Write {
			t.Errorf("write method %q leaked into the default listing", entry.Name)
		}
	}
}

func TestCallRejectsUnknownAndWriteMethods(t *testing.T) {
	server := &Server{
		config:   Config{},
		catalog:  methodCatalog(),
		maxBytes: 1 << 20,
		timeout:  5,
	}

	for _, path := range []string{"/v1/call/no_such_method", "/v1/call/clear_account_watchlist"} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, path, strings.NewReader("{}"))
		server.handleCall(recorder, request)
		if recorder.Code != http.StatusNotFound && recorder.Code != http.StatusForbidden {
			t.Errorf("%s -> %d, want 404 or 403", path, recorder.Code)
		}
		if !strings.Contains(recorder.Body.String(), `"error"`) {
			t.Errorf("%s body lacks an error envelope: %s", path, recorder.Body.String())
		}
	}
}

func TestCallRejectsNonPost(t *testing.T) {
	server := &Server{config: Config{}, catalog: methodCatalog()}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/call/list_security_k_lines", nil)
	server.handleCall(recorder, request)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", recorder.Code)
	}
}

func TestDecodeParamsRejectsBadJSON(t *testing.T) {
	server := &Server{maxBytes: 1 << 20}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/call/query_wencai", strings.NewReader("{not json"))
	if _, err := server.decodeParams(recorder, request); err == nil {
		t.Fatal("expected invalid JSON to be rejected")
	}
}

func TestDecodeParamsAcceptsEmptyBody(t *testing.T) {
	server := &Server{maxBytes: 1 << 20}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/call/get_account_watchlist", nil)
	params, err := server.decodeParams(recorder, request)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if params != nil {
		t.Fatalf("params = %#v, want nil", params)
	}
}
