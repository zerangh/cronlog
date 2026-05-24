package dispatcher_test

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"

	"cronlog/dispatcher"
)

func TestDispatch_ConcurrentWrites(t *testing.T) {
	var buf safeBuffer
	d, err := dispatcher.New(dispatcher.NewRoute("concurrent", &buf, "info"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const workers = 20
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(n int) {
			defer wg.Done()
			d.Dispatch(dispatcher.Entry{
				Level:   "info",
				Message: fmt.Sprintf("worker-%d", n),
			})
		}(i)
	}
	wg.Wait()

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != workers {
		t.Errorf("expected %d lines, got %d", workers, len(lines))
	}
}

// safeBuffer is a thread-safe bytes.Buffer for testing.
type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *safeBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *safeBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}
