// Package metrics provides simple job execution metrics collection.
package metrics

import (
	"time"
)

// Result holds the collected metrics for a single job run.
type Result struct {
	JobName   string        `json:"job_name"`
	StartedAt time.Time     `json:"started_at"`
	FinishedAt time.Time    `json:"finished_at"`
	Duration  time.Duration `json:"duration_ms"`
	ExitCode  int           `json:"exit_code"`
	HasErrors bool          `json:"has_errors"`
	LogLines  int           `json:"log_lines"`
}

// Collector gathers timing and outcome data for a job run.
type Collector struct {
	jobName   string
	startedAt time.Time
	now       func() time.Time
}

// NewCollector creates a new Collector for the given job name.
func NewCollector(jobName string) *Collector {
	return &Collector{
		jobName: jobName,
		now:     time.Now,
	}
}

// Start records the job start time.
func (c *Collector) Start() {
	c.startedAt = c.now()
}

// Finish builds and returns a Result using the provided outcome data.
func (c *Collector) Finish(exitCode int, hasErrors bool, logLines int) Result {
	finished := c.now()
	return Result{
		JobName:    c.jobName,
		StartedAt:  c.startedAt,
		FinishedAt: finished,
		Duration:   finished.Sub(c.startedAt),
		ExitCode:   exitCode,
		HasErrors:  hasErrors,
		LogLines:   logLines,
	}
}
