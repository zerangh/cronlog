package masker

import "errors"

// ErrNoTargets is returned when New is called with no target strings.
var ErrNoTargets = errors.New("masker: at least one target string is required")

// ErrEmptyTarget is returned when one of the provided target strings is empty.
var ErrEmptyTarget = errors.New("masker: target string must not be empty")
