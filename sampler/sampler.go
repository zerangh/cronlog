// Package sampler provides log entry sampling to reduce noise from
// high-frequency cron jobs by allowing only a fraction of log entries through.
package sampler

import (
	"errors"
	"sync/atomic"
)

// Sampler decides whether a given log entry should be emitted based on a
// deterministic counter-based sampling strategy. Every Nth entry is allowed
// through; all others are dropped.
type Sampler struct {
	rate    uint64
	counter atomic.Uint64
}

// ErrInvalidRate is returned when a sampling rate less than 1 is provided.
var ErrInvalidRate = errors.New("sampler: rate must be >= 1")

// New creates a Sampler that allows 1 out of every rate entries through.
// A rate of 1 allows all entries (no sampling). A rate of N allows every
// Nth entry.
func New(rate uint64) (*Sampler, error) {
	if rate < 1 {
		return nil, ErrInvalidRate
	}
	return &Sampler{rate: rate}, nil
}

// Allow returns true if the current entry should be emitted. It increments an
// internal counter on every call and returns true when the counter is a
// multiple of the configured rate.
func (s *Sampler) Allow() bool {
	n := s.counter.Add(1)
	return n%s.rate == 0
}

// Reset resets the internal counter to zero. Useful between job runs to ensure
// consistent sampling behaviour across executions.
func (s *Sampler) Reset() {
	s.counter.Store(0)
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() uint64 {
	return s.rate
}
