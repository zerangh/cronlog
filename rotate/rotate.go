// Package rotate provides log entry rotation based on a maximum entry count.
// When the limit is reached, older entries are discarded to keep memory usage bounded.
package rotate

import (
	"errors"
	"sync"
)

// Entry represents a single log entry stored by the rotator.
type Entry struct {
	Level   string
	Message string
	Fields  map[string]any
}

// Rotator holds a fixed-size ring of log entries, discarding the oldest
// when the capacity is exceeded.
type Rotator struct {
	mu       sync.Mutex
	entries  []Entry
	capacity int
}

// New creates a new Rotator with the given maximum capacity.
// Returns an error if capacity is less than 1.
func New(capacity int) (*Rotator, error) {
	if capacity < 1 {
		return nil, errors.New("rotate: capacity must be at least 1")
	}
	return &Rotator{
		entries:  make([]Entry, 0, capacity),
		capacity: capacity,
	}, nil
}

// Add appends an entry to the rotator. If the rotator is at capacity,
// the oldest entry is evicted before the new one is stored.
func (r *Rotator) Add(e Entry) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.entries) >= r.capacity {
		r.entries = r.entries[1:]
	}
	r.entries = append(r.entries, e)
}

// Entries returns a copy of all currently held entries in insertion order.
func (r *Rotator) Entries() []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := make([]Entry, len(r.entries))
	for i, e := range r.entries {
		copy[i] = e
	}
	return copy
}

// Len returns the number of entries currently stored.
func (r *Rotator) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}

// Reset removes all stored entries without changing the capacity.
func (r *Rotator) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = r.entries[:0]
}
