package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/cronlog/output"
)

func newEntry() output.Entry {
	return output.Entry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     "INFO",
		Message:   "job started",
		JobName:   "backup",
	}
}

func TestWriter_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatText)

	if err := w.Write(newEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "INFO") {
		t.Errorf("expected level in output, got: %s", got)
	}
	if !strings.Contains(got, "backup") {
		t.Errorf("expected job name in output, got: %s", got)
	}
	if !strings.Contains(got, "job started") {
		t.Errorf("expected message in output, got: %s", got)
	}
}

func TestWriter_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatJSON)

	if err := w.Write(newEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got output.Entry
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got.Level != "INFO" {
		t.Errorf("expected level INFO, got %s", got.Level)
	}
	if got.JobName != "backup" {
		t.Errorf("expected job name backup, got %s", got.JobName)
	}
}

func TestWriter_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, "unknown")

	if err := w.Write(newEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if strings.HasPrefix(got, "{") {
		t.Errorf("expected text format, got JSON-like output: %s", got)
	}
}
