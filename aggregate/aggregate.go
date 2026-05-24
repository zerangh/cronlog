// Package aggregate collects log entries and produces a structured summary
// grouped by log level, suitable for inclusion in webhook payloads or reports.
package aggregate

import (
	"sync"
	"time"
)

// Entry represents a single captured log entry.
type Entry struct {
	Level   string
	Message string
	Fields  map[string]any
	At      time.Time
}

// Summary holds aggregated counts and entries grouped by level.
type Summary struct {
	Total    int
	ByLevel  map[string][]Entry
	Counts   map[string]int
}

// Collector accumulates log entries in memory.
type Collector struct {
	mu      sync.Mutex
	entries []Entry
}

// New returns an initialised Collector.
func New() *Collector {
	return &Collector{}
}

// Add appends an entry to the collector.
func (c *Collector) Add(level, message string, fields map[string]any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e := Entry{
		Level:   level,
		Message: message,
		Fields:  fields,
		At:      time.Now(),
	}
	c.entries = append(c.entries, e)
}

// Summarise returns a Summary of all collected entries.
func (c *Collector) Summarise() Summary {
	c.mu.Lock()
	defer c.mu.Unlock()

	s := Summary{
		ByLevel: make(map[string][]Entry),
		Counts:  make(map[string]int),
		Total:   len(c.entries),
	}

	for _, e := range c.entries {
		s.ByLevel[e.Level] = append(s.ByLevel[e.Level], e)
		s.Counts[e.Level]++
	}

	return s
}

// Reset clears all collected entries.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = nil
}
