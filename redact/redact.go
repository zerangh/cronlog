// Package redact provides utilities for scrubbing sensitive values
// (e.g. passwords, tokens, secrets) from log output before it is
// written or transmitted.
package redact

import (
	"regexp"
	"strings"
)

const placeholder = "[REDACTED]"

// defaultPatterns matches common secret-like key names.
var defaultPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|secret|token|api[_-]?key|auth)`),
}

// Redactor scrubs sensitive fields from structured log messages.
type Redactor struct {
	patterns []*regexp.Regexp
}

// New returns a Redactor that uses the built-in sensitive-key patterns
// plus any additional patterns supplied by the caller.
func New(extra ...*regexp.Regexp) *Redactor {
	patterns := make([]*regexp.Regexp, len(defaultPatterns)+len(extra))
	copy(patterns, defaultPatterns)
	copy(patterns[len(defaultPatterns):], extra)
	return &Redactor{patterns: patterns}
}

// Field returns the value unchanged when the key is considered safe,
// or the placeholder string when the key matches a sensitive pattern.
func (r *Redactor) Field(key, value string) string {
	for _, p := range r.patterns {
		if p.MatchString(key) {
			return placeholder
		}
	}
	return value
}

// Map returns a copy of fields with sensitive values replaced.
func (r *Redactor) Map(fields map[string]string) map[string]string {
	out := make(map[string]string, len(fields))
	for k, v := range fields {
		out[k] = r.Field(k, v)
	}
	return out
}

// Line replaces occurrences of any supplied literal secrets inside a
// free-form log line with the placeholder.
func (r *Redactor) Line(line string, secrets []string) string {
	for _, s := range secrets {
		if s == "" {
			continue
		}
		line = strings.ReplaceAll(line, s, placeholder)
	}
	return line
}
