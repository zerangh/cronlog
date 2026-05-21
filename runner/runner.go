package runner

import (
	"context"
	"fmt"

	"github.com/yourorg/cronlog/logger"
	"github.com/yourorg/cronlog/webhook"
)

// JobFunc is the signature for a cron job handler.
type JobFunc func(ctx context.Context, log *logger.Logger) error

// Config holds the runtime configuration for a single cron job execution.
type Config struct {
	// JobName is a human-readable identifier included in logs and notifications.
	JobName string
	// WebhookURL is the endpoint to POST failure notifications to.
	// Leave empty to disable notifications.
	WebhookURL string
	// NotifyOnSuccess sends a webhook even when the job succeeds.
	NotifyOnSuccess bool
}

// Run executes fn with a structured logger, then sends a webhook notification
// if the job fails (or if NotifyOnSuccess is set).
func Run(ctx context.Context, cfg Config, fn JobFunc) error {
	log := logger.New(cfg.JobName, nil)

	log.Info(fmt.Sprintf("job %q started", cfg.JobName), nil)

	runErr := fn(ctx, log)

	if runErr != nil {
		log.Error("job finished with error", map[string]string{"error": runErr.Error()})
	} else {
		log.Info(fmt.Sprintf("job %q completed successfully", cfg.JobName), map[string]string{
			"duration": log.Duration().String(),
		})
	}

	// Notify on failure always; notify on success only when NotifyOnSuccess is
	// set and the job produced no logged errors during its run.
	shouldNotify := runErr != nil || (cfg.NotifyOnSuccess && !log.HasErrors())
	if cfg.WebhookURL != "" && shouldNotify {
		if notifyErr := notify(cfg.WebhookURL, cfg.JobName, log, runErr); notifyErr != nil {
			log.Warn("webhook notification failed", map[string]string{"error": notifyErr.Error()})
		}
	}

	return runErr
}

// notify sends a webhook notification for the completed job.
func notify(webhookURL, jobName string, log *logger.Logger, jobErr error) error {
	n := webhook.NewNotifier(webhookURL)
	payload := buildPayload(jobName, log, jobErr)
	return n.Notify(payload)
}

// buildPayload constructs the notification message from the job summary.
func buildPayload(jobName string, log *logger.Logger, jobErr error) string {
	status := "SUCCESS"
	if jobErr != nil {
		status = "FAILURE"
	}
	return fmt.Sprintf("[cronlog] job=%q status=%s\n%s", jobName, status, log.Summary())
}
