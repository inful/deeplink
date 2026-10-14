package outlook

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"text/template"
)

func TestFuncMapExposesExpectedNames(t *testing.T) {
	fm := FuncMap()
	for _, name := range []string{"outlookEvent"} {
		if _, ok := fm[name]; !ok {
			t.Errorf("FuncMap() missing %q", name)
		}
	}
}

func TestFuncMapOutlookEventFromTemplate(t *testing.T) {
	tpl := `{{ outlookEvent (dict "Subject" "Standup" "StartDate" "2024-01-15T10:00:00.000" "EndDate" "2024-01-15T10:15:00.000") }}`
	got, err := renderTemplate(tpl, FuncMap())
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	want := "https://outlook.office.com/calendar/deeplink/compose?enddt=2024-01-15T10%3A15%3A00.000&startdt=2024-01-15T10%3A00%3A00.000&subject=Standup"
	if got != want {
		t.Errorf("template output =\n  %q\nwant:\n  %q", got, want)
	}
}

func TestFuncMapMissingSubjectReturnsError(t *testing.T) {
	tpl := `{{ outlookEvent (dict "StartDate" "2024-01-15T10:00:00.000" "EndDate" "2024-01-15T11:00:00.000") }}`
	_, err := renderTemplate(tpl, FuncMap())
	if err == nil {
		t.Fatal("expected error for missing Subject, got nil")
	}
}

func TestFuncMapOutlookEmailFromTemplate(t *testing.T) {
	tpl := `{{ outlookEmail (dict "To" (slice "a@x.com" "b@x.com") "Subject" "Hi" "Body" "There") }}`
	got, err := renderTemplate(tpl, FuncMap())
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	want := "mailto:a@x.com,b@x.com?subject=Hi&body=There"
	if got != want {
		t.Errorf("template output =\n  %q\nwant:\n  %q", got, want)
	}
}

func TestFuncMapOutlookEmailAcceptsAnySlice(t *testing.T) {
	// Hugo's .Site.Params often yields []any when the value is parsed from
	// YAML or JSON; make sure the email funcmap handles that.
	tpl := `{{ outlookEmail (dict "To" (slice "a@x.com")) }}`
	got, err := renderTemplate(tpl, FuncMap())
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(got, "mailto:a@x.com") {
		t.Errorf("expected mailto: in %q", got)
	}
}

func TestFuncMapOutlookEmailNoRecipientsReturnsError(t *testing.T) {
	tpl := `{{ outlookEmail (dict "Subject" "Hi") }}`
	_, err := renderTemplate(tpl, FuncMap())
	if err == nil {
		t.Fatal("expected error for no recipients, got nil")
	}
}

func renderTemplate(body string, funcs template.FuncMap) (string, error) {
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
	funcs["slice"] = func(values ...any) []any { return values }

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
