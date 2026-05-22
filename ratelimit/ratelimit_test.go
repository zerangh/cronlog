package ratelimit_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/cronlog/ratelimit"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_InvalidMax(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Config{Max: 0, PerSecond: 1})
	if err == nil {
		t.Fatal("expected error for Max=0")
	}
}

func TestNew_InvalidRate(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Config{Max: 1, PerSecond: 0})
	if err == nil {
		t.Fatal("expected error for PerSecond=0")
	}
}

func TestAllow_BurstConsumption(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Config{Max: 3, PerSecond: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected Allow()=false after burst exhausted")
	}
}

func TestAllow_TokensReplenishOverTime(t *testing.T) {
	base := time.Now()
	l, err := ratelimit.New(ratelimit.Config{Max: 1, PerSecond: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Drain the single token.
	_ = l.Allow()

	// Advance internal clock by injecting via unexported field isn't possible;
	// instead verify that a real 600ms pause replenishes at rate=2.
	_ = base // suppress unused warning
	time.Sleep(600 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("expected token to be replenished after sleep")
	}
}

func TestDo_CallsFnWhenAllowed(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{Max: 1, PerSecond: 1})

	var called int32
	err := l.Do(func() error {
		atomic.StoreInt32(&called, 1)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("expected fn to be called")
	}
}

func TestDo_ReturnsErrRateLimitedWhenExhausted(t *testing.T) {
	l, _ := ratelimit.New(ratelimit.Config{Max: 1, PerSecond: 1})
	_ = l.Allow() // drain

	err := l.Do(func() error { return nil })
	if err != ratelimit.ErrRateLimited {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}
