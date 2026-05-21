package runner_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/cronlog/config"
	"github.com/example/cronlog/runner"
)

func newTestConfig(webhookURL string) config.Config {
	return config.Config{
		JobName:    "test-job",
		WebhookURL: webhookURL,
		TimeoutSec: 5,
	}
}

func TestRun_SuccessfulJob(t *testing.T) {
	var received bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := newTestConfig(server.URL)
	err := runner.Run(cfg, func() error {
		return nil
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Webhook should not fire on success with no errors logged
	if received {
		t.Error("expected webhook not to be called on clean success")
	}
}

func TestRun_JobTimeout(t *testing.T) {
	cfg := newTestConfig("")
	cfg.TimeoutSec = 1

	start := time.Now()
	err := runner.Run(cfg, func() error {
		time.Sleep(3 * time.Second)
		return nil
	})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if elapsed > 2*time.Second {
		t.Errorf("job did not time out in time, elapsed: %v", elapsed)
	}
}

func TestRun_WebhookFiredOnError(t *testing.T) {
	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := newTestConfig(server.URL)
	_ = runner.Run(cfg, func() error {
		return fmt.Errorf("something went wrong")
	})

	if !called {
		t.Error("expected webhook to be called when job returns error")
	}
}
