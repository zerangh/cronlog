// Package buffer provides a thread-safe, capacity-bounded ring buffer for
// storing structured log entries produced during a cron job run.
//
// Typical usage:
//
//	b := buffer.New(200)
//	b.Add(buffer.Entry{Level: "info", Message: "job started"})
//	// ... job runs ...
//	entries := b.Entries() // retrieve all entries for webhook payload
//
// When the buffer reaches its capacity the oldest entry is silently evicted,
// ensuring that memory usage remains bounded even for long-running jobs that
// produce a high volume of log output.
//
// Buffer is safe for concurrent use by multiple goroutines.
package buffer
