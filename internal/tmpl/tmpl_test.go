package tmpl

import (
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	m := map[string]any{
		"present": "value",
		"wrong":   42,
		"empty":   "",
	}

	if got := String(m, "present"); got != "value" {
		t.Errorf("String(present) = %q, want %q", got, "value")
	}
	if got := String(m, "wrong"); got != "" {
		t.Errorf("String(wrong) = %q, want empty", got)
	}
	if got := String(m, "empty"); got != "" {
		t.Errorf("String(empty) = %q, want empty", got)
	}
	if got := String(m, "missing"); got != "" {
		t.Errorf("String(missing) = %q, want empty", got)
	}
}

func TestBool(t *testing.T) {
	m := map[string]any{
		"true":     true,
		"false":    false,
		"strTrue":  "true",
		"strOne":   "1",
		"strYes":   "yes",
		"strFalse": "false",
		"strOther": "maybe",
		"notABool": 42,
	}

	tests := []struct {
		key  string
		want bool
	}{
		{"true", true},
		{"false", false},
		{"strTrue", true},
		{"strOne", true},
		{"strYes", true},
		{"strFalse", false},
		{"strOther", false},
		{"notABool", false},
		{"missing", false},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := Bool(m, tt.key); got != tt.want {
				t.Errorf("Bool(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	m := map[string]any{
		"int":       42,
		"int64":     int64(99),
		"float64":   float64(7),
		"numString": "123",
		"badString": "abc",
		"notANum":   true,
	}

	tests := []struct {
		key  string
		want int
	}{
		{"int", 42},
		{"int64", 99},
		{"float64", 7},
		{"numString", 123},
		{"badString", 0},
		{"notANum", 0},
		{"missing", 0},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := Int(m, tt.key); got != tt.want {
				t.Errorf("Int(%q) = %d, want %d", tt.key, got, tt.want)
			}
		})
	}
}

func TestStringList(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want []string
	}{
		{
			name: "string slice passes through",
			in:   []string{"a", "b"},
			want: []string{"a", "b"},
		},
		{
			name: "any slice is converted",
			in:   []any{"a", "b"},
			want: []string{"a", "b"},
		},
		{
			name: "non-string entries in an any slice are dropped",
			in:   []any{"a", 42, "b"},
			want: []string{"a", "b"},
		},
		{
			name: "a plain string is not a list",
			in:   "a",
			want: nil,
		},
		{
			name: "nil is not a list",
			in:   nil,
			want: nil,
		},
		{
			name: "an int is not a list",
			in:   42,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringList(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringList(%v) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}
