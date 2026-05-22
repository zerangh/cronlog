package dedup

import "errors"

// ErrInvalidWindow is returned when a non-positive deduplication window is provided.
var ErrInvalidWindow = errors.New("dedup: window must be greater than zero")
