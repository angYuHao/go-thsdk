package native_test

import (
	"testing"

	"github.com/angYuHao/go-thsdk/internal/native"
)

func TestRepeatedLoadUnloadCycles(t *testing.T) {
	path := runtimePath(t)
	for i := 0; i < 3; i++ {
		rt, err := native.Open(path)
		if err != nil {
			t.Fatalf("open %d: %v", i, err)
		}
		data, err := rt.Invoke("version", nil, 0)
		if err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
		t.Logf("cycle %d: %s", i, data)
		if err := rt.Close(); err != nil {
			t.Fatalf("close %d: %v", i, err)
		}
	}
}
