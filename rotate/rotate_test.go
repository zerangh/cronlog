package rotate_test

import (
	"sync"
	"testing"

	"cronlog/rotate"
)

func makeEntry(msg string) rotate.Entry {
	return rotate.Entry{Level: "info", Message: msg, Fields: map[string]any{"k": "v"}}
}

func TestNew_InvalidCapacity(t *testing.T) {
	_, err := rotate.New(0)
	if err == nil {
		t.Fatal("expected error for zero capacity")
	}
}

func TestNew_ValidCapacity(t *testing.T) {
	r, err := rotate.New(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 0 {
		t.Fatalf("expected empty rotator, got len %d", r.Len())
	}
}

func TestAdd_StoresEntry(t *testing.T) {
	r, _ := rotate.New(3)
	r.Add(makeEntry("hello"))
	if r.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", r.Len())
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	r, _ := rotate.New(3)
	for i, msg := range []string{"a", "b", "c", "d"} {
		r.Add(makeEntry(msg))
		_ = i
	}
	entries := r.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Message != "b" {
		t.Errorf("expected oldest entry to be 'b', got %q", entries[0].Message)
	}
	if entries[2].Message != "d" {
		t.Errorf("expected newest entry to be 'd', got %q", entries[2].Message)
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	r, _ := rotate.New(5)
	r.Add(makeEntry("x"))
	copy1 := r.Entries()
	copy1[0].Message = "mutated"
	copy2 := r.Entries()
	if copy2[0].Message == "mutated" {
		t.Error("Entries should return an independent copy")
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	r, _ := rotate.New(5)
	r.Add(makeEntry("one"))
	r.Add(makeEntry("two"))
	r.Reset()
	if r.Len() != 0 {
		t.Fatalf("expected 0 entries after reset, got %d", r.Len())
	}
}

func TestAdd_ConcurrentWrites(t *testing.T) {
	r, _ := rotate.New(50)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			r.Add(makeEntry("concurrent"))
		}(i)
	}
	wg.Wait()
	if r.Len() > 50 {
		t.Errorf("rotator exceeded capacity: len=%d", r.Len())
	}
}
