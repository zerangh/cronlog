package truncate_test

import (
	"strings"
	"testing"

	"github.com/example/cronlog/truncate"
)

func TestNew_InvalidMaxLen(t *testing.T) {
	_, err := truncate.New(0)
	if err == nil {
		t.Fatal("expected error for maxLen=0, got nil")
	}
}

func TestNew_ValidMaxLen(t *testing.T) {
	tr, err := truncate.New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestString_ShortInput(t *testing.T) {
	tr, _ := truncate.New(20)
	got := tr.String("hello")
	if got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestString_ExactLength(t *testing.T) {
	tr, _ := truncate.New(5)
	got := tr.String("hello")
	if got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestString_ExceedsLimit(t *testing.T) {
	tr, _ := truncate.New(10)
	input := strings.Repeat("a", 20)
	got := tr.String(input)
	if len(got) > 10 {
		t.Errorf("expected len <= 10, got %d", len(got))
	}
	if !strings.HasSuffix(got, truncate.Ellipsis) {
		t.Errorf("expected ellipsis suffix, got %q", got)
	}
}

func TestString_VerySmallMaxLen(t *testing.T) {
	tr, _ := truncate.New(1)
	got := tr.String("hello")
	if len(got) > 1 {
		t.Errorf("expected len <= 1, got %d", len(got))
	}
}

func TestFields_TruncatesStringValues(t *testing.T) {
	tr, _ := truncate.New(8)
	fields := map[string]any{
		"short": "hi",
		"long":  "this is a very long value",
		"count": 42,
	}
	out := tr.Fields(fields)

	if out["short"] != "hi" {
		t.Errorf("short field modified unexpectedly: %v", out["short"])
	}
	if s, ok := out["long"].(string); !ok || len(s) > 8 {
		t.Errorf("long field not truncated correctly: %v", out["long"])
	}
	if out["count"] != 42 {
		t.Errorf("non-string field mutated: %v", out["count"])
	}
}

func TestFields_DoesNotMutateInput(t *testing.T) {
	tr, _ := truncate.New(5)
	orig := "original long string"
	fields := map[string]any{"msg": orig}
	tr.Fields(fields)
	if fields["msg"] != orig {
		t.Error("Fields mutated the input map")
	}
}
