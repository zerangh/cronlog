package pipeline_test

import (
	"strings"
	"testing"

	"github.com/cronlog/pipeline"
)

func identity(e pipeline.Entry) (pipeline.Entry, bool) { return e, true }

func uppercaseMsg(e pipeline.Entry) (pipeline.Entry, bool) {
	e.Message = strings.ToUpper(e.Message)
	return e, true
}

func dropErrors(e pipeline.Entry) (pipeline.Entry, bool) {
	if e.Level == "error" {
		return e, false
	}
	return e, true
}

func TestNew_NoProcessors(t *testing.T) {
	_, err := pipeline.New()
	if err == nil {
		t.Fatal("expected error for empty processor list")
	}
}

func TestNew_ValidProcessors(t *testing.T) {
	p, err := pipeline.New(identity)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 1 {
		t.Fatalf("expected 1 processor, got %d", p.Len())
	}
}

func TestRun_PassesThroughAllProcessors(t *testing.T) {
	p, _ := pipeline.New(uppercaseMsg)
	entry := pipeline.Entry{Level: "info", Message: "hello"}
	out, ok := p.Run(entry)
	if !ok {
		t.Fatal("expected entry to pass through")
	}
	if out.Message != "HELLO" {
		t.Fatalf("expected HELLO, got %s", out.Message)
	}
}

func TestRun_DropsEntryWhenProcessorReturnsFalse(t *testing.T) {
	p, _ := pipeline.New(dropErrors)
	entry := pipeline.Entry{Level: "error", Message: "boom"}
	_, ok := p.Run(entry)
	if ok {
		t.Fatal("expected entry to be dropped")
	}
}

func TestRun_StopsAtFirstDrop(t *testing.T) {
	called := false
	marker := func(e pipeline.Entry) (pipeline.Entry, bool) {
		called = true
		return e, true
	}
	p, _ := pipeline.New(dropErrors, marker)
	p.Run(pipeline.Entry{Level: "error", Message: "boom"})
	if called {
		t.Fatal("subsequent processor should not have been called")
	}
}

func TestRun_ChainedTransformations(t *testing.T) {
	addField := func(e pipeline.Entry) (pipeline.Entry, bool) {
		if e.Fields == nil {
			e.Fields = map[string]string{}
		}
		e.Fields["env"] = "test"
		return e, true
	}
	p, _ := pipeline.New(uppercaseMsg, addField)
	out, ok := p.Run(pipeline.Entry{Level: "info", Message: "hi"})
	if !ok {
		t.Fatal("expected entry to pass")
	}
	if out.Message != "HI" {
		t.Fatalf("expected HI, got %s", out.Message)
	}
	if out.Fields["env"] != "test" {
		t.Fatalf("expected env=test, got %s", out.Fields["env"])
	}
}
