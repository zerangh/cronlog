package backoff

import (
	"testing"
	"time"
)

func TestValidate_InvalidBase(t *testing.T) {
	c := Config{Base: 0}
	if err := c.Validate(); err != ErrInvalidBase {
		t.Fatalf("expected ErrInvalidBase, got %v", err)
	}
}

func TestValidate_Valid(t *testing.T) {
	c := Config{Base: time.Second}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelay_Fixed(t *testing.T) {
	c := Config{Strategy: Fixed, Base: 100 * time.Millisecond}
	for _, attempt := range []int{1, 2, 5} {
		d, err := c.Delay(attempt)
		if err != nil {
			t.Fatalf("attempt %d: unexpected error: %v", attempt, err)
		}
		if d != 100*time.Millisecond {
			t.Errorf("attempt %d: expected 100ms, got %v", attempt, d)
		}
	}
}

func TestDelay_Linear(t *testing.T) {
	c := Config{Strategy: Linear, Base: 100 * time.Millisecond}
	cases := []struct {
		attempt int
		want    time.Duration
	}{
		{1, 100 * time.Millisecond},
		{2, 200 * time.Millisecond},
		{3, 300 * time.Millisecond},
	}
	for _, tc := range cases {
		d, err := c.Delay(tc.attempt)
		if err != nil {
			t.Fatalf("attempt %d: unexpected error: %v", tc.attempt, err)
		}
		if d != tc.want {
			t.Errorf("attempt %d: expected %v, got %v", tc.attempt, tc.want, d)
		}
	}
}

func TestDelay_Exponential(t *testing.T) {
	c := Config{Strategy: Exponential, Base: 100 * time.Millisecond}
	cases := []struct {
		attempt int
		want    time.Duration
	}{
		{1, 100 * time.Millisecond},
		{2, 200 * time.Millisecond},
		{3, 400 * time.Millisecond},
		{4, 800 * time.Millisecond},
	}
	for _, tc := range cases {
		d, err := c.Delay(tc.attempt)
		if err != nil {
			t.Fatalf("attempt %d: unexpected error: %v", tc.attempt, err)
		}
		if d != tc.want {
			t.Errorf("attempt %d: expected %v, got %v", tc.attempt, tc.want, d)
		}
	}
}

func TestDelay_MaxCap(t *testing.T) {
	c := Config{Strategy: Exponential, Base: 100 * time.Millisecond, Max: 300 * time.Millisecond}
	d, err := c.Delay(4) // would be 800ms without cap
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d > 300*time.Millisecond {
		t.Errorf("expected delay <= 300ms, got %v", d)
	}
}

func TestDelay_JitterIncreasesDelay(t *testing.T) {
	c := Config{Strategy: Fixed, Base: 200 * time.Millisecond, Jitter: true}
	d, err := c.Delay(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d < 200*time.Millisecond {
		t.Errorf("jitter should not reduce delay below base: got %v", d)
	}
	if d > 300*time.Millisecond {
		t.Errorf("jitter should not exceed base + 50%%: got %v", d)
	}
}

func TestDelay_InvalidBase_ReturnsError(t *testing.T) {
	c := Config{Base: -1}
	_, err := c.Delay(1)
	if err != ErrInvalidBase {
		t.Fatalf("expected ErrInvalidBase, got %v", err)
	}
}
