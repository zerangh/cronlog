package buffer_test

import (
	"testing"

	"github.com/example/cronlog/buffer"
)

func TestNew_DefaultCapacity(t *testing.T) {
	b := buffer.New(0)
	if b == nil {
		t.Fatal("expected non-nil buffer")
	}
}

func TestAdd_StoresEntry(t *testing.T) {
	b := buffer.New(10)
	b.Add(buffer.Entry{Level: "info", Message: "hello"})

	if b.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", b.Len())
	}

	entries := b.Entries()
	if entries[0].Message != "hello" {
		t.Errorf("expected message 'hello', got %q", entries[0].Message)
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	b := buffer.New(3)
	b.Add(buffer.Entry{Message: "first"})
	b.Add(buffer.Entry{Message: "second"})
	b.Add(buffer.Entry{Message: "third"})
	b.Add(buffer.Entry{Message: "fourth"})

	if b.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", b.Len())
	}

	entries := b.Entries()
	if entries[0].Message != "second" {
		t.Errorf("expected oldest entry to be 'second', got %q", entries[0].Message)
	}
	if entries[2].Message != "fourth" {
		t.Errorf("expected newest entry to be 'fourth', got %q", entries[2].Message)
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	b := buffer.New(5)
	b.Add(buffer.Entry{Message: "original"})

	entries := b.Entries()
	entries[0].Message = "mutated"

	got := b.Entries()
	if got[0].Message != "original" {
		t.Errorf("buffer was mutated via returned slice; got %q", got[0].Message)
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	b := buffer.New(5)
	b.Add(buffer.Entry{Message: "a"})
	b.Add(buffer.Entry{Message: "b"})
	b.Reset()

	if b.Len() != 0 {
		t.Errorf("expected 0 entries after reset, got %d", b.Len())
	}
}

func TestAdd_WithFields(t *testing.T) {
	b := buffer.New(5)
	b.Add(buffer.Entry{
		Level:   "error",
		Message: "something failed",
		Fields:  map[string]any{"code": 500},
	})

	entries := b.Entries()
	if entries[0].Fields["code"] != 500 {
		t.Errorf("expected field code=500, got %v", entries[0].Fields["code"])
	}
}
