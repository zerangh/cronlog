// Package sanitize provides utilities for cleaning log field values
// before they are written to output. It removes or replaces control
// characters, null bytes, and excessively long strings that may corrupt
// structured log output or cause downstream parsing issues.
package sanitize

import (
	"strings"
	"unicode"
)

const defaultMaxRunes = 1024

// Sanitizer cleans string values for safe log output.
type Sanitizer struct {
	maxRunes int
}

// New returns a Sanitizer that truncates values to maxRunes runes.
// If maxRunes is <= 0, a default of 1024 is used.
func New(maxRunes int) *Sanitizer {
	if maxRunes <= 0 {
		maxRunes = defaultMaxRunes
	}
	return &Sanitizer{maxRunes: maxRunes}
}

// String sanitizes a single string value. It strips null bytes and
// non-printable control characters (except tab and newline), then
// truncates the result to the configured maximum rune count.
func (s *Sanitizer) String(v string) string {
	var b strings.Builder
	b.Grow(len(v))

	count := 0
	for _, r := range v {
		if count >= s.maxRunes {
			break
		}
		if r == '\x00' {
			continue
		}
		if unicode.IsControl(r) && r != '\t' && r != '\n' {
			continue
		}
		b.WriteRune(r)
		count++
	}

	return b.String()
}

// Map sanitizes all string values in a map, returning a new map.
// Non-string values are passed through unchanged.
func (s *Sanitizer) Map(fields map[string]any) map[string]any {
	out := make(map[string]any, len(fields))
	for k, v := range fields {
		if str, ok := v.(string); ok {
			out[k] = s.String(str)
		} else {
			out[k] = v
		}
	}
	return out
}
