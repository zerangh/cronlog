package redact_test

import (
	"regexp"
	"testing"

	"github.com/example/cronlog/redact"
)

func TestField_SafeKey(t *testing.T) {
	r := redact.New()
	got := r.Field("job_name", "backup")
	if got != "backup" {
		t.Errorf("expected 'backup', got %q", got)
	}
}

func TestField_SensitiveKey(t *testing.T) {
	r := redact.New()
	for _, key := range []string{"password", "PASSWORD", "api_key", "auth_token", "secret"} {
		got := r.Field(key, "supersecret")
		if got != "[REDACTED]" {
			t.Errorf("key %q: expected [REDACTED], got %q", key, got)
		}
	}
}

func TestField_CustomPattern(t *testing.T) {
	extra := regexp.MustCompile(`(?i)ssn`)
	r := redact.New(extra)
	got := r.Field("ssn", "123-45-6789")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestMap_RedactsSensitiveFields(t *testing.T) {
	r := redact.New()
	input := map[string]string{
		"job":      "nightly-backup",
		"password": "hunter2",
		"status":   "ok",
		"token":    "abc123",
	}
	out := r.Map(input)

	if out["job"] != "nightly-backup" {
		t.Errorf("job should be unchanged, got %q", out["job"])
	}
	if out["status"] != "ok" {
		t.Errorf("status should be unchanged, got %q", out["status"])
	}
	if out["password"] != "[REDACTED]" {
		t.Errorf("password should be redacted, got %q", out["password"])
	}
	if out["token"] != "[REDACTED]" {
		t.Errorf("token should be redacted, got %q", out["token"])
	}
}

func TestMap_DoesNotMutateInput(t *testing.T) {
	r := redact.New()
	input := map[string]string{"password": "secret"}
	_ = r.Map(input)
	if input["password"] != "secret" {
		t.Error("Map must not mutate the input map")
	}
}

func TestLine_ReplacesSecrets(t *testing.T) {
	r := redact.New()
	line := r.Line("connecting with token abc123 to host db", []string{"abc123"})
	expected := "connecting with token [REDACTED] to host db"
	if line != expected {
		t.Errorf("expected %q, got %q", expected, line)
	}
}

func TestLine_EmptySecretIgnored(t *testing.T) {
	r := redact.New()
	original := "no secrets here"
	got := r.Line(original, []string{""})
	if got != original {
		t.Errorf("expected line unchanged, got %q", got)
	}
}
