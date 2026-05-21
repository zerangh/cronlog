package logger

import (
	"bytes"
	"strings"
	"testing"
)

func newTestLogger(name string) (*Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	return New(name, &buf), &buf
}

func TestInfo_WritesEntry(t *testing.T) {
	l, buf := newTestLogger("test-job")
	l.Info("starting up", nil)

	if len(l.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.Entries))
	}
	if l.Entries[0].Level != LevelInfo {
		t.Errorf("expected INFO, got %s", l.Entries[0].Level)
	}
	if !strings.Contains(buf.String(), "starting up") {
		t.Errorf("output missing message: %s", buf.String())
	}
}

func TestError_SetsHasErrors(t *testing.T) {
	l, _ := newTestLogger("test-job")
	if l.HasErrors() {
		t.Fatal("expected no errors initially")
	}
	l.Error("something broke", map[string]string{"code": "500"})
	if !l.HasErrors() {
		t.Fatal("expected HasErrors to be true after Error call")
	}
}

func TestWarn_DoesNotSetHasErrors(t *testing.T) {
	l, _ := newTestLogger("test-job")
	l.Warn("low disk", nil)
	if l.HasErrors() {
		t.Error("WARN should not set HasErrors")
	}
}

func TestSummary_ContainsJobName(t *testing.T) {
	l, _ := newTestLogger("my-cron-job")
	l.Info("done", nil)
	s := l.Summary()
	if !strings.Contains(s, "my-cron-job") {
		t.Errorf("summary missing job name: %s", s)
	}
}

func TestSummary_ContainsAllEntries(t *testing.T) {
	l, _ := newTestLogger("job")
	l.Info("step one", nil)
	l.Warn("step two", nil)
	l.Error("step three", nil)
	s := l.Summary()
	for _, msg := range []string{"step one", "step two", "step three"} {
		if !strings.Contains(s, msg) {
			t.Errorf("summary missing %q", msg)
		}
	}
}

func TestFields_AppearInOutput(t *testing.T) {
	l, buf := newTestLogger("job")
	l.Info("processed", map[string]string{"rows": "42"})
	if !strings.Contains(buf.String(), `rows="42"`) {
		t.Errorf("fields not found in output: %s", buf.String())
	}
}

func TestDuration_NonNegative(t *testing.T) {
	l, _ := newTestLogger("job")
	if l.Duration() < 0 {
		t.Error("duration should be non-negative")
	}
}
