package pipeline_test

import (
	"testing"

	"github.com/cronlog/pipeline"
)

func TestLevelFilter_AllowsMatchingLevel(t *testing.T) {
	proc := pipeline.LevelFilter("info", "warn")
	e := pipeline.Entry{Level: "info", Message: "ok"}
	_, ok := proc(e)
	if !ok {
		t.Fatal("expected info to be allowed")
	}
}

func TestLevelFilter_DropsUnmatchedLevel(t *testing.T) {
	proc := pipeline.LevelFilter("info")
	e := pipeline.Entry{Level: "debug", Message: "verbose"}
	_, ok := proc(e)
	if ok {
		t.Fatal("expected debug to be dropped")
	}
}

func TestLevelFilter_CaseInsensitive(t *testing.T) {
	proc := pipeline.LevelFilter("ERROR")
	e := pipeline.Entry{Level: "error", Message: "boom"}
	_, ok := proc(e)
	if !ok {
		t.Fatal("expected case-insensitive match")
	}
}

func TestAddField_InjectsKeyValue(t *testing.T) {
	proc := pipeline.AddField("job", "backup")
	e := pipeline.Entry{Level: "info", Message: "done"}
	out, ok := proc(e)
	if !ok {
		t.Fatal("expected entry to pass")
	}
	if out.Fields["job"] != "backup" {
		t.Fatalf("expected job=backup, got %s", out.Fields["job"])
	}
}

func TestAddField_DoesNotMutateOriginal(t *testing.T) {
	origFields := map[string]string{"x": "1"}
	proc := pipeline.AddField("y", "2")
	e := pipeline.Entry{Fields: origFields}
	proc(e)
	if _, exists := origFields["y"]; exists {
		t.Fatal("original fields map should not be mutated")
	}
}

func TestMessagePrefix_PrependsString(t *testing.T) {
	proc := pipeline.MessagePrefix("[cron] ")
	e := pipeline.Entry{Level: "info", Message: "starting"}
	out, ok := proc(e)
	if !ok {
		t.Fatal("expected entry to pass")
	}
	if out.Message != "[cron] starting" {
		t.Fatalf("unexpected message: %s", out.Message)
	}
}
