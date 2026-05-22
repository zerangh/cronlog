package output_test

import (
	"testing"

	"github.com/yourorg/cronlog/output"
)

func TestParseFormat_JSON(t *testing.T) {
	f, err := output.ParseFormat("json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != output.FormatJSON {
		t.Errorf("expected FormatJSON, got %s", f)
	}
}

func TestParseFormat_Text(t *testing.T) {
	f, err := output.ParseFormat("text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != output.FormatText {
		t.Errorf("expected FormatText, got %s", f)
	}
}

func TestParseFormat_EmptyDefaultsToText(t *testing.T) {
	f, err := output.ParseFormat("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != output.FormatText {
		t.Errorf("expected FormatText for empty string, got %s", f)
	}
}

func TestParseFormat_CaseInsensitive(t *testing.T) {
	f, err := output.ParseFormat("JSON")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != output.FormatJSON {
		t.Errorf("expected FormatJSON, got %s", f)
	}
}

func TestParseFormat_Unknown(t *testing.T) {
	_, err := output.ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestFormat_IsValid(t *testing.T) {
	if !output.FormatJSON.IsValid() {
		t.Error("FormatJSON should be valid")
	}
	if !output.FormatText.IsValid() {
		t.Error("FormatText should be valid")
	}
	if output.Format("csv").IsValid() {
		t.Error("unknown format should not be valid")
	}
}
