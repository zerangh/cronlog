// Package fanout provides a log entry dispatcher that writes to multiple
// destinations simultaneously, collecting any errors that occur.
package fanout

import (
	"errors"
	"fmt"
	"strings"
)

// Writer is any destination that can receive a log entry represented as a
// byte slice (e.g. an io.Writer adapter, a webhook sink, etc.).
type Writer interface {
	Write(p []byte) (n int, err error)
}

// Fanout dispatches each write to all registered writers.
type Fanout struct {
	writers []Writer
}

// New returns a Fanout that dispatches to the supplied writers.
// At least one writer must be provided.
func New(writers ...Writer) (*Fanout, error) {
	if len(writers) == 0 {
		return nil, errors.New("fanout: at least one writer is required")
	}
	return &Fanout{writers: writers}, nil
}

// Write sends p to every registered writer. All writers are attempted
// regardless of individual failures. If one or more writers fail, a combined
// error is returned; the number of bytes reported is from the first
// successful write, or 0 if all fail.
func (f *Fanout) Write(p []byte) (int, error) {
	var errs []string
	n := 0
	firstSuccess := false

	for i, w := range f.writers {
		wrote, err := w.Write(p)
		if err != nil {
			errs = append(errs, fmt.Sprintf("writer[%d]: %s", i, err.Error()))
			continue
		}
		if !firstSuccess {
			n = wrote
			firstSuccess = true
		}
	}

	if len(errs) > 0 {
		return n, errors.New(strings.Join(errs, "; "))
	}
	return n, nil
}

// Len returns the number of registered writers.
func (f *Fanout) Len() int {
	return len(f.writers)
}
