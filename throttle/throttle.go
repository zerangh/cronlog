// Package throttle provides a concurrency limiter that restricts the number
// of simultaneous log-processing goroutines, preventing resource exhaustion
// during high-volume cron job output.
package throttle

import (
	"context"
	"errors"
)

// ErrThrottled is returned when the throttle limit is reached and the context
// is cancelled before a slot becomes available.
var ErrThrottled = errors.New("throttle: limit reached, context cancelled")

// Throttle limits the number of concurrent operations.
type Throttle struct {
	sem chan struct{}
}

// New creates a Throttle that allows at most max concurrent operations.
// Returns an error if max is less than 1.
func New(max int) (*Throttle, error) {
	if max < 1 {
		return nil, errors.New("throttle: max must be at least 1")
	}
	return &Throttle{
		sem: make(chan struct{}, max),
	}, nil
}

// Acquire blocks until a slot is available or the context is cancelled.
// Returns ErrThrottled if the context is cancelled while waiting.
func (t *Throttle) Acquire(ctx context.Context) error {
	select {
	case t.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ErrThrottled
	}
}

// Release frees a previously acquired slot. It is a no-op if called more
// times than Acquire has succeeded.
func (t *Throttle) Release() {
	select {
	case <-t.sem:
	default:
	}
}

// Available returns the number of free slots remaining.
func (t *Throttle) Available() int {
	return cap(t.sem) - len(t.sem)
}

// Cap returns the maximum concurrency level.
func (t *Throttle) Cap() int {
	return cap(t.sem)
}
