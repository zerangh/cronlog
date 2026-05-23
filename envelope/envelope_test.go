package envelope_test

import (
	"testing"
	"time"

	"github.com/cronlog/envelope"
)

var (
	t0 = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	t1 = t0.Add(2500 * time.Millisecond)
)

func baseEntries() []envelope.Entry {
	return []envelope.Entry{
		{Level: "info", Message: "starting"},
		{Level: "error", Message: "connection refused", Fields: map[string]any{"host": "db"}},
		{Level: "warn", Message: "retrying"},
		{Level: "error", Message: "timeout"},
	}
}

func TestNew_PopulatesFields(t *testing.T) {
	e := envelope.New("backup", t0, t1, false, 1, baseEntries())

	if e.JobName != "backup" {
		t.Errorf("expected job name 'backup', got %q", e.JobName)
	}
	if e.DurationMs != 2500 {
		t.Errorf("expected duration 2500ms, got %d", e.DurationMs)
	}
	if e.Success {
		t.Error("expected success=false")
	}
	if e.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", e.ExitCode)
	}
}

func TestErrorEntries_FiltersCorrectly(t *testing.T) {
	e := envelope.New("backup", t0, t1, false, 1, baseEntries())
	errs := e.ErrorEntries()

	if len(errs) != 2 {
		t.Fatalf("expected 2 error entries, got %d", len(errs))
	}
	for _, en := range errs {
		if en.Level != "error" {
			t.Errorf("unexpected level %q in error entries", en.Level)
		}
	}
}

func TestErrorEntries_EmptyWhenNone(t *testing.T) {
	entries := []envelope.Entry{
		{Level: "info", Message: "ok"},
	}
	e := envelope.New("ping", t0, t1, true, 0, entries)
	if got := e.ErrorEntries(); len(got) != 0 {
		t.Errorf("expected no error entries, got %d", len(got))
	}
}

func TestSummary_SuccessSeconds(t *testing.T) {
	e := envelope.New("sync", t0, t1, true, 0, nil)
	got := e.Summary()
	expected := "sync succeeded in 2.50s"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSummary_FailureMilliseconds(t *testing.T) {
	start := t0
	end := t0.Add(450 * time.Millisecond)
	e := envelope.New("check", start, end, false, 2, nil)
	got := e.Summary()
	expected := "check failed in 450ms"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
