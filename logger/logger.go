package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a log entry.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Entry holds a single structured log record.
type Entry struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Fields    map[string]string
}

// Logger captures structured output for a cron job run.
type Logger struct {
	JobName string
	Entries []Entry
	out     io.Writer
	startAt time.Time
}

// New creates a Logger that writes human-readable output to w.
// Pass nil to use os.Stdout.
func New(jobName string, w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{
		JobName: jobName,
		out:     w,
		startAt: time.Now(),
	}
}

// log appends an entry and writes it to the output writer.
func (l *Logger) log(level Level, msg string, fields map[string]string) {
	e := Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    fields,
	}
	l.Entries = append(l.Entries, e)
	fmt.Fprintln(l.out, formatEntry(e))
}

func (l *Logger) Info(msg string, fields map[string]string)  { l.log(LevelInfo, msg, fields) }
func (l *Logger) Warn(msg string, fields map[string]string)  { l.log(LevelWarn, msg, fields) }
func (l *Logger) Error(msg string, fields map[string]string) { l.log(LevelError, msg, fields) }

// HasErrors returns true if any ERROR-level entries were recorded.
func (l *Logger) HasErrors() bool {
	for _, e := range l.Entries {
		if e.Level == LevelError {
			return true
		}
	}
	return false
}

// Duration returns the elapsed time since the logger was created.
func (l *Logger) Duration() time.Duration {
	return time.Since(l.startAt)
}

// Summary returns a plain-text summary of all captured entries.
func (l *Logger) Summary() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "=== Job: %s | duration: %s ===\n", l.JobName, l.Duration().Round(time.Millisecond))
	for _, e := range l.Entries {
		fmt.Fprintln(&buf, formatEntry(e))
	}
	return buf.String()
}

func formatEntry(e Entry) string {
	base := fmt.Sprintf("%s [%s] %s", e.Timestamp.Format(time.RFC3339), e.Level, e.Message)
	for k, v := range e.Fields {
		base += fmt.Sprintf(" %s=%q", k, v)
	}
	return base
}
