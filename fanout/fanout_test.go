package fanout_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/cronlog/fanout"
)

// stubWriter records what was written to it and can be configured to fail.
type stubWriter struct {
	recorded []byte
	failWith error
}

func (s *stubWriter) Write(p []byte) (int, error) {
	if s.failWith != nil {
		return 0, s.failWith
	}
	s.recorded = append(s.recorded, p...)
	return len(p), nil
}

func TestNew_NoWriters(t *testing.T) {
	_, err := fanout.New()
	if err == nil {
		t.Fatal("expected error when no writers provided")
	}
}

func TestNew_ValidWriters(t *testing.T) {
	f, err := fanout.New(&stubWriter{}, &stubWriter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Len() != 2 {
		t.Fatalf("expected 2 writers, got %d", f.Len())
	}
}

func TestWrite_DispatchesToAllWriters(t *testing.T) {
	a, b := &stubWriter{}, &stubWriter{}
	f, _ := fanout.New(a, b)

	msg := []byte("hello fanout")
	n, err := f.Write(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(msg) {
		t.Fatalf("expected n=%d, got %d", len(msg), n)
	}
	if string(a.recorded) != string(msg) {
		t.Errorf("writer a: got %q, want %q", a.recorded, msg)
	}
	if string(b.recorded) != string(msg) {
		t.Errorf("writer b: got %q, want %q", b.recorded, msg)
	}
}

func TestWrite_ContinuesAfterPartialFailure(t *testing.T) {
	good := &stubWriter{}
	bad := &stubWriter{failWith: errors.New("disk full")}
	f, _ := fanout.New(bad, good)

	_, err := f.Write([]byte("data"))
	if err == nil {
		t.Fatal("expected error from failing writer")
	}
	if !strings.Contains(err.Error(), "writer[0]") {
		t.Errorf("error should identify the failing writer index, got: %v", err)
	}
	if string(good.recorded) != "data" {
		t.Errorf("good writer should still receive data, got %q", good.recorded)
	}
}

func TestWrite_AllFail(t *testing.T) {
	a := &stubWriter{failWith: errors.New("err a")}
	b := &stubWriter{failWith: errors.New("err b")}
	f, _ := fanout.New(a, b)

	n, err := f.Write([]byte("x"))
	if err == nil {
		t.Fatal("expected error when all writers fail")
	}
	if n != 0 {
		t.Errorf("expected n=0 when all writers fail, got %d", n)
	}
	if !strings.Contains(err.Error(), "writer[0]") || !strings.Contains(err.Error(), "writer[1]") {
		t.Errorf("error should mention both failing writers: %v", err)
	}
}
