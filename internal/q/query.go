// Package q contains internal query-string and path-joining helpers
// shared by every deeplink sub-package. It is internal so the public
// surface of github.com/inful/deeplink stays focused on the deep-link
// builders themselves.
package q

import (
	"net/url"
	"strings"
)

// Pair is a single key/value pair to be encoded into a URL query string.
//
// Keys MAY contain literal '[' and ']' characters; they are preserved as-is
// (required by systems like GitLab, which use namespaced params such as
// "issue[title]"). All other characters in keys and values are URL-encoded
// unless Raw is true.
type Pair struct {
	Key   string
	Value string
	// Raw, when true, inserts the value verbatim with no URL encoding.
	// Use this for values that the destination system expects to receive
	// in a specific non-encoded form (for example, the Teams `users`
	// parameter, which the documentation shows as a comma-separated list
	// rather than a single comma-encoded value). The Key is still
	// encoded; only the Value is preserved literally.
	Raw bool
}

// Encode builds a URL query string from the given pairs. It differs
// from url.Values.Encode() in two ways that match the format used in most
// vendor documentation:
//
//  1. Spaces are encoded as "%20" (not "+"), matching the canonical URLs
//     published by GitLab, Microsoft, and others.
//  2. Literal '[' and ']' in keys are preserved (not encoded as "%5B" /
//     "%5D"), again to match vendor docs.
//
// Pairs with an empty Value are omitted entirely. Callers that need to
// distinguish "not set" from "set to empty" should not rely on this
// function.
func Encode(pairs []Pair) string {
	if len(pairs) == 0 {
		return ""
	}

	first := true
	var b strings.Builder
	for _, p := range pairs {
		if p.Value == "" {
			continue
		}
		if !first {
			b.WriteByte('&')
		}
		first = false
		b.WriteString(encodeKey(p.Key))
		b.WriteByte('=')
		if p.Raw {
			b.WriteString(p.Value)
		} else {
			b.WriteString(encodeValue(p.Value))
		}
	}
	return b.String()
}

// JoinPath joins a base URL and one or more path segments, inserting a
// single "/" between them. Leading and trailing slashes on both the base
// and the segments are normalized away so the result never contains "//".
func JoinPath(base string, segments ...string) string {
	out := strings.TrimRight(base, "/")
	for _, s := range segments {
		s = strings.Trim(s, "/")
		if s == "" {
			continue
		}
		out += "/" + s
	}
	return out
}

// PathSegment escapes a string for use as a single URL path segment,
// URL-encoding every reserved character. This matches the encoding used
// in vendor documentation such as Microsoft Teams deep links (which use
// query-style encoding even for path components).
//
// Unlike url.PathEscape, this escapes '&' (a sub-delim) and '/' (the
// path separator), so a path segment can never be confused with a query
// parameter separator or split into multiple segments. It also uses %20
// (not '+') for spaces, matching query encoding conventions.
func PathSegment(s string) string {
	encoded := url.QueryEscape(s)
	return strings.ReplaceAll(encoded, "+", "%20")
}

// encodeKey URL-encodes a query key while preserving literal '[' and ']'.
func encodeKey(s string) string {
	encoded := url.QueryEscape(s)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "%5B", "[")
	encoded = strings.ReplaceAll(encoded, "%5D", "]")
	return encoded
}

// encodeValue URL-encodes a query value, using %20 for spaces.
func encodeValue(s string) string {
	encoded := url.QueryEscape(s)
	return strings.ReplaceAll(encoded, "+", "%20")
}
