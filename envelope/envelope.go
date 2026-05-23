// Package envelope wraps a log entry with job metadata and execution context
// for use in webhook payloads and structured output.
package envelope

import (
	"fmt"
	"time"
)

// Entry represents a single log line captured during a job run.
type Entry struct {
	Level   string         `json:"level"`
	Message string         `json:"message"`
	Fields  map[string]any `json:"fields,omitempty"`
}

// Envelope wraps job execution metadata together with captured log entries.
type Envelope struct {
	JobName    string    `json:"job_name"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	DurationMs int64     `json:"duration_ms"`
	Success    bool      `json:"success"`
	ExitCode   int       `json:"exit_code"`
	Entries    []Entry   `json:"entries"`
}

// New constructs an Envelope from the provided job metadata and log entries.
func New(jobName string, startedAt, finishedAt time.Time, success bool, exitCode int, entries []Entry) Envelope {
	return Envelope{
		JobName:    jobName,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		DurationMs: finishedAt.Sub(startedAt).Milliseconds(),
		Success:    success,
		ExitCode:   exitCode,
		Entries:    entries,
	}
}

// ErrorEntries returns only the entries with level "error".
func (e Envelope) ErrorEntries() []Entry {
	var out []Entry
	for _, en := range e.Entries {
		if en.Level == "error" {
			out = append(out, en)
		}
	}
	return out
}

// Summary returns a short human-readable description of the envelope.
func (e Envelope) Summary() string {
	status := "succeeded"
	if !e.Success {
		status = "failed"
	}
	return fmt.Sprintf("%s %s in %s", e.JobName, status, formatDuration(e.DurationMs))
}

func formatDuration(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}
	return fmt.Sprintf("%.2fs", float64(ms)/1000)
}
