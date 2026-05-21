package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the notification payload sent to the webhook endpoint.
type Payload struct {
	JobName   string    `json:"job_name"`
	Status    string    `json:"status"` // "success" or "failure"
	Message   string    `json:"message"`
	ExitCode  int       `json:"exit_code"`
	StartedAt time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Duration  string    `json:"duration"`
}

// Client sends webhook notifications.
type Client struct {
	URL        string
	HTTPClient *http.Client
}

// NewClient creates a new webhook Client with the given URL.
func NewClient(url string) *Client {
	return &Client{
		URL: url,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send marshals the payload and POSTs it to the webhook URL.
func (c *Client) Send(p Payload) error {
	if c.URL == "" {
		return nil
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: failed to marshal payload: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
