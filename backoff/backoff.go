// Package backoff provides configurable delay strategies for retry loops.
// It supports fixed, linear, and exponential backoff with optional jitter.
package backoff

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

// Strategy defines how delays are calculated between retry attempts.
type Strategy int

const (
	// Fixed applies the same base delay on every attempt.
	Fixed Strategy = iota
	// Linear increases the delay linearly: base * attempt.
	Linear
	// Exponential doubles the delay each attempt: base * 2^(attempt-1).
	Exponential
)

// Config holds the parameters for a backoff policy.
type Config struct {
	Strategy Strategy
	// Base is the initial delay duration.
	Base time.Duration
	// Max caps the computed delay. Zero means no cap.
	Max time.Duration
	// Jitter adds a random fraction of the computed delay when true.
	Jitter bool
}

// ErrInvalidBase is returned when Base is not positive.
var ErrInvalidBase = errors.New("backoff: base duration must be positive")

// Validate checks that the Config is well-formed.
func (c Config) Validate() error {
	if c.Base <= 0 {
		return ErrInvalidBase
	}
	return nil
}

// Delay returns the computed wait duration for the given attempt number
// (1-indexed). It applies the configured strategy, cap, and jitter.
func (c Config) Delay(attempt int) (time.Duration, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}
	if attempt < 1 {
		attempt = 1
	}

	var d time.Duration
	switch c.Strategy {
	case Linear:
		d = c.Base * time.Duration(attempt)
	case Exponential:
		mult := math.Pow(2, float64(attempt-1))
		d = time.Duration(float64(c.Base) * mult)
	default: // Fixed
		d = c.Base
	}

	if c.Max > 0 && d > c.Max {
		d = c.Max
	}

	if c.Jitter && d > 0 {
		// Add up to 50 % of d as random jitter.
		jitter := time.Duration(rand.Int63n(int64(d) / 2))
		d += jitter
	}

	return d, nil
}
