// Package filter provides log entry filtering based on severity level.
package filter

import "strings"

// Level represents a log severity level.
type Level int

const (
	// LevelDebug includes all log entries.
	LevelDebug Level = iota
	// LevelInfo includes info, warn, and error entries.
	LevelInfo
	// LevelWarn includes warn and error entries.
	LevelWarn
	// LevelError includes only error entries.
	LevelError
)

// levelNames maps string representations to Level values.
var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
}

// ParseLevel parses a string into a Level.
// Returns LevelInfo and an error if the string is unrecognized.
func ParseLevel(s string) (Level, error) {
	if s == "" {
		return LevelInfo, nil
	}
	l, ok := levelNames[strings.ToLower(strings.TrimSpace(s))]
	if !ok {
		return LevelInfo, &UnknownLevelError{Input: s}
	}
	return l, nil
}

// String returns the string representation of a Level.
func (l Level) String() string {
	for name, level := range levelNames {
		if level == l {
			return name
		}
	}
	return "info"
}

// Filter decides whether a log entry at the given entry level
// should be included given the configured minimum level.
type Filter struct {
	minLevel Level
}

// New creates a new Filter with the given minimum level.
func New(min Level) *Filter {
	return &Filter{minLevel: min}
}

// Allow returns true if the entry level meets or exceeds the minimum level.
func (f *Filter) Allow(entryLevel Level) bool {
	return entryLevel >= f.minLevel
}

// MinLevel returns the configured minimum level.
func (f *Filter) MinLevel() Level {
	return f.minLevel
}
