package tee_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/cronlog/tee"
)

// errWriter always returns an error on Write.
type errWriter struct{ err error }

func (e *errWriter) Write(_ []byte) (int, error) { return 0, e.err }

func TestNew_NoDestinations(t *testing.T) {
	_, err := tee.New()
	if err == nil {
		t.Fatal("expected error for zero destinations, got nil")
	}
}

func TestNew_NilDestination(t *testing.T) {
	_, err := tee.New(nil)
	if err == nil {
		t.Fatal("expected error for nil destination, got nil")
	}
}

func TestNew_ValidDestinations(t *testing.T) {
	var a, b bytes.Buffer
	w, err := tee.New(&a, &b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", w.Len())
	}
}

func TestWrite_DuplicatesToAllDestinations(t *testing.T) {
	var a, b bytes.Buffer
	w, _ := tee.New(&a, &b)

	msg := []byte("hello tee")
	n, err := w.Write(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(msg) {
		t.Fatalf("expected n=%d, got %d", len(msg), n)
	}
	if a.String() != string(msg) {
		t.Errorf("destination a: got %q, want %q", a.String(), string(msg))
	}
	if b.String() != string(msg) {
		t.Errorf("destination b: got %q, want %q", b.String(), string(msg))
	}
}

func TestWrite_ContinuesAfterPartialFailure(t *testing.T) {
	var good bytes.Buffer
	bad := &errWriter{err: errors.New("disk full")}

	w, _ := tee.New(bad, &good)

	msg := []byte("partial")
	_, err := w.Write(msg)
	if err == nil {
		t.Fatal("expected error from failing destination, got nil")
	}
	if good.String() != string(msg) {
		t.Errorf("healthy destination did not receive write: got %q", good.String())
	}
}

func TestWrite_AllFailures_ReturnsJoinedError(t *testing.T) {
	e1 := &errWriter{err: errors.New("err1")}
	e2 := &errWriter{err: errors.New("err2")}

	w, _ := tee.New(e1, e2)
	_, err := w.Write([]byte("data"))
	if err == nil {
		t.Fatal("expected joined error, got nil")
	}
}

func TestWrite_DiscardsToDevNull(t *testing.T) {
	w, _ := tee.New(io.Discard)
	_, err := w.Write([]byte("silent"))
	if err != nil {
		t.Fatalf("unexpected error writing to discard: %v", err)
	}
}
