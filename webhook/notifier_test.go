package webhook_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/cronlog/webhook"
)

func TestNotify_Success(t *testing.T) {
	var received webhook.Payload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := webhook.NewNotifier(server.URL)
	payload := webhook.Payload{
		JobName:   "backup",
		Status:    "failed",
		Message:   "exit status 1",
		ExitCode:  1,
		Timestamp: time.Now().UTC(),
		Duration:  "2m30s",
	}

	if err := n.Notify(payload); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if received.JobName != "backup" {
		t.Errorf("expected job_name 'backup', got '%s'", received.JobName)
	}
	if received.ExitCode != 1 {
		t.Errorf("expected exit_code 1, got %d", received.ExitCode)
	}
}

func TestNotify_NonSuccessStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := webhook.NewNotifier(server.URL)
	err := n.Notify(webhook.Payload{JobName: "test"})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestNotify_EmptyURL(t *testing.T) {
	n := webhook.NewNotifier("")
	err := n.Notify(webhook.Payload{JobName: "test"})
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}
