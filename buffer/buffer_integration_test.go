package buffer_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/example/cronlog/buffer"
)

// TestAdd_ConcurrentWrites verifies that concurrent Add calls do not race or
// corrupt internal state.
func TestAdd_ConcurrentWrites(t *testing.T) {
	const goroutines = 20
	const perGoroutine = 50

	b := buffer.New(500)

	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				b.Add(buffer.Entry{
					Level:   "info",
					Message: fmt.Sprintf("goroutine %d entry %d", id, i),
				})
			}
		}(g)
	}
	wg.Wait()

	total := goroutines * perGoroutine // 1000
	want := 500                        // capped at capacity
	if b.Len() != want {
		t.Errorf("expected %d entries (capacity), got %d (total written: %d)", want, b.Len(), total)
	}
}

// TestAdd_EvictionPreservesCapacity checks that repeated over-capacity writes
// never cause Len to exceed the configured capacity.
func TestAdd_EvictionPreservesCapacity(t *testing.T) {
	cap := 10
	b := buffer.New(cap)

	for i := 0; i < 100; i++ {
		b.Add(buffer.Entry{Message: fmt.Sprintf("entry %d", i)})
		if b.Len() > cap {
			t.Fatalf("Len %d exceeded capacity %d after %d inserts", b.Len(), cap, i+1)
		}
	}
}
