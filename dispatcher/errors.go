package dispatcher

import "errors"

// ErrNilWriter is returned when a Route is constructed with a nil writer.
var ErrNilWriter = errors.New("dispatcher: writer must not be nil")

// ErrEmptyRouteName is returned when a Route has an empty name.
var ErrEmptyRouteName = errors.New("dispatcher: route name must not be empty")

// Validate checks that a Route is properly configured.
func (r *Route) Validate() error {
	if r.Name == "" {
		return ErrEmptyRouteName
	}
	if r.Writer == nil {
		return ErrNilWriter
	}
	return nil
}
