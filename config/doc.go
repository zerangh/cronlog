// Package config provides configuration loading for cronlog.
//
// Configuration can be constructed programmatically via the Config struct or
// loaded automatically from environment variables using FromEnv.
//
// Environment variables:
//
//	CRONLOG_JOB_NAME        — (required) human-readable name for the cron job
//	CRONLOG_WEBHOOK_URL     — (optional) webhook endpoint for failure notifications
//	CRONLOG_TIMEOUT_SECONDS — (optional) job execution timeout in seconds (0 = no timeout)
//	CRONLOG_NOTIFY_SUCCESS  — (optional) send webhook even on successful runs (default: false)
//
// Example usage:
//
//	cfg, err := config.FromEnv()
//	if err != nil {
//	    log.Fatalf("cronlog: bad config: %v", err)
//	}
//	runner.Run(cfg, myJobFunc)
package config
