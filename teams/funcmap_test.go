package teams

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"text/template"
)

func TestFuncMapExposesExpectedNames(t *testing.T) {
	fm := FuncMap()
	for _, name := range []string{"teamsNewChat"} {
		if _, ok := fm[name]; !ok {
			t.Errorf("FuncMap() missing %q", name)
		}
	}
}

func TestFuncMapTeamsNewChatFromTemplate(t *testing.T) {
	tpl := `{{ teamsNewChat (dict "Users" (slice "joe@x.com" "bob@x.com") "Message" "Quick q") }}`
	got, err := renderTemplate(tpl, FuncMap())
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	want := "https://teams.microsoft.com/l/chat/0/0?users=joe@x.com,bob@x.com&message=Quick%20q"
	if got != want {
		t.Errorf("template output =\n  %q\nwant:\n  %q", got, want)
	}
}

func TestFuncMapAcceptsAnySliceForUsers(t *testing.T) {
	// Hugo commonly hands back []any from .Site.Params; make sure we accept it.
	tpl := `{{ teamsNewChat (dict "Users" (slice "a@x.com" "b@x.com")) }}`
	got, err := renderTemplate(tpl, FuncMap())
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(got, "users=a@x.com,b@x.com") {
		t.Errorf("expected users parameter in URL, got %q", got)
	}
}

func TestFuncMapNoUsersReturnsError(t *testing.T) {
	tpl := `{{ teamsNewChat (dict) }}`
	_, err := renderTemplate(tpl, FuncMap())
	if err == nil {
		t.Fatal("expected error for missing Users, got nil")
	}
	if !strings.Contains(err.Error(), "at least one user") {
		t.Errorf("error = %q, want it to mention users", err)
	}
}

func TestFuncMapTeamsChannelFromTemplate(t *testing.T) {
	tpl := `{{ teamsChannel (dict "ChannelID" "19:abc@thread.tacv2" "ChannelName" "General" "GroupID" "g" "TenantID" "t") }}`
	got, err := renderTemplate(tpl, FuncMap())
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	want := "https://teams.microsoft.com/l/channel/19%3Aabc%40thread.tacv2/General?groupId=g&tenantId=t"
	if got != want {
		t.Errorf("template output =\n  %q\nwant:\n  %q", got, want)
	}
}

func TestFuncMapChannelPrivateFlag(t *testing.T) {
	tpl := `{{ teamsChannel (dict "ChannelID" "19:abc@thread.tacv2" "ChannelName" "x" "GroupID" "g" "TenantID" "t" "Private" true) }}`
	got, err := renderTemplate(tpl, FuncMap())
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(got, "ngc=true") {
		t.Errorf("expected ngc=true in URL, got %q", got)
	}
}

func TestFuncMapChannelMissingIDReturnsError(t *testing.T) {
	tpl := `{{ teamsChannel (dict "ChannelName" "x") }}`
	_, err := renderTemplate(tpl, FuncMap())
	if err == nil {
		t.Fatal("expected error for missing ChannelID, got nil")
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
	funcs["slice"] = func(values ...any) []any {
		return values
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
