package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the webhook notification body sent on cron job failure.
type Payload struct {
	JobName   string    `json:"job_name"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	ExitCode  int       `json:"exit_code"`
	Timestamp time.Time `json:"timestamp"`
	Duration  string    `json:"duration"`
}

// Notifier sends failure notifications to a configured webhook URL.
type Notifier struct {
	URL        string
	HTTPClient *http.Client
}

// NewNotifier creates a Notifier with a default HTTP client timeout.
func NewNotifier(url string) *Notifier {
	return &Notifier{
		URL: url,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Notify sends a POST request with the given payload to the webhook URL.
func (n *Notifier) Notify(p Payload) error {
	if n.URL == "" {
		return fmt.Errorf("webhook: URL is not configured")
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
