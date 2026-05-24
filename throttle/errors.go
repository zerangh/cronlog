package throttle

// sentinel errors are defined in throttle.go to keep the package surface
// minimal. This file is reserved for future structured error types.
//
// Example of a future typed error:
//
//	type LimitError struct {
//		Cap int
//	}
//
//	func (e *LimitError) Error() string {
//		return fmt.Sprintf("throttle: all %d slots occupied", e.Cap)
//	}
