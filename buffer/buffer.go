// Package buffer provides an in-memory log entry buffer with a configurable
// capacity cap. When the buffer is full, the oldest entries are evicted to
// make room for new ones (ring-buffer semantics).
package buffer

import "sync"

// Entry represents a single buffered log line.
type Entry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Fields  map[string]any `json:"fields,omitempty"`
}

// Buffer holds log entries up to a fixed capacity.
type Buffer struct {
	mu       sync.Mutex
	entries  []Entry
	capacity int
}

// New returns a Buffer that keeps at most capacity entries.
// If capacity is less than 1 it defaults to 100.
func New(capacity int) *Buffer {
	if capacity < 1 {
		capacity = 100
	}
	return &Buffer{
		entries:  make([]Entry, 0, capacity),
		capacity: capacity,
	}
}

// Add appends an entry to the buffer, evicting the oldest entry when full.
func (b *Buffer) Add(e Entry) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.entries) >= b.capacity {
		// evict oldest
		b.entries = b.entries[1:]
	}
	b.entries = append(b.entries, e)
}

// Entries returns a shallow copy of all buffered entries.
func (b *Buffer) Entries() []Entry {
	b.mu.Lock()
	defer b.mu.Unlock()

	copy := make([]Entry, len(b.entries))
	for i, e := range b.entries {
		copy[i] = e
	}
	return copy
}

// Len returns the current number of buffered entries.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.entries)
}

// Reset discards all buffered entries.
func (b *Buffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = b.entries[:0]
}
