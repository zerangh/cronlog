package runner

import (
	"context"
	"fmt"
	"time"
)

// runWithTimeout executes fn within the given duration.
// If the deadline is exceeded, the context is cancelled and an error is returned.
// A zero or negative timeout means no limit is applied.
func runWithTimeout(timeout time.Duration, fn func() error) error {
	if timeout <= 0 {
		return fn()
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	type result struct {
		err error
	}

	ch := make(chan result, 1)

	go func() {
		ch <- result{err: fn()}
	}()

	select {
	case res := <-ch:
		return res.err
	case <-ctx.Done():
		return fmt.Errorf("job timed out after %s", timeout)
	}
}
