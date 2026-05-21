package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the configuration for a cronlog run.
type Config struct {
	// JobName is a human-readable identifier for the cron job.
	JobName string

	// WebhookURL is the endpoint to POST failure notifications to.
	WebhookURL string

	// TimeoutSeconds is the maximum number of seconds to allow the job to run.
	TimeoutSeconds int

	// NotifyOnSuccess controls whether a webhook call is made even on success.
	NotifyOnSuccess bool
}

// Timeout returns the timeout as a time.Duration.
func (c *Config) Timeout() time.Duration {
	return time.Duration(c.TimeoutSeconds) * time.Second
}

// Validate returns an error if required fields are missing or invalid.
func (c *Config) Validate() error {
	if c.JobName == "" {
		return fmt.Errorf("config: job_name must not be empty")
	}
	if c.TimeoutSeconds < 0 {
		return fmt.Errorf("config: timeout_seconds must be non-negative, got %d", c.TimeoutSeconds)
	}
	return nil
}

// FromEnv builds a Config by reading well-known environment variables.
// CRONLOG_JOB_NAME        — required
// CRONLOG_WEBHOOK_URL     — optional
// CRONLOG_TIMEOUT_SECONDS — optional, defaults to 0 (no timeout)
// CRONLOG_NOTIFY_SUCCESS  — optional, defaults to false
func FromEnv() (*Config, error) {
	cfg := &Config{
		JobName:    os.Getenv("CRONLOG_JOB_NAME"),
		WebhookURL: os.Getenv("CRONLOG_WEBHOOK_URL"),
	}

	if raw := os.Getenv("CRONLOG_TIMEOUT_SECONDS"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("config: invalid CRONLOG_TIMEOUT_SECONDS %q: %w", raw, err)
		}
		cfg.TimeoutSeconds = v
	}

	if raw := os.Getenv("CRONLOG_NOTIFY_SUCCESS"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("config: invalid CRONLOG_NOTIFY_SUCCESS %q: %w", raw, err)
		}
		cfg.NotifyOnSuccess = v
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}
