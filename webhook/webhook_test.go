package webhook_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/cronlog/webhook"
)

func TestSend_Success(t *testing.T) {
	var received webhook.Payload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := webhook.NewClient(server.URL)
	now := time.Now()
	p := webhook.Payload{
		JobName:    "backup",
		Status:     "failure",
		Message:    "exit status 1",
		ExitCode:   1,
		StartedAt:  now,
		FinishedAt: now.Add(5 * time.Second),
		Duration:   "5s",
	}

	if err := client.Send(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.JobName != "backup" {
		t.Errorf("expected job_name 'backup', got '%s'", received.JobName)
	}
	if received.Status != "failure" {
		t.Errorf("expected status 'failure', got '%s'", received.Status)
	}
}

func TestSend_NonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := webhook.NewClient(server.URL)
	err := client.Send(webhook.Payload{JobName: "test"})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSend_EmptyURL(t *testing.T) {
	client := webhook.NewClient("")
	if err := client.Send(webhook.Payload{JobName: "test"}); err != nil {
		t.Fatalf("expected no error for empty URL, got: %v", err)
	}
}
