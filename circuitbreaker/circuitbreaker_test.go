package circuitbreaker

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNew_InvalidResetTimeout(t *testing.T) {
	_, err := New(3, 0)
	if err == nil {
		t.Fatal("expected error for zero resetTimeout")
	}
}

func TestAllow_ClosedByDefault(t *testing.T) {
	b, _ := New(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b, _ := New(3, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != StateClosed {
		t.Fatal("expected closed before threshold")
	}
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatalf("expected open, got %v", b.State())
	}
}

func TestAllow_RejectsWhenOpen(t *testing.T) {
	b, _ := New(1, time.Minute)
	b.RecordFailure()
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_HalfOpenAfterTimeout(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	b, _ := New(1, time.Minute)
	b.now = fixedNow(base)
	b.RecordFailure()

	// advance past reset timeout
	b.now = fixedNow(base.Add(2 * time.Minute))
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil in half-open, got %v", err)
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("expected half-open, got %v", b.State())
	}
}

func TestRecordSuccess_ClosesCircuit(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	b, _ := New(1, time.Minute)
	b.now = fixedNow(base)
	b.RecordFailure()

	b.now = fixedNow(base.Add(2 * time.Minute))
	_ = b.Allow() // transition to half-open
	b.RecordSuccess()

	if b.State() != StateClosed {
		t.Fatalf("expected closed after success, got %v", b.State())
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil after close, got %v", err)
	}
}

func TestRecordFailure_HalfOpenReopens(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	b, _ := New(1, time.Minute)
	b.now = fixedNow(base)
	b.RecordFailure()

	b.now = fixedNow(base.Add(2 * time.Minute))
	_ = b.Allow()
	b.RecordFailure() // fail again in half-open

	if b.State() != StateOpen {
		t.Fatalf("expected open after half-open failure, got %v", b.State())
	}
}
