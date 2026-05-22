package sampler

import (
	"sync"
	"testing"
)

func TestNew_InvalidRate(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for rate=0, got nil")
	}
	if err != ErrInvalidRate {
		t.Fatalf("expected ErrInvalidRate, got %v", err)
	}
}

func TestNew_ValidRate(t *testing.T) {
	s, err := New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 3 {
		t.Fatalf("expected rate 3, got %d", s.Rate())
	}
}

func TestAllow_RateOne_AllowsAll(t *testing.T) {
	s, _ := New(1)
	for i := 0; i < 10; i++ {
		if !s.Allow() {
			t.Fatalf("expected Allow()=true at call %d with rate=1", i+1)
		}
	}
}

func TestAllow_RateN_AllowsEveryNth(t *testing.T) {
	const rate = 4
	s, _ := New(rate)

	allowed := 0
	for i := 0; i < 20; i++ {
		if s.Allow() {
			allowed++
		}
	}

	expected := 20 / rate
	if allowed != int(expected) {
		t.Fatalf("expected %d allowed entries, got %d", expected, allowed)
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	s, _ := New(5)

	// advance counter to 5 — the 5th call should be allowed
	for i := 0; i < 4; i++ {
		s.Allow()
	}
	if !s.Allow() {
		t.Fatal("expected 5th call to be allowed")
	}

	s.Reset()

	// after reset, next single call should NOT be allowed (counter=1, 1%5 != 0)
	if s.Allow() {
		t.Fatal("expected first call after reset to be denied")
	}
}

func TestAllow_ConcurrentSafe(t *testing.T) {
	s, _ := New(3)

	var wg sync.WaitGroup
	const goroutines = 50

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Allow()
		}()
	}
	wg.Wait()
}
