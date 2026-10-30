//go:build linux || darwin

package native

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

// unsafePointer is the neutral spelling of a raw address handed over the ABI.
type unsafePointer = unsafe.Pointer

func mapLibrary(path string) (platformLibrary, error) {
	cpath := C.CString(path)
	defer C.free(cpath)
	handle := C.dlopen(cpath, C.RTLD_NOW)
	if handle == nil {
		return nil, fmt.Errorf("thsdk: 运行组件加载失败，请确认安装包与当前环境匹配: %s", C.GoString(C.dlerror()))
	}
	return &unixLibrary{handle: unsafe.Pointer(handle)}, nil
}

// funcPtr holds an exported function's address resolved through dlsym.
type funcPtr struct{ addr unsafePointer }

func (p *funcPtr) call(request []byte) (unsafePointer, error) {
	type callFn func(unsafePointer) unsafePointer
	fn := *(*callFn)(unsafePointer(p))
	var arg unsafePointer
	if len(request) > 0 {
		arg = unsafePointer(&request[0])
	}
	return fn(arg), nil
}

func (p *funcPtr) free(ptr unsafePointer) error {
	if ptr == nil {
		return nil
	}
	type freeFn func(unsafePointer)
	fn := *(*freeFn)(unsafePointer(p))
	fn(ptr)
	return nil
}

type unixLibrary struct{ handle unsafePointer }

func (l *unixLibrary) lookup(name string) (*funcPtr, error) {
	cname := C.CString(name)
	defer C.free(cname)
	C.dlerror()
	symbol := C.dlsym(unsafe.Pointer(l.handle), cname)
	if symbol == nil {
		return nil, fmt.Errorf("thsdk: 运行组件缺少导出 %s: %s", name, C.GoString(C.dlerror()))
	}
	return &funcPtr{addr: symbol}, nil
}

func (l *unixLibrary) readResult(ptr unsafePointer) ([]byte, error) {
	if ptr == nil {
		return nil, fmt.Errorf("thsdk: 运行组件返回了空指针")
	}
	return []byte(C.GoString((*C.char)(ptr))), nil
}

// unmap is deliberately empty; see mappedLibrary.release.
func (l *unixLibrary) unmap() error { return nil }

func sleepDuration(d time.Duration) { time.Sleep(d) }
