// Package output provides structured log output formatting for cronlog.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Format represents the output format for log entries.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// Entry represents a single log line to be written.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	JobName   string    `json:"job_name"`
}

// Writer writes log entries to an io.Writer in the configured format.
type Writer struct {
	out    io.Writer
	format Format
}

// NewWriter creates a new Writer with the given output destination and format.
func NewWriter(out io.Writer, format Format) *Writer {
	return &Writer{out: out, format: format}
}

// Write formats and writes a single Entry to the underlying writer.
func (w *Writer) Write(e Entry) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(e)
	default:
		return w.writeText(e)
	}
}

func (w *Writer) writeJSON(e Entry) error {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("output: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(w.out, "%s\n", b)
	return err
}

func (w *Writer) writeText(e Entry) error {
	_, err := fmt.Fprintf(w.out, "%s [%s] [%s] %s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.JobName,
		e.Message,
	)
	return err
}
