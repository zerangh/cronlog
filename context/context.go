// Package context provides job-scoped context enrichment for cronlog.
// It attaches structured metadata (job name, run ID, attempt number) to
// a standard context.Context so all downstream components share a
// consistent set of fields without explicit parameter threading.
package context

import (
	"context"
	"fmt"
	"time"
)

type key int

const jobKey key = iota

// JobMeta holds metadata about the current cron job execution.
type JobMeta struct {
	JobName   string
	RunID     string
	Attempt   int
	StartedAt time.Time
}

// WithJob returns a new context carrying the provided JobMeta.
func WithJob(ctx context.Context, meta JobMeta) context.Context {
	return context.WithValue(ctx, jobKey, meta)
}

// JobFrom retrieves the JobMeta stored in ctx.
// The second return value reports whether a value was present.
func JobFrom(ctx context.Context) (JobMeta, bool) {
	v, ok := ctx.Value(jobKey).(JobMeta)
	return v, ok
}

// Fields returns a map of log-ready key/value pairs derived from the
// JobMeta stored in ctx. If no meta is present the map is empty.
func Fields(ctx context.Context) map[string]string {
	meta, ok := JobFrom(ctx)
	if !ok {
		return map[string]string{}
	}
	return map[string]string{
		"job_name":   meta.JobName,
		"run_id":     meta.RunID,
		"attempt":    fmt.Sprintf("%d", meta.Attempt),
		"started_at": meta.StartedAt.UTC().Format(time.RFC3339),
	}
}
