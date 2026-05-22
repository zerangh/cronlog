// Package filter provides log-level filtering for cronlog.
//
// It defines a set of severity levels (debug, info, warn, error) and a
// Filter type that decides whether a given log entry should be emitted
// based on a configured minimum level.
//
// Usage:
//
//	f, err := filter.ParseLevel(os.Getenv("CRONLOG_LEVEL"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	filter := filter.New(f)
//	if filter.Allow(filter.LevelWarn) {
//		// emit the entry
//	}
package filter
