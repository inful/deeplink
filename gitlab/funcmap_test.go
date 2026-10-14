package gitlab

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"text/template"
)

func TestFuncMapExposesExpectedNames(t *testing.T) {
	fm := FuncMap()
	for _, name := range []string{"gitlabIssue", "gitlabMergeRequest"} {
		if _, ok := fm[name]; !ok {
			t.Errorf("FuncMap() missing %q", name)
		}
	}
}

func TestFuncMapGitlabIssueFromTemplate(t *testing.T) {
	tpl := `{{ gitlabIssue (dict "BaseURL" "https://gitlab.com/g/p" "Title" "Hi" "IssueType" "issue") }}`
	got, err := renderTemplate(tpl, FuncMap())
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	want := "https://gitlab.com/g/p/-/work_items/new?issue[title]=Hi&issue[issue_type]=issue"
	if got != want {
		t.Errorf("template output =\n  %q\nwant:\n  %q", got, want)
	}
}

func TestFuncMapMissingBaseURLReturnsError(t *testing.T) {
	tpl := `{{ gitlabIssue (dict "Title" "Hi") }}`
	_, err := renderTemplate(tpl, FuncMap())
	if err == nil {
		t.Fatal("expected error for missing BaseURL, got nil")
	}
	if !strings.Contains(err.Error(), "BaseURL is required") {
		t.Errorf("error = %q, want it to mention BaseURL", err)
	}
}

func TestFuncMapGitlabMergeRequestFromTemplate(t *testing.T) {
	tpl := `{{ gitlabMergeRequest (dict "BaseURL" "https://gitlab.com/g/p" "Title" "Refactor" "SourceBranch" "feat/x" "TargetBranch" "main") }}`
	got, err := renderTemplate(tpl, FuncMap())
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	want := "https://gitlab.com/g/p/-/merge_requests/new?merge_request[title]=Refactor&merge_request[source_branch]=feat%2Fx&merge_request[target_branch]=main"
	if got != want {
		t.Errorf("template output =\n  %q\nwant:\n  %q", got, want)
	}
}

func TestFuncMapAcceptsStringBoolFromYAML(t *testing.T) {
	// Hugo/YAML often hand us strings for booleans; make sure we accept them.
	tpl := `{{ gitlabIssue (dict "BaseURL" "https://gitlab.com/g/p" "Title" "Secret" "Confidential" "true") }}`
	got, err := renderTemplate(tpl, FuncMap())
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(got, "issue[confidential]=true") {
		t.Errorf("expected confidential=true in URL, got %q", got)
	}
}

func TestFuncMapAcceptsVariousIntForms(t *testing.T) {
	// JSON decodes all numbers to float64; YAML often uses int or int64.
	// We accept all of them for AddRelatedIssue.
	cases := []struct {
		name string
		val  any
		want string
	}{
		{"int", 42, "add_related_issue=42"},
		{"int64", int64(99), "add_related_issue=99"},
		{"float64", float64(7), "add_related_issue=7"},
		{"numeric string", "123", "add_related_issue=123"},
		{"non-numeric string", "abc", ""}, // Atoi fails, returns 0 → omitted
	}

	dyn := template.New("dyn").Funcs(FuncMap()).Funcs(template.FuncMap{
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("dict: odd number of arguments")
			}
			m := map[string]any{}
			for i := 0; i < len(values); i += 2 {
				m[values[i].(string)] = values[i+1] //nolint:errcheck // test helper; misuse panics, which is fine
			}
			return m, nil
		},
	})
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body := `{{ gitlabIssue (dict "BaseURL" "https://gitlab.com/g/p" "Title" "x" "AddRelatedIssue" .v) }}`
			tpl, err := dyn.Parse(body)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			var buf bytes.Buffer
			if err := tpl.Execute(&buf, map[string]any{"v": c.val}); err != nil {
				t.Fatalf("execute: %v", err)
			}
			got := buf.String()
			if c.want != "" {
				if !strings.Contains(got, c.want) {
					t.Errorf("expected %q in %q", c.want, got)
				}
			} else if strings.Contains(got, "add_related_issue") {
				t.Errorf("did not expect add_related_issue in %q", got)
			}
		})
	}
}

func renderTemplate(body string, funcs template.FuncMap) (string, error) {
	// Inject a `dict` helper for tests; Hugo has this built-in but Go's
	// text/template does not.
	if funcs == nil {
		funcs = template.FuncMap{}
	}
	funcs["dict"] = func(values ...any) (map[string]any, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("dict: odd number of arguments")
		}
		m := make(map[string]any, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, errors.New("dict: keys must be strings")
			}
			m[key] = values[i+1]
		}
		return m, nil
	}

	tpl, err := template.New("test").Funcs(funcs).Parse(body)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, nil); err != nil {
		return "", err
	}
	return buf.String(), nil
}
