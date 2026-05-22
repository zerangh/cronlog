// Package retry provides configurable retry logic for cron job execution.
package retry

import (
	"context"
	"fmt"
	"time"
)

// Policy defines the retry behaviour for a job.
type Policy struct {
	// MaxAttempts is the total number of times the job will be tried (including
	// the first attempt). A value of 1 means no retries.
	MaxAttempts int

	// Delay is the duration to wait between consecutive attempts.
	Delay time.Duration
}

// Validate returns an error if the policy contains invalid values.
func (p Policy) Validate() error {
	if p.MaxAttempts < 1 {
		return fmt.Errorf("retry: MaxAttempts must be at least 1, got %d", p.MaxAttempts)
	}
	if p.Delay < 0 {
		return fmt.Errorf("retry: Delay must be non-negative, got %s", p.Delay)
	}
	return nil
}

// JobFunc is the function signature expected by Run.
type JobFunc func(ctx context.Context) error

// Run executes fn according to p, retrying on non-nil errors until the attempt
// limit is reached or ctx is cancelled. It returns the last error encountered,
// or nil on success.
func Run(ctx context.Context, p Policy, fn JobFunc) error {
	if err := p.Validate(); err != nil {
		return err
	}

	var last error
	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry: context cancelled before attempt %d: %w", attempt, err)
		}

		last = fn(ctx)
		if last == nil {
			return nil
		}

		if attempt < p.MaxAttempts {
			select {
			case <-time.After(p.Delay):
			case <-ctx.Done():
				return fmt.Errorf("retry: context cancelled during delay after attempt %d: %w", attempt, ctx.Err())
			}
		}
	}

	return fmt.Errorf("retry: all %d attempt(s) failed, last error: %w", p.MaxAttempts, last)
}
