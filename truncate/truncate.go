// Package truncate provides log field and message truncation to prevent
// excessively large payloads from being sent to webhook notifiers or
// written to output streams.
package truncate

import "fmt"

const (
	DefaultMaxLen = 1024
	Ellipsis      = "..."
)

// Truncator truncates string values that exceed a configured maximum length.
type Truncator struct {
	maxLen int
}

// New returns a Truncator that truncates strings longer than maxLen bytes.
// If maxLen is less than 1, an error is returned.
func New(maxLen int) (*Truncator, error) {
	if maxLen < 1 {
		return nil, fmt.Errorf("truncate: maxLen must be at least 1, got %d", maxLen)
	}
	return &Truncator{maxLen: maxLen}, nil
}

// String truncates s to the configured maximum length. If s exceeds the
// limit, it is cut and an ellipsis is appended. The returned string will
// never exceed maxLen bytes.
func (t *Truncator) String(s string) string {
	if len(s) <= t.maxLen {
		return s
	}
	cutAt := t.maxLen - len(Ellipsis)
	if cutAt < 0 {
		cutAt = 0
	}
	return s[:cutAt] + Ellipsis
}

// Fields applies truncation to every string value in the provided map.
// Non-string values are left unchanged. The original map is not mutated;
// a new map is returned.
func (t *Truncator) Fields(fields map[string]any) map[string]any {
	out := make(map[string]any, len(fields))
	for k, v := range fields {
		if s, ok := v.(string); ok {
			out[k] = t.String(s)
		} else {
			out[k] = v
		}
	}
	return out
}
