// Package output handles formatting and writing of structured log entries
// produced by cronlog jobs.
//
// Two formats are supported:
//
//   - text: human-readable lines suitable for terminal output or syslog.
//     Each line follows the pattern:
//     <RFC3339 timestamp> [LEVEL] [job_name] message
//
//   - json: newline-delimited JSON objects, useful for log aggregation
//     pipelines (e.g. Loki, Datadog, CloudWatch Logs).
//
// Usage:
//
//	w := output.NewWriter(os.Stdout, output.FormatJSON)
//	w.Write(output.Entry{
//	    Timestamp: time.Now(),
//	    Level:     "INFO",
//	    Message:   "backup completed",
//	    JobName:   "nightly-backup",
//	})
//
// The desired format can be parsed from an environment variable or CLI flag
// using [ParseFormat].
package output
