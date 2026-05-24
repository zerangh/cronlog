// Package tee provides a log entry writer that duplicates writes to
// multiple destinations, similar to the Unix tee command.
package tee

import (
	"errors"
	"fmt"
	"io"
)

// Writer duplicates each Write call to all registered destinations.
// If any destination returns an error the write continues to the
// remaining destinations; all errors are joined and returned.
type Writer struct {
	dsts []io.Writer
}

// New creates a Writer that fans out to each of the provided
// destinations. At least one destination must be supplied.
func New(dsts ...io.Writer) (*Writer, error) {
	if len(dsts) == 0 {
		return nil, errors.New("tee: at least one destination is required")
	}
	for i, d := range dsts {
		if d == nil {
			return nil, fmt.Errorf("tee: destination at index %d is nil", i)
		}
	}
	return &Writer{dsts: dsts}, nil
}

// Write writes p to every destination. It returns the length of p on
// success. If one or more destinations fail their errors are joined
// and returned alongside the number of bytes written by the first
// successful destination (or 0 if all fail).
func (w *Writer) Write(p []byte) (int, error) {
	var errs []error
	n := 0
	for _, d := range w.dsts {
		wrote, err := d.Write(p)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if n == 0 {
			n = wrote
		}
	}
	if len(errs) > 0 {
		return n, errors.Join(errs...)
	}
	return n, nil
}

// Len returns the number of destinations registered with the writer.
func (w *Writer) Len() int { return len(w.dsts) }
