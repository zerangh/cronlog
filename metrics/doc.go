// Package metrics provides lightweight job execution metrics collection
// for cronlog.
//
// A Collector is created at the start of a job run, records timing
// information, and produces a Result once the job finishes. The Result
// captures the job name, wall-clock start and finish times, total
// duration, exit code, whether any errors were logged, and the total
// number of log lines emitted.
//
// Typical usage:
//
//	c := metrics.NewCollector(cfg.JobName)
//	c.Start()
//	// ... run job ...
//	result := c.Finish(exitCode, logger.HasErrors(), logger.LineCount())
//
// The Result can be serialised to JSON and included in webhook
// notification payloads or written to structured log output.
package metrics
