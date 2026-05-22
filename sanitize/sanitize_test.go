package sanitize_test

import (
	"strings"
	"testing"

	"github.com/cronlog/sanitize"
)

func TestNew_DefaultMaxRunes(t *testing.T) {
	s := sanitize.New(0)
	long := strings.Repeat("a", 2000)
	got := s.String(long)
	if len([]rune(got)) != 1024 {
		t.Fatalf("expected 1024 runes, got %d", len([]rune(got)))
	}
}

func TestString_NullBytesRemoved(t *testing.T) {
	s := sanitize.New(100)
	got := s.String("hel\x00lo")
	if got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestString_ControlCharsRemoved(t *testing.T) {
	s := sanitize.New(100)
	// \x01 and \x1b (ESC) are control chars; tab and newline should survive
	got := s.String("a\x01b\x1bc\td\ne")
	if got != "abc\td\ne" {
		t.Fatalf("expected %q, got %q", "abc\td\ne", got)
	}
}

func TestString_TruncatesAtMaxRunes(t *testing.T) {
	s := sanitize.New(5)
	got := s.String("abcdefgh")
	if got != "abcde" {
		t.Fatalf("expected %q, got %q", "abcde", got)
	}
}

func TestString_ShortInputUnchanged(t *testing.T) {
	s := sanitize.New(100)
	got := s.String("hello world")
	if got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}
}

func TestString_EmptyInput(t *testing.T) {
	s := sanitize.New(100)
	got := s.String("")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestMap_SanitizesStringValues(t *testing.T) {
	s := sanitize.New(10)
	input := map[string]any{
		"msg":   "hel\x00lo",
		"count": 42,
		"flag":  true,
	}
	out := s.Map(input)

	if out["msg"] != "hello" {
		t.Fatalf("expected %q, got %q", "hello", out["msg"])
	}
	if out["count"] != 42 {
		t.Fatalf("expected 42, got %v", out["count"])
	}
	if out["flag"] != true {
		t.Fatalf("expected true, got %v", out["flag"])
	}
}

func TestMap_DoesNotMutateInput(t *testing.T) {
	s := sanitize.New(100)
	input := map[string]any{"key": "val\x00ue"}
	_ = s.Map(input)
	if input["key"] != "val\x00ue" {
		t.Fatal("input map was mutated")
	}
}
