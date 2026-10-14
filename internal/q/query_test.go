package q

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		name  string
		pairs []Pair
		want  string
	}{
		{
			name:  "empty",
			pairs: nil,
			want:  "",
		},
		{
			name:  "single pair",
			pairs: []Pair{{Key: "a", Value: "b"}},
			want:  "a=b",
		},
		{
			name:  "multiple pairs",
			pairs: []Pair{{Key: "a", Value: "b"}, {Key: "c", Value: "d"}},
			want:  "a=b&c=d",
		},
		{
			name:  "space becomes %20 not +",
			pairs: []Pair{{Key: "title", Value: "hello world"}},
			want:  "title=hello%20world",
		},
		{
			name:  "brackets preserved in keys (GitLab style)",
			pairs: []Pair{{Key: "issue[title]", Value: "My issue"}},
			want:  "issue[title]=My%20issue",
		},
		{
			name:  "brackets preserved alongside other pairs",
			pairs: []Pair{{Key: "issue[title]", Value: "Foo"}, {Key: "issue[description]", Value: "Bar"}},
			want:  "issue[title]=Foo&issue[description]=Bar",
		},
		{
			name:  "special characters in values",
			pairs: []Pair{{Key: "x", Value: "a&b=c"}},
			want:  "x=a%26b%3Dc",
		},
		{
			name:  "empty values skipped",
			pairs: []Pair{{Key: "a", Value: "b"}, {Key: "c", Value: ""}, {Key: "d", Value: "e"}},
			want:  "a=b&d=e",
		},
		{
			name:  "all empty values yields empty string",
			pairs: []Pair{{Key: "a", Value: ""}, {Key: "b", Value: ""}},
			want:  "",
		},
		{
			name:  "unicode preserved",
			pairs: []Pair{{Key: "subject", Value: "Møte med teamet"}},
			want:  "subject=M%C3%B8te%20med%20teamet",
		},
		{
			name: "canonical GitLab example",
			pairs: []Pair{
				{Key: "issue[title]", Value: "Whoa, we're half-way there"},
				{Key: "issue[description]", Value: "Whoa, livin' in a URL"},
			},
			want: "issue[title]=Whoa%2C%20we%27re%20half-way%20there&issue[description]=Whoa%2C%20livin%27%20in%20a%20URL",
		},
		{
			name:  "boolean serialized as true/false",
			pairs: []Pair{{Key: "issue[confidential]", Value: "true"}},
			want:  "issue[confidential]=true",
		},
		{
			name:  "raw value is not URL-encoded (Teams users parameter)",
			pairs: []Pair{{Key: "users", Value: "joe@contoso.com,bob@contoso.com", Raw: true}},
			want:  "users=joe@contoso.com,bob@contoso.com",
		},
		{
			name:  "raw value with empty result is still skipped",
			pairs: []Pair{{Key: "users", Value: "", Raw: true}},
			want:  "",
		},
		{
			name: "mixed raw and encoded pairs",
			pairs: []Pair{
				{Key: "users", Value: "a@x.com,b@x.com", Raw: true},
				{Key: "message", Value: "Hi folks"},
			},
			want: "users=a@x.com,b@x.com&message=Hi%20folks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.pairs)
			if got != tt.want {
				t.Errorf("Encode() =\n  %q\nwant:\n  %q", got, tt.want)
			}
		})
	}
}

func TestJoinPath(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		segments []string
		want     string
	}{
		{
			name:     "no segments returns trimmed base",
			base:     "https://example.com/foo/",
			segments: nil,
			want:     "https://example.com/foo",
		},
		{
			name:     "single segment",
			base:     "https://example.com",
			segments: []string{"foo"},
			want:     "https://example.com/foo",
		},
		{
			name:     "segment with leading slash normalised",
			base:     "https://example.com/",
			segments: []string{"/foo"},
			want:     "https://example.com/foo",
		},
		{
			name:     "multiple segments",
			base:     "https://example.com",
			segments: []string{"foo", "bar", "baz"},
			want:     "https://example.com/foo/bar/baz",
		},
		{
			name:     "empty segments skipped",
			base:     "https://example.com",
			segments: []string{"foo", "", "bar"},
			want:     "https://example.com/foo/bar",
		},
		{
			name:     "trailing slashes on base and segments normalised",
			base:     "https://example.com/",
			segments: []string{"/foo/", "/bar"},
			want:     "https://example.com/foo/bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JoinPath(tt.base, tt.segments...)
			if got != tt.want {
				t.Errorf("JoinPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPathSegment(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "plain text passes through",
			in:   "hello",
			want: "hello",
		},
		{
			name: "space becomes %20",
			in:   "hello world",
			want: "hello%20world",
		},
		{
			name: "slash IS escaped (so the path segment can never split)",
			in:   "a/b",
			want: "a%2Fb",
		},
		{
			name: "ampersand IS escaped (unlike url.PathEscape)",
			in:   "R&D",
			want: "R%26D",
		},
		{
			name: "colon and at-sign escaped (Microsoft Teams style)",
			in:   "19:9be3de4e70874c71a608dee9ba803ed3@thread.tacv2",
			want: "19%3A9be3de4e70874c71a608dee9ba803ed3%40thread.tacv2",
		},
		{
			name: "empty string",
			in:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PathSegment(tt.in)
			if got != tt.want {
				t.Errorf("PathSegment() = %q, want %q", got, tt.want)
			}
		})
	}
}
