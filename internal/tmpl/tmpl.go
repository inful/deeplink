// Package tmpl provides helpers for coercing template data
// (map[string]any) into the typed values the deep-link builders need.
//
// The helpers are deliberately forgiving about the types a value may
// arrive as, because template data reaches a library from YAML, JSON,
// Hugo's `dict`, or plain Go maps, and each represents booleans and
// numbers differently. A value that cannot be interpreted is treated as
// absent rather than as an error.
package tmpl

import "strconv"

// String returns the value at key as a string, or "" if it is missing or
// is not a string.
func String(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Bool returns the value at key as a bool. Strings are accepted so that
// values coming from YAML or JSON (where booleans are sometimes quoted)
// still work. Anything else is false.
func Bool(m map[string]any, key string) bool {
	switch b := m[key].(type) {
	case bool:
		return b
	case string:
		return b == "true" || b == "1" || b == "yes"
	}
	return false
}

// Int returns the value at key as an int. Strings are parsed, so a
// numeric string from YAML or JSON still works. A non-numeric string is
// treated as absent (0).
func Int(m map[string]any, key string) int {
	switch n := m[key].(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case string:
		// strconv.Atoi reports failure for non-numeric input; that
		// returns 0, which is the documented "absent" value here.
		i, _ := strconv.Atoi(n) //nolint:errcheck // Atoi failure (non-numeric input) returns 0, the documented "absent" sentinel
		return i
	}
	return 0
}

// StringList returns v as a slice of strings. Both []string and []any
// are accepted — the latter is what Hugo's `slice` function and values
// parsed from YAML or JSON produce. Any other type yields nil.
func StringList(v any) []string {
	switch xs := v.(type) {
	case []string:
		return xs
	case []any:
		out := make([]string, 0, len(xs))
		for _, x := range xs {
			if s, ok := x.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}
