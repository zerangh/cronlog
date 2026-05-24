// Package dispatcher routes log entries to one or more named destinations
// based on configurable routing rules.
package dispatcher

import (
	"errors"
	"io"
)

// ErrNoRoutes is returned when a Dispatcher is created with no routes.
var ErrNoRoutes = errors.New("dispatcher: at least one route is required")

// ErrDuplicateRoute is returned when two routes share the same name.
var ErrDuplicateRoute = errors.New("dispatcher: duplicate route name")

// Entry represents a log entry to be dispatched.
type Entry struct {
	Level   string
	Message string
	Fields  map[string]string
}

// Route pairs a named destination with a writer and an optional level filter.
type Route struct {
	Name       string
	Writer     io.Writer
	MinLevel   string
	levelOrder map[string]int
}

// Dispatcher routes entries to registered destinations.
type Dispatcher struct {
	routes []*Route
}

var levelOrder = map[string]int{
	"debug": 0,
	"info":  1,
	"warn":  2,
	"error": 3,
}

// New creates a Dispatcher from the provided routes.
func New(routes ...*Route) (*Dispatcher, error) {
	if len(routes) == 0 {
		return nil, ErrNoRoutes
	}
	seen := make(map[string]struct{}, len(routes))
	for _, r := range routes {
		if _, dup := seen[r.Name]; dup {
			return nil, ErrDuplicateRoute
		}
		seen[r.Name] = struct{}{}
		r.levelOrder = levelOrder
	}
	return &Dispatcher{routes: routes}, nil
}

// NewRoute constructs a Route with the given name, writer, and minimum level.
func NewRoute(name string, w io.Writer, minLevel string) *Route {
	return &Route{Name: name, Writer: w, MinLevel: minLevel}
}

// Dispatch sends e to every route whose minimum level is satisfied.
func (d *Dispatcher) Dispatch(e Entry) {
	for _, r := range d.routes {
		if r.allows(e.Level) {
			_, _ = io.WriteString(r.Writer, format(e))
		}
	}
}

func (r *Route) allows(level string) bool {
	if r.MinLevel == "" {
		return true
	}
	min, ok := r.levelOrder[r.MinLevel]
	if !ok {
		return true
	}
	got, ok := r.levelOrder[level]
	if !ok {
		return true
	}
	return got >= min
}

func format(e Entry) string {
	return e.Level + " " + e.Message + "\n"
}
