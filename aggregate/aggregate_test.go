package aggregate_test

import (
	"sync"
	"testing"

	"github.com/cronlog/aggregate"
)

func TestAdd_StoresEntry(t *testing.T) {
	c := aggregate.New()
	c.Add("info", "job started", nil)

	s := c.Summarise()
	if s.Total != 1 {
		t.Fatalf("expected total 1, got %d", s.Total)
	}
}

func TestSummarise_GroupsByLevel(t *testing.T) {
	c := aggregate.New()
	c.Add("info", "step one", nil)
	c.Add("error", "something failed", map[string]any{"code": 500})
	c.Add("info", "step two", nil)

	s := c.Summarise()

	if s.Counts["info"] != 2 {
		t.Errorf("expected 2 info entries, got %d", s.Counts["info"])
	}
	if s.Counts["error"] != 1 {
		t.Errorf("expected 1 error entry, got %d", s.Counts["error"])
	}
	if len(s.ByLevel["error"]) != 1 {
		t.Errorf("expected 1 entry in ByLevel[error], got %d", len(s.ByLevel["error"]))
	}
}

func TestSummarise_EmptyCollector(t *testing.T) {
	c := aggregate.New()
	s := c.Summarise()

	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
	if len(s.ByLevel) != 0 {
		t.Errorf("expected empty ByLevel map")
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	c := aggregate.New()
	c.Add("warn", "low disk", nil)
	c.Reset()

	s := c.Summarise()
	if s.Total != 0 {
		t.Errorf("expected total 0 after reset, got %d", s.Total)
	}
}

func TestAdd_ConcurrentWrites(t *testing.T) {
	c := aggregate.New()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Add("info", "concurrent", nil)
		}()
	}
	wg.Wait()

	s := c.Summarise()
	if s.Total != 50 {
		t.Errorf("expected 50 entries, got %d", s.Total)
	}
}
