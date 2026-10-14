package deeplink

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"text/template"

	"github.com/inful/deeplink/gitlab"
	"github.com/inful/deeplink/outlook"
	"github.com/inful/deeplink/teams"
)

func TestFuncmapCombinesAllSubPackages(t *testing.T) {
	combined := Funcmap()

	want := map[string]bool{}
	for _, src := range []map[string]any{
		gitlab.FuncMap(),
		teams.FuncMap(),
		outlook.FuncMap(),
	} {
		for k := range src {
			want[k] = true
		}
	}

	for name := range want {
		if _, ok := combined[name]; !ok {
			t.Errorf("Funcmap() missing %q from sub-package", name)
		}
	}
}

func TestNoFuncNameCollisions(t *testing.T) {
	combined := Funcmap()
	seen := make(map[string]string, len(combined))
	for name, fn := range combined {
		// every value must be a function (template.FuncMap constraint)
		if fn == nil {
			t.Errorf("Funcmap() entry %q is nil", name)
		}
		if prev, ok := seen[name]; ok {
			t.Errorf("Funcmap() name %q declared by both %s and another source", name, prev)
		}
		seen[name] = name
	}
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
