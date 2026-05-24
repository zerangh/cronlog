package throttle_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cronlog/throttle"
)

func TestNew_InvalidMax(t *testing.T) {
	_, err := throttle.New(0)
	if err == nil {
		t.Fatal("expected error for max=0, got nil")
	}
}

func TestNew_ValidMax(t *testing.T) {
	th, err := throttle.New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.Cap() != 3 {
		t.Fatalf("expected cap 3, got %d", th.Cap())
	}
}

func TestAcquire_GrantsSlot(t *testing.T) {
	th, _ := throttle.New(2)
	ctx := context.Background()

	if err := th.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.Available() != 1 {
		t.Fatalf("expected 1 available, got %d", th.Available())
	}
}

func TestRelease_FreesSlot(t *testing.T) {
	th, _ := throttle.New(1)
	ctx := context.Background()

	_ = th.Acquire(ctx)
	th.Release()

	if th.Available() != 1 {
		t.Fatalf("expected 1 available after release, got %d", th.Available())
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	th, _ := throttle.New(1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = th.Acquire(context.Background()) // fill the slot

	err := th.Acquire(ctx)
	if err != throttle.ErrThrottled {
		t.Fatalf("expected ErrThrottled, got %v", err)
	}
}

func TestAcquire_ConcurrentSafety(t *testing.T) {
	const workers = 20
	const limit = 5

	th, _ := throttle.New(limit)
	ctx := context.Background()

	var wg sync.WaitGroup
	var mu sync.Mutex
	peak := 0
	current := 0

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := th.Acquire(ctx); err != nil {
				return
			}
			defer th.Release()

			mu.Lock()
			current++
			if current > peak {
				peak = current
			}
			mu.Unlock()

			time.Sleep(5 * time.Millisecond)

			mu.Lock()
			current--
			mu.Unlock()
		}()
	}
	wg.Wait()

	if peak > limit {
		t.Fatalf("peak concurrency %d exceeded limit %d", peak, limit)
	}
}
