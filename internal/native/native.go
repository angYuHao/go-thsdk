// Package native binds the THSDK runtime component (hq.dll / hq.so / hq.dylib).
//
// The component exposes a JSON call ABI:
//
//	void* CallAlloc(const char* request); // JSON in, JSON out, caller must free
//	void  FreeResult(void* result);
//
// Request envelope: {"method":name,"params":...,"timeout_ms":n}
// Response envelope: {"ok":bool,"data":...,"error":{code,message,retry_after_ms}}
package native

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ABIVersion is the component ABI this SDK speaks.
const ABIVersion = "2.0.0"

const (
	maxRequestBytes  = 1 << 20
	maxMethodBytes   = 128
	maxJSONDepth     = 64
	MaxCallTimeoutMS = 600_000

	// defaultMinSpacing matches the component's own throttle message:
	// "请保持至少 50ms 的调用间隔".
	defaultMinSpacing = 50 * time.Millisecond
)

var fieldNamePattern = regexp.MustCompile(`^[a-z0-9_]+$`)

// APIError is a business error reported by the runtime component.
type APIError struct {
	Code         string
	Message      string
	retryAfterMS int
	hasRetryTime bool
}

func (e *APIError) Error() string {
	if e.Code == "" {
		return e.Message
	}
	return e.Code + ": " + e.Message
}

// RetryAfterMS reports the runtime's rate-limit hint, and whether it was present.
func (e *APIError) RetryAfterMS() (int, bool) { return e.retryAfterMS, e.hasRetryTime }

func newAPIError(raw json.RawMessage) *APIError {
	err := &APIError{Code: "unknown_error", Message: "unknown error"}
	if len(raw) == 0 || string(raw) == "null" {
		return err
	}
	var fields struct {
		Code         *string         `json:"code"`
		Message      *string         `json:"message"`
		RetryAfterMS json.RawMessage `json:"retry_after_ms"`
	}
	if jsonErr := json.Unmarshal(raw, &fields); jsonErr != nil {
		err.Message = strings.TrimSpace(string(raw))
		return err
	}
	if fields.Code != nil && *fields.Code != "" {
		err.Code = *fields.Code
	}
	if fields.Message != nil && *fields.Message != "" {
		err.Message = *fields.Message
	}
	if len(fields.RetryAfterMS) > 0 {
		var candidate int
		if json.Unmarshal(fields.RetryAfterMS, &candidate) == nil && candidate > 0 {
			err.retryAfterMS = candidate
			err.hasRetryTime = true
		}
	}
	return err
}

// IsRateLimited reports whether err is the component's throttling response.
func IsRateLimited(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.Code == "rate_limited"
}

type responseEnvelope struct {
	OK    *bool           `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error json.RawMessage `json:"error"`
}

// Runtime is a handle on the loaded component. The component throttles calls
// below a 50 ms interval, so calls are serialized and spaced here.
type Runtime struct {
	callMu     sync.Mutex
	lib        *mappedLibrary
	callAlloc  *funcPtr
	freeResult *funcPtr
	MinSpacing time.Duration
}

// Open loads the runtime component at path and completes the ABI handshake.
func Open(path string) (*Runtime, error) {
	lib, err := retainLibrary(path)
	if err != nil {
		return nil, err
	}
	rt := &Runtime{lib: lib, MinSpacing: defaultMinSpacing}
	if rt.callAlloc, err = lib.lookup("CallAlloc"); err != nil {
		lib.release()
		return nil, err
	}
	if rt.freeResult, err = lib.lookup("FreeResult"); err != nil {
		lib.release()
		return nil, err
	}
	if err := rt.handshake(); err != nil {
		lib.release()
		return nil, err
	}
	return rt, nil
}

// Close releases this handle. The component itself stays mapped for the lifetime
// of the process because a Go c-shared image cannot be reloaded after unloading.
func (r *Runtime) Close() error {
	r.callMu.Lock()
	defer r.callMu.Unlock()
	if r == nil || r.lib == nil {
		return nil
	}
	err := r.lib.release()
	r.lib = nil
	return err
}

func (r *Runtime) handshake() error {
	data, err := r.handshakeCall("version", nil)
	if err != nil {
		return err
	}
	var payload struct {
		ABIVersion string `json:"abi_version"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("thsdk: 无法解析运行组件版本: %w", err)
	}
	if payload.ABIVersion != ABIVersion {
		return fmt.Errorf("thsdk: 运行组件 ABI 版本 %q 与本 SDK 要求的 %q 不兼容", payload.ABIVersion, ABIVersion)
	}

	// ABI 2.0.0 predates strict request decoding and kept the same version string,
	// so probe the behavior instead of trusting it.
	if _, probeErr := r.handshakeCall("version", map[string]any{"thsdk_protocol_probe": true}); probeErr == nil {
		return errors.New("thsdk: 运行组件请求协议过旧，请安装最新版运行组件")
	} else {
		var strict *APIError
		if !errors.As(probeErr, &strict) || strict.Code != "invalid_request" {
			return fmt.Errorf("thsdk: 运行组件请求协议探测失败: %w", probeErr)
		}
	}
	return nil
}

