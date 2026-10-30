package native

import (
	"sync"
	"time"
)

// platformLibrary is the mapped-image half of the binding; the cache around it
// decides when a platform may actually unmap.
type platformLibrary interface {
	lookup(name string) (*funcPtr, error)
	readResult(ptr unsafePointer) ([]byte, error)
	unmap() error
}

// mappedLibrary keeps the component resident. A Go c-shared image that has been
// unloaded once cannot be reloaded safely in the same process, so the mapping is
// reference-counted and only ever released at process exit.
type mappedLibrary struct {
	mu       sync.Mutex
	platform platformLibrary
	refs     int

	// The component throttles the whole process, not one handle, so the spacing
	// clock is shared by every Runtime mapped from the same file.
	paceMu   sync.Mutex
	lastCall time.Time
}

var libraryCache = struct {
	sync.Mutex
	paths map[string]*mappedLibrary
}{paths: map[string]*mappedLibrary{}}

func retainLibrary(path string) (*mappedLibrary, error) {
	libraryCache.Lock()
	defer libraryCache.Unlock()
	if entry, ok := libraryCache.paths[path]; ok {
		entry.mu.Lock()
		entry.refs++
		entry.mu.Unlock()
		return entry, nil
	}
	platform, err := mapLibrary(path)
	if err != nil {
		return nil, err
	}
	entry := &mappedLibrary{platform: platform, refs: 1}
	libraryCache.paths[path] = entry
	return entry, nil
}

func (m *mappedLibrary) release() error {
	m.mu.Lock()
	m.refs--
	last := m.refs <= 0
	m.mu.Unlock()
	if !last {
		return nil
	}
	// Keep the image mapped even at zero references: reloading is unsafe, and a
	// later Open() should reuse it rather than pay for a fresh mapping.
	return nil
}

func (m *mappedLibrary) lookup(name string) (*funcPtr, error) {
	return m.platform.lookup(name)
}

func (m *mappedLibrary) readResult(ptr unsafePointer) ([]byte, error) {
	return m.platform.readResult(ptr)
}

// pace blocks until the shared minimum interval since the last call has elapsed.
func (m *mappedLibrary) pace(spacing time.Duration) {
	if spacing <= 0 {
		return
	}
	m.paceMu.Lock()
	defer m.paceMu.Unlock()
	if !m.lastCall.IsZero() {
		if wait := spacing - time.Since(m.lastCall); wait > 0 {
			sleepDuration(wait)
		}
	}
	m.lastCall = time.Now()
}
