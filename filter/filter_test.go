package filter_test

import (
	"testing"

	"github.com/cronlog/filter"
)

func TestParseLevel_Debug(t *testing.T) {
	l, err := filter.ParseLevel("debug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l != filter.LevelDebug {
		t.Errorf("expected LevelDebug, got %v", l)
	}
}

func TestParseLevel_CaseInsensitive(t *testing.T) {
	for _, s := range []string{"WARN", "Warn", "warn"} {
		l, err := filter.ParseLevel(s)
		if err != nil {
			t.Fatalf("unexpected error for input %q: %v", s, err)
		}
		if l != filter.LevelWarn {
			t.Errorf("expected LevelWarn for %q, got %v", s, l)
		}
	}
}

func TestParseLevel_EmptyDefaultsToInfo(t *testing.T) {
	l, err := filter.ParseLevel("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l != filter.LevelInfo {
		t.Errorf("expected LevelInfo, got %v", l)
	}
}

func TestParseLevel_Unknown(t *testing.T) {
	_, err := filter.ParseLevel("verbose")
	if err == nil {
		t.Fatal("expected error for unknown level, got nil")
	}
	var ule *filter.UnknownLevelError
	if !errorAs(err, &ule) {
		t.Errorf("expected UnknownLevelError, got %T", err)
	}
}

func TestFilter_Allow(t *testing.T) {
	tests := []struct {
		min      filter.Level
		entry    filter.Level
		expected bool
	}{
		{filter.LevelInfo, filter.LevelDebug, false},
		{filter.LevelInfo, filter.LevelInfo, true},
		{filter.LevelInfo, filter.LevelWarn, true},
		{filter.LevelInfo, filter.LevelError, true},
		{filter.LevelError, filter.LevelWarn, false},
		{filter.LevelDebug, filter.LevelDebug, true},
	}
	for _, tt := range tests {
		f := filter.New(tt.min)
		got := f.Allow(tt.entry)
		if got != tt.expected {
			t.Errorf("Allow(%v) with min=%v: expected %v, got %v",
				tt.entry, tt.min, tt.expected, got)
		}
	}
}

func TestFilter_MinLevel(t *testing.T) {
	f := filter.New(filter.LevelWarn)
	if f.MinLevel() != filter.LevelWarn {
		t.Errorf("expected LevelWarn, got %v", f.MinLevel())
	}
}

// errorAs is a helper to avoid importing errors in test.
func errorAs(err error, target interface{}) bool {
	if ule, ok := target.(**filter.UnknownLevelError); ok {
		if e, ok2 := err.(*filter.UnknownLevelError); ok2 {
			*ule = e
			return true
		}
	}
	return false
}
