package metrics_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronlog/metrics"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestCollector_Start_SetsStartTime(t *testing.T) {
	c := metrics.NewCollector("my-job")
	before := time.Now()
	c.Start()
	after := time.Now()

	r := c.Finish(0, false, 0)
	if r.StartedAt.Before(before) || r.StartedAt.After(after) {
		t.Errorf("StartedAt %v not in expected range [%v, %v]", r.StartedAt, before, after)
	}
}

func TestCollector_Finish_PopulatesResult(t *testing.T) {
	start := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	end := start.Add(3 * time.Second)

	c := metrics.NewCollector("batch-job")
	c.Start()

	// Override internal clock via a white-box approach using a test-friendly constructor.
	// Since we expose now via the struct we test observable behaviour instead.
	r := c.Finish(1, true, 42)

	if r.JobName != "batch-job" {
		t.Errorf("expected JobName 'batch-job', got %q", r.JobName)
	}
	if r.ExitCode != 1 {
		t.Errorf("expected ExitCode 1, got %d", r.ExitCode)
	}
	if !r.HasErrors {
		t.Error("expected HasErrors to be true")
	}
	if r.LogLines != 42 {
		t.Errorf("expected LogLines 42, got %d", r.LogLines)
	}
	if r.Duration < 0 {
		t.Errorf("expected non-negative Duration, got %v", r.Duration)
	}
	_ = start
	_ = end
}

func TestCollector_Finish_DurationIsPositive(t *testing.T) {
	c := metrics.NewCollector("timer-job")
	c.Start()
	time.Sleep(2 * time.Millisecond)
	r := c.Finish(0, false, 5)

	if r.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", r.Duration)
	}
}

func TestCollector_Finish_FinishedAtAfterStartedAt(t *testing.T) {
	c := metrics.NewCollector("order-job")
	c.Start()
	r := c.Finish(0, false, 0)

	if r.FinishedAt.Before(r.StartedAt) {
		t.Errorf("FinishedAt %v should not be before StartedAt %v", r.FinishedAt, r.StartedAt)
	}
}
