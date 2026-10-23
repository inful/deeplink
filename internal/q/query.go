// Package q contains internal query-string and path-joining helpers
// shared by every deeplink sub-package. It is internal so the public
// surface of github.com/inful/deeplink stays focused on the deep-link
// builders themselves.
package q

import (
	"fmt"
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
	//
	// Raw values are validated by Encode: a value containing a character
	// that would let it escape its own field is rejected rather than
	// written out. See RawSafe.
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
//
// Encode returns an error if a Raw pair holds a value that could not be
// written without encoding (see RawSafe).
func Encode(pairs []Pair) (string, error) {
	first := true
	var b strings.Builder
	for _, p := range pairs {
		if p.Value == "" {
			continue
		}

		var value string
		if p.Raw {
			safe, err := RawSafe(p.Value)
			if err != nil {
				return "", fmt.Errorf("%s: %w", p.Key, err)
			}
			value = safe
		} else {
			value = encodeValue(p.Value)
		}

		if !first {
			b.WriteByte('&')
		}
		first = false
		b.WriteString(encodeKey(p.Key))
		b.WriteByte('=')
		b.WriteString(value)
	}
	return b.String(), nil
}

// RawSafe validates a value that must be inserted into a URL without
// encoding, and returns it unchanged.
//
// It rejects the characters that would let a value escape its own field:
// '&' starts another query parameter, '?' starts the query string, '#'
// starts the fragment, and whitespace (including CR and LF) either breaks
// the URL or, in a mailto: URL, enables header injection.
//
// Legitimate values for the raw fields this library uses — email
// addresses, comma-separated recipient lists, and Teams bot IDs of the
// form "28:<guid>" — never contain these characters.
func RawSafe(value string) (string, error) {
	if strings.ContainsAny(value, "&?#\r\n\t ") {
		return "", fmt.Errorf("value contains a URL-structural or whitespace character: %q", value)
	}
	return value, nil
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
	return encodeValue(s)
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
