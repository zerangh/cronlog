package pipeline

import "strings"

// LevelFilter returns a Processor that drops entries whose level is not
// in the allowed set. Comparison is case-insensitive.
func LevelFilter(allowed ...string) Processor {
	set := make(map[string]struct{}, len(allowed))
	for _, l := range allowed {
		set[strings.ToLower(l)] = struct{}{}
	}
	return func(e Entry) (Entry, bool) {
		_, ok := set[strings.ToLower(e.Level)]
		return e, ok
	}
}

// AddField returns a Processor that injects a static key/value pair into
// every entry's Fields map without mutating the original map.
func AddField(key, value string) Processor {
	return func(e Entry) (Entry, bool) {
		merged := make(map[string]string, len(e.Fields)+1)
		for k, v := range e.Fields {
			merged[k] = v
		}
		merged[key] = value
		e.Fields = merged
		return e, true
	}
}

// MessagePrefix returns a Processor that prepends a string to every
// entry's message.
func MessagePrefix(prefix string) Processor {
	return func(e Entry) (Entry, bool) {
		e.Message = prefix + e.Message
		return e, true
	}
}
