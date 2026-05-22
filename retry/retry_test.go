package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cronlog/retry"
)

var errBoom = errors.New("boom")

func TestRun_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Run(context.Background(), retry.Policy{MaxAttempts: 3, Delay: 0}, func(_ context.Context) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRun_RetriesOnError(t *testing.T) {
	calls := 0
	err := retry.Run(context.Background(), retry.Policy{MaxAttempts: 3, Delay: 0}, func(_ context.Context) error {
		calls++
		if calls < 3 {
			return errBoom
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error after eventual success, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRun_ExhaustsAttempts(t *testing.T) {
	calls := 0
	err := retry.Run(context.Background(), retry.Policy{MaxAttempts: 2, Delay: 0}, func(_ context.Context) error {
		calls++
		return errBoom
	})
	if err == nil {
		t.Fatal("expected an error when all attempts fail")
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
	if !errors.Is(err, errBoom) {
		t.Fatalf("expected wrapped errBoom, got %v", err)
	}
}

func TestRun_ContextCancelledDuringDelay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	err := retry.Run(ctx, retry.Policy{MaxAttempts: 5, Delay: 200 * time.Millisecond}, func(_ context.Context) error {
		calls++
		if calls == 1 {
			cancel()
		}
		return errBoom
	})
	if err == nil {
		t.Fatal("expected error after context cancellation")
	}
	if calls != 1 {
		t.Fatalf("expected exactly 1 call before cancellation, got %d", calls)
	}
}

func TestValidate_InvalidMaxAttempts(t *testing.T) {
	p := retry.Policy{MaxAttempts: 0, Delay: 0}
	if err := p.Validate(); err == nil {
		t.Fatal("expected validation error for MaxAttempts=0")
	}
}

func TestValidate_NegativeDelay(t *testing.T) {
	p := retry.Policy{MaxAttempts: 1, Delay: -time.Second}
	if err := p.Validate(); err == nil {
		t.Fatal("expected validation error for negative Delay")
	}
}
