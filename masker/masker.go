// Package masker provides utilities for masking sensitive substrings
// within log messages before they are written to any output.
package masker

import (
	"strings"
)

// Masker replaces occurrences of registered sensitive values with a
// fixed placeholder so they are never written to log output.
type Masker struct {
	placeholder string
	targets     []string
}

// New creates a Masker that replaces any of the provided target strings
// with placeholder. If placeholder is empty it defaults to "[REDACTED]".
func New(placeholder string, targets ...string) (*Masker, error) {
	if len(targets) == 0 {
		return nil, ErrNoTargets
	}
	for _, t := range targets {
		if t == "" {
			return nil, ErrEmptyTarget
		}
	}
	if placeholder == "" {
		placeholder = "[REDACTED]"
	}
	copy := make([]string, len(targets))
	for i, t := range targets {
		copy[i] = t
	}
	return &Masker{placeholder: placeholder, targets: copy}, nil
}

// Mask replaces every occurrence of each registered target within s and
// returns the sanitised string.
func (m *Masker) Mask(s string) string {
	for _, t := range m.targets {
		s = strings.ReplaceAll(s, t, m.placeholder)
	}
	return s
}

// MaskMap returns a new map where every string value has been passed
// through Mask. The original map is never mutated.
func (m *Masker) MaskMap(fields map[string]string) map[string]string {
	out := make(map[string]string, len(fields))
	for k, v := range fields {
		out[k] = m.Mask(v)
	}
	return out
}

// Targets returns a copy of the registered target strings.
func (m *Masker) Targets() []string {
	copy := make([]string, len(m.targets))
	for i, t := range m.targets {
		copy[i] = t
	}
	return copy
}
