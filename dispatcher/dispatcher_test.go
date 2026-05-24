package dispatcher_test

import (
	"bytes"
	"strings"
	"testing"

	"cronlog/dispatcher"
)

func TestNew_NoRoutes(t *testing.T) {
	_, err := dispatcher.New()
	if err != dispatcher.ErrNoRoutes {
		t.Fatalf("expected ErrNoRoutes, got %v", err)
	}
}

func TestNew_DuplicateRoute(t *testing.T) {
	w := &bytes.Buffer{}
	r1 := dispatcher.NewRoute("sink", w, "info")
	r2 := dispatcher.NewRoute("sink", w, "warn")
	_, err := dispatcher.New(r1, r2)
	if err != dispatcher.ErrDuplicateRoute {
		t.Fatalf("expected ErrDuplicateRoute, got %v", err)
	}
}

func TestNew_ValidRoutes(t *testing.T) {
	w := &bytes.Buffer{}
	d, err := dispatcher.New(dispatcher.NewRoute("a", w, "info"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Dispatcher")
	}
}

func TestDispatch_WritesToMatchingRoute(t *testing.T) {
	buf := &bytes.Buffer{}
	d, _ := dispatcher.New(dispatcher.NewRoute("out", buf, "info"))
	d.Dispatch(dispatcher.Entry{Level: "info", Message: "hello"})
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected 'hello' in output, got %q", buf.String())
	}
}

func TestDispatch_FiltersLowLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	d, _ := dispatcher.New(dispatcher.NewRoute("out", buf, "error"))
	d.Dispatch(dispatcher.Entry{Level: "debug", Message: "noisy"})
	if buf.Len() != 0 {
		t.Errorf("expected no output for debug entry below error threshold")
	}
}

func TestDispatch_MultipleRoutes(t *testing.T) {
	var a, b bytes.Buffer
	d, _ := dispatcher.New(
		dispatcher.NewRoute("all", &a, "debug"),
		dispatcher.NewRoute("errors", &b, "error"),
	)
	d.Dispatch(dispatcher.Entry{Level: "info", Message: "msg"})
	if a.Len() == 0 {
		t.Error("expected 'all' route to receive entry")
	}
	if b.Len() != 0 {
		t.Error("expected 'errors' route to skip info entry")
	}
}

func TestRoute_Validate_EmptyName(t *testing.T) {
	w := &bytes.Buffer{}
	r := dispatcher.NewRoute("", w, "info")
	if err := r.Validate(); err != dispatcher.ErrEmptyRouteName {
		t.Fatalf("expected ErrEmptyRouteName, got %v", err)
	}
}

func TestRoute_Validate_NilWriter(t *testing.T) {
	r := dispatcher.NewRoute("sink", nil, "info")
	if err := r.Validate(); err != dispatcher.ErrNilWriter {
		t.Fatalf("expected ErrNilWriter, got %v", err)
	}
}
