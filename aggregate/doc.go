// Package aggregate provides an in-memory log entry collector that groups
// captured entries by severity level and exposes a structured Summary.
//
// Typical usage within a cron job run:
//
//	col := aggregate.New()
//	col.Add("info",  "starting backup", nil)
//	col.Add("error", "disk full",       map[string]any{"path": "/var"})
//
//	s := col.Summarise()
//	fmt.Println(s.Counts["error"]) // 1
//
// Collector is safe for concurrent use.
package aggregate
