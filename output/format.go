package output

import (
	"fmt"
	"strings"
)

// ParseFormat parses a string into a Format, returning an error for unknown values.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "json":
		return FormatJSON, nil
	case "text", "":
		return FormatText, nil
	default:
		return "", fmt.Errorf("output: unknown format %q, must be \"json\" or \"text\"", s)
	}
}

// String returns the string representation of a Format.
func (f Format) String() string {
	return string(f)
}

// IsValid reports whether the Format is a recognised value.
func (f Format) IsValid() bool {
	switch f {
	case FormatJSON, FormatText:
		return true
	default:
		return false
	}
}
