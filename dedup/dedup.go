// Package dedup provides log entry deduplication to suppress repeated
// identical messages within a configurable time window.
package dedup

import (
	"sync"
	"time"
)

// entry tracks the last time a message was seen and how many times it was suppressed.
type entry struct {
	lastSeen    time.Time
	suppressed  int
}

// Deduplicator suppresses repeated log messages within a sliding time window.
type Deduplicator struct {
	mu      sync.Mutex
	window  time.Duration
	seen    map[string]*entry
	now     func() time.Time
}

// New creates a Deduplicator with the given deduplication window.
// Messages with the same key seen within the window are suppressed.
func New(window time.Duration) (*Deduplicator, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Deduplicator{
		window: window,
		seen:   make(map[string]*entry),
		now:    time.Now,
	}, nil
}

// Allow returns true if the message with the given key should be logged.
// Duplicate messages within the window are suppressed and false is returned.
func (d *Deduplicator) Allow(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if e, ok := d.seen[key]; ok {
		if now.Sub(e.lastSeen) < d.window {
			e.suppressed++
			return false
		}
	}
	d.seen[key] = &entry{lastSeen: now}
	return true
}

// Suppressed returns the number of times the given key was suppressed since
// it was last allowed through. Returns 0 if the key is unknown.
func (d *Deduplicator) Suppressed(key string) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	if e, ok := d.seen[key]; ok {
		return e.suppressed
	}
	return 0
}

// Reset clears all tracked entries, allowing all messages through again.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]*entry)
}
