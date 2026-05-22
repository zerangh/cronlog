package context_test

import (
	"context"
	"testing"
	"time"

	jobctx "github.com/cronlog/context"
)

func fixedTime() time.Time {
	return time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
}

func baseMeta() jobctx.JobMeta {
	return jobctx.JobMeta{
		JobName:   "backup",
		RunID:     "abc-123",
		Attempt:   1,
		StartedAt: fixedTime(),
	}
}

func TestWithJob_StoreAndRetrieve(t *testing.T) {
	ctx := jobctx.WithJob(context.Background(), baseMeta())
	meta, ok := jobctx.JobFrom(ctx)
	if !ok {
		t.Fatal("expected meta to be present")
	}
	if meta.JobName != "backup" {
		t.Errorf("job name: got %q, want %q", meta.JobName, "backup")
	}
	if meta.RunID != "abc-123" {
		t.Errorf("run id: got %q, want %q", meta.RunID, "abc-123")
	}
	if meta.Attempt != 1 {
		t.Errorf("attempt: got %d, want 1", meta.Attempt)
	}
}

func TestJobFrom_MissingReturnsNotOK(t *testing.T) {
	_, ok := jobctx.JobFrom(context.Background())
	if ok {
		t.Fatal("expected ok=false for empty context")
	}
}

func TestFields_ContainsAllKeys(t *testing.T) {
	ctx := jobctx.WithJob(context.Background(), baseMeta())
	fields := jobctx.Fields(ctx)

	expected := []string{"job_name", "run_id", "attempt", "started_at"}
	for _, k := range expected {
		if _, ok := fields[k]; !ok {
			t.Errorf("missing field %q in Fields output", k)
		}
	}
}

func TestFields_EmptyContextReturnsEmptyMap(t *testing.T) {
	fields := jobctx.Fields(context.Background())
	if len(fields) != 0 {
		t.Errorf("expected empty map, got %v", fields)
	}
}

func TestFields_AttemptIsStringified(t *testing.T) {
	meta := baseMeta()
	meta.Attempt = 3
	ctx := jobctx.WithJob(context.Background(), meta)
	fields := jobctx.Fields(ctx)
	if fields["attempt"] != "3" {
		t.Errorf("attempt field: got %q, want %q", fields["attempt"], "3")
	}
}

func TestFields_StartedAtIsRFC3339(t *testing.T) {
	ctx := jobctx.WithJob(context.Background(), baseMeta())
	fields := jobctx.Fields(ctx)
	want := "2024-01-15T10:00:00Z"
	if fields["started_at"] != want {
		t.Errorf("started_at: got %q, want %q", fields["started_at"], want)
	}
}
