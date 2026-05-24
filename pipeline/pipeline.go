// Package pipeline chains multiple log entry processors together,
// passing each entry through a sequence of handlers in order.
package pipeline

import "fmt"

// Entry represents a single log entry passed through the pipeline.
type Entry struct {
	Level   string
	Message string
	Fields  map[string]string
}

// Processor is a function that transforms or filters a log entry.
// Returning (entry, false) drops the entry from the pipeline.
type Processor func(entry Entry) (Entry, bool)

// Pipeline holds an ordered sequence of processors.
type Pipeline struct {
	processors []Processor
}

// New creates a Pipeline with the given processors applied in order.
// Returns an error if no processors are provided.
func New(processors ...Processor) (*Pipeline, error) {
	if len(processors) == 0 {
		return nil, fmt.Errorf("pipeline: at least one processor is required")
	}
	return &Pipeline{processors: processors}, nil
}

// Run passes entry through each processor in sequence.
// If any processor drops the entry (returns false), Run returns the
// partially-transformed entry and false immediately.
func (p *Pipeline) Run(entry Entry) (Entry, bool) {
	for _, proc := range p.processors {
		var keep bool
		entry, keep = proc(entry)
		if !keep {
			return entry, false
		}
	}
	return entry, true
}

// Len returns the number of processors in the pipeline.
func (p *Pipeline) Len() int {
	return len(p.processors)
}
