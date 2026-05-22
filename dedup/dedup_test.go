package dedup_test

import (
	"testing"
	"time"

	"github.com/cronlog/dedup"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := dedup.New(0)
	if err == nil {
		t.Fatal("expected error for zero window, got nil")
	}
	_, err = dedup.New(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative window, got nil")
	}
}

func TestNew_ValidWindow(t *testing.T) {
	d, err := dedup.New(5 * time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Deduplicator")
	}
}

func TestAllow_FirstOccurrenceAllowed(t *testing.T) {
	d, _ := dedup.New(5 * time.Second)
	if !d.Allow("job.started") {
		t.Error("expected first occurrence to be allowed")
	}
}

func TestAllow_DuplicateWithinWindowSuppressed(t *testing.T) {
	base := time.Now()
	d, _ := dedup.New(10 * time.Second)
	// inject fixed clock via unexported field workaround: use public API only
	d.Allow("db.error") // first — allowed

	// second call immediately after should be suppressed
	if d.Allow("db.error") {
		t.Error("expected duplicate within window to be suppressed")
	}
	_ = base
}

func TestAllow_AllowedAfterWindowExpires(t *testing.T) {
	d, _ := dedup.New(1 * time.Millisecond)
	d.Allow("timeout.warn")
	time.Sleep(5 * time.Millisecond)
	if !d.Allow("timeout.warn") {
		t.Error("expected message to be allowed after window expires")
	}
}

func TestSuppressed_CountsCorrectly(t *testing.T) {
	d, _ := dedup.New(10 * time.Second)
	d.Allow("repeated") // allowed, count = 0
	d.Allow("repeated") // suppressed, count = 1
	d.Allow("repeated") // suppressed, count = 2

	if got := d.Suppressed("repeated"); got != 2 {
		t.Errorf("expected suppressed=2, got %d", got)
	}
}

func TestSuppressed_UnknownKeyReturnsZero(t *testing.T) {
	d, _ := dedup.New(5 * time.Second)
	if got := d.Suppressed("unknown"); got != 0 {
		t.Errorf("expected 0 for unknown key, got %d", got)
	}
}

func TestReset_ClearsAllEntries(t *testing.T) {
	d, _ := dedup.New(10 * time.Second)
	d.Allow("event.x")
	d.Allow("event.x") // suppressed
	d.Reset()

	if !d.Allow("event.x") {
		t.Error("expected Allow to return true after Reset")
	}
	if got := d.Suppressed("event.x"); got != 0 {
		t.Errorf("expected suppressed=0 after Reset, got %d", got)
	}
}
