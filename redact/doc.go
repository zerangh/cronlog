// Package redact provides a Redactor type that prevents sensitive
// values from appearing in structured log entries or webhook payloads.
//
// # Key-based redaction
//
// The Redactor inspects field names against a set of regular expressions.
// Any field whose key matches a pattern (e.g. "password", "token",
// "api_key") has its value replaced with the literal string [REDACTED].
//
// # Line-based redaction
//
// For free-form log lines, the Redactor accepts a slice of known secret
// strings and performs a simple string replacement so that raw secrets
// are never persisted or transmitted.
//
// # Usage
//
//	r := redact.New()
//	safeFields := r.Map(rawFields)
//	safeLine  := r.Line(rawLine, []string{os.Getenv("DB_PASSWORD")})
package redact
