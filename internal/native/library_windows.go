//go:build windows

package native

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

// unsafePointer is the neutral spelling of a raw address handed over the ABI.
type unsafePointer = unsafe.Pointer

// funcPtr holds an exported function's address. It is a real Go object, so it
// stays valid for as long as the Runtime references it.
type funcPtr struct{ addr uintptr }

func (p *funcPtr) call(request []byte) (unsafePointer, error) {
	var arg unsafePointer
	if len(request) > 0 {
		arg = unsafePointer(&request[0])
	}
	ret, _, err := syscall.Syscall(p.addr, 1, uintptr(arg), 0, 0)
	if err != syscall.Errno(0) {
		return nil, err
	}
	return buffer(ret), nil
}

func (p *funcPtr) free(ptr unsafePointer) error {
	_, _, err := syscall.Syscall(p.addr, 1, uintptr(ptr), 0, 0)
	if err != syscall.Errno(0) {
		return fmt.Errorf("thsdk: 释放结果缓冲区失败: %w", err)
	}
	return nil
}

// buffer adopts a raw address returned by the component. The component, not the
// Go runtime, owns the result buffer, so the uintptr round-trip is intentional.
func buffer(addr uintptr) unsafePointer { return unsafePointer(addr) }

type windowsLibrary struct{ dll *syscall.DLL }

func mapLibrary(path string) (platformLibrary, error) {
	dll, err := syscall.LoadDLL(path)
	if err != nil {
		return nil, fmt.Errorf("thsdk: 运行组件加载失败，请确认安装包与当前环境匹配: %w", err)
	}
	return &windowsLibrary{dll: dll}, nil
}

func (l *windowsLibrary) lookup(name string) (*funcPtr, error) {
	proc, err := l.dll.FindProc(name)
	if err != nil {
		return nil, fmt.Errorf("thsdk: 运行组件缺少导出 %s: %w", name, err)
	}
	return &funcPtr{addr: proc.Addr()}, nil
}

// readResult copies a NUL-terminated UTF-8 buffer owned by the component.
func (l *windowsLibrary) readResult(ptr unsafePointer) ([]byte, error) {
	if ptr == nil {
		return nil, fmt.Errorf("thsdk: 运行组件返回了空指针")
	}
	length := 0
	for {
		// The component cannot produce responses beyond its own write bound;
		// a much larger scan means the pointer is not a C string.
		if length > 64<<20 {
			return nil, fmt.Errorf("thsdk: 运行组件返回了未终止的结果缓冲区")
		}
		if *(*byte)(unsafe.Add(ptr, length)) == 0 {
			break
		}
		length++
	}
	out := make([]byte, length)
	copy(out, unsafe.Slice((*byte)(ptr), length))
	return out, nil
}

// unmap is deliberately empty; see mappedLibrary.release.
func (l *windowsLibrary) unmap() error { return nil }

func sleepDuration(d time.Duration) { time.Sleep(d) }
