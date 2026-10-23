package deeplink

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"text/template"
)

func TestFuncmapExposesExactlyTheExpectedFunctions(t *testing.T) {
	// This is the library's public template API. Pinning the exact set
	// means adding, renaming, or removing a function has to be a
	// deliberate change to this list.
	want := []string{
		"gitlabIssue",
		"gitlabMergeRequest",
		"teamsNewChat",
		"teamsChannel",
		"outlookEvent",
		"outlookEmail",
	}

	got := Funcmap()
	for _, name := range want {
		fn, ok := got[name]
		if !ok {
			t.Errorf("Funcmap() missing %q", name)
			continue
		}
		if fn == nil {
			t.Errorf("Funcmap() entry %q is nil", name)
		}
	}
	if len(got) != len(want) {
		t.Errorf("Funcmap() has %d functions, want %d (update the want list if the API changed)", len(got), len(want))
	}
}

func TestMergePanicsOnDuplicateName(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("merge() did not panic on a duplicate function name")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic value is %T, want string", r)
		}
		if !strings.Contains(msg, "duplicate") || !strings.Contains(msg, "dupName") {
			t.Errorf("panic message %q should name the duplicated function", msg)
		}
	}()

	merge(
		template.FuncMap{"dupName": func() {}},
		template.FuncMap{"dupName": func() {}},
	)
}

// TestFuncmapEndToEnd exercises the full top-level Funcmap() the way a
// consumer would: register it with a text/template instance and render a
// template that calls one function from every sub-package.
func TestFuncmapEndToEnd(t *testing.T) {
	funcs := Funcmap()
	// Inject Hugo-like helpers (dict, slice) so the test is realistic.
	funcs["dict"] = func(values ...any) (map[string]any, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("dict: odd args")
		}
		m := make(map[string]any, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			m[values[i].(string)] = values[i+1] //nolint:errcheck // test helper; misuse panics, which is fine
		}
		return m, nil
	}
	funcs["slice"] = func(values ...any) []any { return values }

	tpl, err := template.New("page").Funcs(funcs).Parse(strings.Join([]string{
		`{{ $u := gitlabIssue (dict "BaseURL" "https://gitlab.com/g/p" "Title" "Hi") }}`,
		`{{ $v := teamsNewChat (dict "Users" (slice "a@x.com")) }}`,
		`{{ $w := outlookEvent (dict "Subject" "Standup") }}`,
		`GITLAB={{ $u }}`,
		`TEAMS={{ $v }}`,
		`OUTLOOK={{ $w }}`,
	}, "\n"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, nil); err != nil {
		t.Fatalf("execute: %v", err)
	}
	got := buf.String()
	for _, want := range []string{
		"GITLAB=https://gitlab.com/g/p/-/work_items/new?issue[title]=Hi",
		"TEAMS=https://teams.microsoft.com/l/chat/0/0?users=a@x.com",
		"OUTLOOK=https://outlook.office.com/calendar/deeplink/compose?subject=Standup",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q\nfull output:\n%s", want, got)
		}
	}
}