// handshakeCall keeps startup probes readable; spacing is enforced per call.
func (r *Runtime) handshakeCall(method string, params any) (json.RawMessage, error) {
	return r.Invoke(method, params, 0)
}

// Invoke performs one JSON call and returns the untouched `data` payload.
// A timeoutMS of zero lets the component apply its own default.
func (r *Runtime) Invoke(method string, params any, timeoutMS int) (json.RawMessage, error) {
	encoded, err := encodeRequest(method, params, timeoutMS)
	if err != nil {
		return nil, err
	}

	r.callMu.Lock()
	defer r.callMu.Unlock()
	if r.lib == nil {
		return nil, errors.New("thsdk: 运行组件已关闭")
	}

	payload, err := r.callLocked(encoded)
	if err != nil {
		return nil, err
	}

	var envelope responseEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return nil, errors.New("thsdk: 运行组件返回了无效数据")
	}
	if envelope.OK == nil {
		return nil, errors.New("thsdk: 运行组件返回了不兼容的数据")
	}
	if !*envelope.OK {
		return nil, newAPIError(envelope.Error)
	}
	return envelope.Data, nil
}

func (r *Runtime) callLocked(encoded []byte) ([]byte, error) {
	r.lib.pace(r.MinSpacing)

	// The component reads a NUL-terminated C string.
	request := make([]byte, len(encoded)+1)
	copy(request, encoded)

	raw, err := r.callAlloc.call(request)
	if err != nil {
		return nil, fmt.Errorf("thsdk: 调用运行组件失败: %w", err)
	}
	if raw == nil {
		return nil, errors.New("thsdk: 运行组件未返回结果")
	}
	payload, readErr := r.lib.readResult(raw)
	freeErr := r.freeResult.free(raw)
	if readErr != nil {
		return nil, readErr
	}
	if freeErr != nil {
		return nil, freeErr
	}
	return payload, nil
}

func encodeRequest(method string, params any, timeoutMS int) ([]byte, error) {
	if err := validateMethod(method); err != nil {
		return nil, err
	}
	if timeoutMS != 0 && (timeoutMS < 0 || timeoutMS > MaxCallTimeoutMS) {
		return nil, fmt.Errorf("thsdk: timeout_ms 必须是 0（默认）或 1 到 %d 之间的整数", MaxCallTimeoutMS)
	}

	// Marshal first so any Go value with a json tag is accepted, then validate
	// the document that actually goes over the ABI.
	request := struct {
		Method    string `json:"method"`
		Params    any    `json:"params,omitempty"`
		TimeoutMS int    `json:"timeout_ms,omitempty"`
	}{Method: method, Params: params, TimeoutMS: timeoutMS}

	encoded, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("thsdk: 请求无法编码为严格 JSON: %w", err)
	}
	if len(encoded) > maxRequestBytes {
		return nil, errors.New("thsdk: 请求超过 1 MiB 限制")
	}
	if err := validateJSONDocument(encoded); err != nil {
		return nil, err
	}
	return encoded, nil
}

func validateMethod(method string) error {
	if method == "" || method != strings.TrimSpace(method) {
		return errors.New("thsdk: method 必须是非空规范名称")
	}
	if len(method) > maxMethodBytes || !fieldNamePattern.MatchString(method) {
		return errors.New("thsdk: method 必须是不超过 128 字节的 ASCII 小写 snake_case")
	}
	return nil
}

// validateJSONDocument enforces the ABI's field-name and nesting rules on the
// encoded request. json.Marshal already rejects NaN and Infinity.
func validateJSONDocument(encoded []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()
	fieldStack := make([]string, 0, maxJSONDepth)
	return walkValue(decoder, &fieldStack, 0)
}

func walkValue(decoder *json.Decoder, fields *[]string, depth int) error {
	if depth > maxJSONDepth {
		return errors.New("thsdk: 请求的 JSON 嵌套超过 64 层")
	}
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	switch delimiter := token.(type) {
	case json.Delim:
		switch delimiter {
		case '{':
			for decoder.More() {
				key, err := decoder.Token()
				if err != nil {
					return err
				}
				name, ok := key.(string)
				if !ok {
					return errors.New("thsdk: 请求字段必须是字符串")
				}
				if !allowedFieldName(name) {
					return errors.New("thsdk: 请求字段必须使用 ASCII 小写 snake_case")
				}
				*fields = append(*fields, name)
				if err := walkValue(decoder, fields, depth+1); err != nil {
					return err
				}
				*fields = (*fields)[:len(*fields)-1]
			}
			_, err := decoder.Token() // closing brace
			return err
		case '[':
			for decoder.More() {
				if err := walkValue(decoder, fields, depth+1); err != nil {
					return err
				}
			}
			_, err := decoder.Token() // closing bracket
			return err
		default:
			return fmt.Errorf("thsdk: 请求包含非法 JSON 定界符 %v", delimiter)
		}
	case string, float64, json.Number, bool, nil:
		return nil
	default:
		return fmt.Errorf("thsdk: 请求包含不可序列化的值类型: %T", token)
	}
}

// allowedFieldName accepts the envelope's own camelCase keys and snake_case
// everywhere below them.
func allowedFieldName(name string) bool {
	switch name {
	case "method", "params", "timeout_ms":
		return true
	}
	return fieldNamePattern.MatchString(name)
}
