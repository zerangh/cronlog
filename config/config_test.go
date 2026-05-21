package config_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronlog/config"
)

func TestValidate_MissingJobName(t *testing.T) {
	cfg := &config.Config{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing job name, got nil")
	}
}

func TestValidate_NegativeTimeout(t *testing.T) {
	cfg := &config.Config{JobName: "test-job", TimeoutSeconds: -1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative timeout, got nil")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := &config.Config{JobName: "test-job", TimeoutSeconds: 30}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestTimeout_ReturnsDuration(t *testing.T) {
	cfg := &config.Config{JobName: "j", TimeoutSeconds: 45}
	if got := cfg.Timeout(); got != 45*time.Second {
		t.Fatalf("expected 45s, got %v", got)
	}
}

func TestFromEnv_MissingJobName(t *testing.T) {
	t.Setenv("CRONLOG_JOB_NAME", "")
	_, err := config.FromEnv()
	if err == nil {
		t.Fatal("expected error when CRONLOG_JOB_NAME is empty")
	}
}

func TestFromEnv_InvalidTimeout(t *testing.T) {
	t.Setenv("CRONLOG_JOB_NAME", "my-job")
	t.Setenv("CRONLOG_TIMEOUT_SECONDS", "not-a-number")
	t.Cleanup(func() { t.Setenv("CRONLOG_TIMEOUT_SECONDS", "") })
	_, err := config.FromEnv()
	if err == nil {
		t.Fatal("expected error for invalid timeout value")
	}
}

func TestFromEnv_ValidDefaults(t *testing.T) {
	t.Setenv("CRONLOG_JOB_NAME", "nightly-backup")
	t.Setenv("CRONLOG_WEBHOOK_URL", "https://hooks.example.com/abc")
	t.Setenv("CRONLOG_TIMEOUT_SECONDS", "")
	t.Setenv("CRONLOG_NOTIFY_SUCCESS", "")

	cfg, err := config.FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.JobName != "nightly-backup" {
		t.Errorf("expected job name 'nightly-backup', got %q", cfg.JobName)
	}
	if cfg.TimeoutSeconds != 0 {
		t.Errorf("expected default timeout 0, got %d", cfg.TimeoutSeconds)
	}
	if cfg.NotifyOnSuccess {
		t.Error("expected NotifyOnSuccess to default to false")
	}
}

func TestFromEnv_NotifyOnSuccess(t *testing.T) {
	t.Setenv("CRONLOG_JOB_NAME", "hourly-sync")
	t.Setenv("CRONLOG_NOTIFY_SUCCESS", "true")
	t.Cleanup(func() { t.Setenv("CRONLOG_NOTIFY_SUCCESS", "") })

	cfg, err := config.FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.NotifyOnSuccess {
		t.Error("expected NotifyOnSuccess to be true")
	}
}
