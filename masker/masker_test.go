package masker_test

import (
	"testing"

	"github.com/cronlog/masker"
)

func TestNew_NoTargets(t *testing.T) {
	_, err := masker.New("[REDACTED]")
	if err != masker.ErrNoTargets {
		t.Fatalf("expected ErrNoTargets, got %v", err)
	}
}

func TestNew_EmptyTarget(t *testing.T) {
	_, err := masker.New("[REDACTED]", "valid", "")
	if err != masker.ErrEmptyTarget {
		t.Fatalf("expected ErrEmptyTarget, got %v", err)
	}
}

func TestNew_DefaultPlaceholder(t *testing.T) {
	m, err := masker.New("", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Mask("my secret value")
	want := "my [REDACTED] value"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_ReplacesTarget(t *testing.T) {
	m, _ := masker.New("***", "tok_abc123")
	got := m.Mask("Authorization: Bearer tok_abc123")
	want := "Authorization: Bearer ***"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_MultipleTargets(t *testing.T) {
	m, _ := masker.New("[REDACTED]", "password", "secret")
	got := m.Mask("password=hunter2 secret=xyz")
	want := "[REDACTED]=hunter2 [REDACTED]=xyz"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMask_NoMatch(t *testing.T) {
	m, _ := masker.New("[REDACTED]", "token")
	input := "nothing sensitive here"
	if got := m.Mask(input); got != input {
		t.Errorf("expected unchanged string, got %q", got)
	}
}

func TestMaskMap_DoesNotMutateInput(t *testing.T) {
	m, _ := masker.New("[REDACTED]", "s3cr3t")
	orig := map[string]string{"key": "value with s3cr3t inside"}
	origVal := orig["key"]
	out := m.MaskMap(orig)
	if orig["key"] != origVal {
		t.Error("original map was mutated")
	}
	if out["key"] == origVal {
		t.Error("output map was not masked")
	}
}

func TestTargets_ReturnsCopy(t *testing.T) {
	m, _ := masker.New("[REDACTED]", "alpha", "beta")
	targets := m.Targets()
	targets[0] = "mutated"
	if m.Targets()[0] == "mutated" {
		t.Error("Targets() returned a reference to internal slice")
	}
}
