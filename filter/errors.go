package filter

import "fmt"

// UnknownLevelError is returned when a level string cannot be parsed.
type UnknownLevelError struct {
	Input string
}

func (e *UnknownLevelError) Error() string {
	return fmt.Sprintf("filter: unknown log level %q: must be one of debug, info, warn, error", e.Input)
}
