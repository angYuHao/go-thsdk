package native_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/angYuHao/go-thsdk/internal/native"
)

func runtimePath(t *testing.T) string {
	t.Helper()
	path := os.Getenv("THSDK_RUNTIME_PATH")
	if path == "" {
		name := map[string]string{"windows": "hq.dll", "darwin": "hq.dylib", "linux": "hq.so"}[runtime.GOOS]
		path = filepath.Join("..", "..", "libs", runtime.GOOS, runtime.GOARCH, name)
	}
	if _, err := os.Stat(path); err != nil {
		t.Skipf("runtime component not available: %v", err)
	}
	return path
}

func TestOpenAndVersion(t *testing.T) {
	rt, err := native.Open(runtimePath(t))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer rt.Close()

	data, err := rt.Invoke("version", nil, 0)
	if err != nil {
		t.Fatalf("version: %v", err)
	}
	var pretty map[string]any
	if err := json.Unmarshal(data, &pretty); err != nil {
		t.Fatalf("version payload: %v (%s)", err, data)
	}
	t.Logf("version data: %s", data)

	if _, err := rt.Invoke("methods", nil, 0); err != nil {
		t.Fatalf("methods: %v", err)
	}
}

func TestClientRejectsCamelCaseField(t *testing.T) {
	rt, err := native.Open(runtimePath(t))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer rt.Close()

	// The SDK validates before the ABI is touched, so no request is billed.
	if _, err := rt.Invoke("version", map[string]any{"BadKey": true}, 0); err == nil {
		t.Fatalf("expected a camelCase field to be rejected")
	}
}

func TestStrictRequestValidation(t *testing.T) {
	rt, err := native.Open(runtimePath(t))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer rt.Close()

	_, err = rt.Invoke("version", map[string]any{"local_only_field": true}, 0)
	var apiErr *native.APIError
	if err == nil {
		t.Fatalf("expected the component to reject an unknown request field")
	}
	if !asAPI(err, &apiErr) || apiErr.Code != "invalid_request" {
		t.Fatalf("expected invalid_request, got %v", err)
	}
}

func asAPI(err error, target **native.APIError) bool {
	apiErr, ok := err.(*native.APIError)
	if ok {
		*target = apiErr
	}
	return ok
}
