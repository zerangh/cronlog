// Package truncate provides utilities for capping the length of log message
// strings and field values before they are written to output or forwarded
// to webhook notifiers.
//
// # Usage
//
//	tr, err := truncate.New(512)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Truncate a single string.
//	safe := tr.String(longMessage)
//
//	// Truncate all string values in a fields map.
//	safeFields := tr.Fields(rawFields)
//
// Strings that exceed the configured maximum length are trimmed and
// suffixed with "..." so that the result never exceeds maxLen bytes.
package truncate
