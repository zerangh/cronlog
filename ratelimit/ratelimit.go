// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how frequently webhook notifications are dispatched.
package ratelimit

import (
	"errors"
	"sync"
	"time"
)

// ErrRateLimited is returned when a call is rejected due to rate limiting.
var ErrRateLimited = errors.New("ratelimit: too many requests")

// Limiter controls the rate of events using a token-bucket algorithm.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	now      func() time.Time
}

// Config holds configuration for a Limiter.
type Config struct {
	// Max is the maximum number of tokens (burst size).
	Max int
	// PerSecond is the token replenishment rate.
	PerSecond float64
}

// New creates a Limiter with the given configuration.
// Max must be >= 1 and PerSecond must be > 0.
func New(cfg Config) (*Limiter, error) {
	if cfg.Max < 1 {
		return nil, errors.New("ratelimit: Max must be at least 1")
	}
	if cfg.PerSecond <= 0 {
		return nil, errors.New("ratelimit: PerSecond must be greater than 0")
	}
	return &Limiter{
		tokens:   float64(cfg.Max),
		max:      float64(cfg.Max),
		rate:     cfg.PerSecond,
		lastTick: time.Now(),
		now:      time.Now,
	}, nil
}

// Allow reports whether an event may proceed. It is safe for concurrent use.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Do calls fn if the rate limit allows, otherwise returns ErrRateLimited.
func (l *Limiter) Do(fn func() error) error {
	if !l.Allow() {
		return ErrRateLimited
	}
	return fn()
}
