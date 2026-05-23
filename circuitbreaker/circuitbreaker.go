// Package circuitbreaker provides a simple circuit breaker for webhook
// notifications and other external calls made during cron job execution.
//
// The circuit breaker transitions between three states:
//
//   - Closed: normal operation, calls are allowed through.
//   - Open: failure threshold exceeded, calls are rejected immediately.
//   - HalfOpen: a single probe call is allowed to test recovery.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is open and the call is
// rejected without executing the underlying operation.
var ErrOpen = errors.New("circuitbreaker: circuit is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // rejecting calls
	StateHalfOpen              // testing recovery
)

// Breaker is a circuit breaker that tracks consecutive failures and opens
// the circuit when the failure threshold is exceeded.
type Breaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	threshold    int
	resetTimeout time.Duration
	openedAt     time.Time
	now          func() time.Time
}

// New creates a Breaker that opens after threshold consecutive failures and
// attempts recovery after resetTimeout.
func New(threshold int, resetTimeout time.Duration) (*Breaker, error) {
	if threshold <= 0 {
		return nil, errors.New("circuitbreaker: threshold must be positive")
	}
	if resetTimeout <= 0 {
		return nil, errors.New("circuitbreaker: resetTimeout must be positive")
	}
	return &Breaker{
		threshold:    threshold,
		resetTimeout: resetTimeout,
		now:          time.Now,
	}, nil
}

// Allow returns nil if the call should proceed, or ErrOpen if the circuit is
// open. Callers must follow up with a call to RecordSuccess or RecordFailure.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateOpen:
		if b.now().Sub(b.openedAt) >= b.resetTimeout {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	default:
		return nil
	}
}

// RecordSuccess records a successful call. If the circuit is half-open it
// transitions back to closed.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure records a failed call and may open the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.state == StateHalfOpen || b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = b.now()
	}
}

// State returns the current state of the circuit breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
