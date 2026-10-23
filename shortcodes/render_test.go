package shortcodes_test

import (
	"html"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/inful/deeplink/gitlab"
	"github.com/inful/deeplink/outlook"
	"github.com/inful/deeplink/shortcodes"
	"github.com/inful/deeplink/teams"
)

// The Hugo shortcodes are a second implementation of the URL-building
// logic in the Go library. This test renders every shortcode with Hugo
// and checks the result is byte-identical to the Go function it mirrors.
//
// It is skipped when Hugo is not installed.
func TestShortcodesMatchTheGoLibrary(t *testing.T) {
	hugo, err := exec.LookPath("hugo")
	if err != nil {
		t.Skip("hugo is not installed; skipping the shortcode rendering test")
	}

	cases := []struct {
		key        string
		shortcode  string
		attributes string
		want       string
	}{
		{
			key:        "GITLAB_ISSUE",
			shortcode:  "deeplink-gitlab-issue",
			attributes: `baseURL="https://gitlab.com/g/p" title="Bug: login broken" description="Steps to reproduce" issueType="issue" confidential="true" addRelatedIssue="42"`,
			want:       mustIssue(t),
		},
		{
			key:        "GITLAB_MR",
			shortcode:  "deeplink-gitlab-mergerequest",
			attributes: `baseURL="https://gitlab.com/g/p" title="Refactor login" description="See design doc" sourceBranch="feature/login" targetBranch="main" squash="true" navSource="docs"`,
			want:       mustMergeRequest(t),
		},
		{
			key:        "TEAMS_CHAT",
			shortcode:  "deeplink-teams-new-chat",
			attributes: `users="joe@contoso.com,bob@contoso.com" tenantId="f4c2b8e1-7a3d-4e6f-9b12-8c5d0e3a1f47" topicName="Prep For Meeting Tomorrow" message="Hi folks, kicking off a chat about our meeting tomorrow"`,
			want:       mustNewChat(t),
		},
		{
			key:        "TEAMS_CHANNEL",
			shortcode:  "deeplink-teams-channel",
			attributes: `channelId="19:9be3de4e70874c71a608dee9ba803ed3@thread.tacv2" channelName="My example channel" groupId="72602e12-78ac-474c-99d6-f619710353a9" tenantId="aaaabbbb-0000-cccc-1111-dddd2222eeee"`,
			want:       mustChannel(t),
		},
		{
			key:        "OUTLOOK_EVENT",
			shortcode:  "deeplink-outlook-event",
			attributes: `subject="Sprint review" body="Agenda: 1. intro" startDate="2024-01-15T10:00:00.000" endDate="2024-01-15T11:00:00.000" location="Conf room 4 / Teams"`,
			want:       mustEvent(t),
		},
		{
			key:        "OUTLOOK_EMAIL",
			shortcode:  "deeplink-outlook-email",
			attributes: `to="alice@example.com" subject="Question about the rollout" body="Hi Alice," cc="c@x.com"`,
			want:       mustEmail(t),
		},
	}

	// Build a minimal Hugo site that uses every shortcode once.
	site := t.TempDir()
	writeFile(t, filepath.Join(site, "hugo.toml"), "baseURL = \"https://example.com/\"\ntitle = \"test\"\n")
	writeFile(t, filepath.Join(site, "layouts", "_default", "single.html"), "<!DOCTYPE html><html><body>{{ .Content }}</body></html>\n")
	writeFile(t, filepath.Join(site, "layouts", "_default", "list.html"), "<!DOCTYPE html><html><body>list</body></html>\n")
	writeFile(t, filepath.Join(site, "layouts", "index.html"), "<!DOCTYPE html><html><body>home</body></html>\n")

	shortcodeDir := filepath.Join(site, "layouts", "shortcodes")
	if err := os.MkdirAll(shortcodeDir, 0o750); err != nil {
		t.Fatalf("mkdir shortcodes: %v", err)
	}
	for _, name := range shortcodes.Names() {
		content, err := shortcodes.Files.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		writeFile(t, filepath.Join(shortcodeDir, name), string(content))
	}

	var page strings.Builder
	page.WriteString("---\ntitle: Probe\n---\n")
	for _, c := range cases {
		page.WriteString(c.key + "={{< " + c.shortcode + " " + c.attributes + " >}}\n")
	}
	writeFile(t, filepath.Join(site, "content", "probe.md"), page.String())

	cmd := exec.Command(hugo, "--quiet", "--destination", "public")
	cmd.Dir = site
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("hugo failed: %v\n%s", err, out)
	}

	rendered, err := os.ReadFile(filepath.Join(site, "public", "probe", "index.html"))
	if err != nil {
		t.Fatalf("read rendered page: %v", err)
	}

	// Each rendered shortcode looks like: KEY=<a href="URL">text</a>.
	// The key is not necessarily at the start of a line — Goldmark wraps
	// it in a <p> — so anchor on the key itself.
	re := regexp.MustCompile(`([A-Z_]+)=<a href="([^"]*)"`)
	got := map[string]string{}
	for _, m := range re.FindAllStringSubmatch(string(rendered), -1) {
		// Hugo escapes the href for HTML; undo that so we compare the
		// URL itself, not its HTML representation.
		got[m[1]] = html.UnescapeString(m[2])
	}

	for _, c := range cases {
		t.Run(c.key, func(t *testing.T) {
			actual, ok := got[c.key]
			if !ok {
				t.Fatalf("shortcode %q produced no link in the rendered page", c.shortcode)
			}
			if actual != c.want {
				t.Errorf("shortcode %q produced\n  %q\nbut the Go library produces\n  %q",
					c.shortcode, actual, c.want)
			}
		})
	}
}

// TestShortcodesRejectUnsafeRecipients checks that the shortcodes refuse
// a recipients value that could add parameters of its own, matching the
// Go library.
func TestShortcodesRejectUnsafeRecipients(t *testing.T) {
	hugo, err := exec.LookPath("hugo")
	if err != nil {
		t.Skip("hugo is not installed; skipping the shortcode rendering test")
	}

	cases := []struct {
		name       string
		shortcode  string
		attributes string
	}{
		{
			name:       "teams users inject a parameter",
			shortcode:  "deeplink-teams-new-chat",
			attributes: `users="a@x.com&tenantId=EVIL"`,
		},
		{
			name:       "outlook to injects a parameter",
			shortcode:  "deeplink-outlook-email",
			attributes: `to="a@x.com?bcc=attacker@evil.com"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			site := t.TempDir()
			writeFile(t, filepath.Join(site, "hugo.toml"), "baseURL = \"https://example.com/\"\ntitle = \"test\"\n")
			writeFile(t, filepath.Join(site, "layouts", "_default", "single.html"), "<!DOCTYPE html><html><body>{{ .Content }}</body></html>\n")
			writeFile(t, filepath.Join(site, "layouts", "_default", "list.html"), "<!DOCTYPE html><html><body>list</body></html>\n")
			writeFile(t, filepath.Join(site, "layouts", "index.html"), "<!DOCTYPE html><html><body>home</body></html>\n")

			shortcodeDir := filepath.Join(site, "layouts", "shortcodes")
			if err := os.MkdirAll(shortcodeDir, 0o750); err != nil {
				t.Fatalf("mkdir shortcodes: %v", err)
			}
			for _, name := range shortcodes.Names() {
				content, err := shortcodes.Files.ReadFile(name)
				if err != nil {
					t.Fatalf("read %s: %v", name, err)
				}
				writeFile(t, filepath.Join(shortcodeDir, name), string(content))
			}

			writeFile(t, filepath.Join(site, "content", "probe.md"),
				"---\ntitle: Probe\n---\n"+
					"{{< "+c.shortcode+" "+c.attributes+" >}}\n")

			// No --quiet here: Hugo only prints the errorf message
			// without it, and that message is what we assert on.
			cmd := exec.Command(hugo, "--destination", "public")
			cmd.Dir = site
			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("hugo built the site; expected the shortcode to fail the build\n%s", out)
			}
			if !strings.Contains(string(out), "must not contain") {
				t.Errorf("hugo failed for the wrong reason:\n%s", out)
			}
		})
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mustIssue(t *testing.T) string {
	t.Helper()
	u, err := gitlab.Issue(gitlab.IssueParams{
		BaseURL: "https://gitlab.com/g/p", Title: "Bug: login broken",
		Description: "Steps to reproduce", IssueType: gitlab.IssueTypeIssue,
		Confidential: true, AddRelatedIssue: 42,
	})
	if err != nil {
		t.Fatalf("gitlab.Issue: %v", err)
	}
	return u
}

func mustMergeRequest(t *testing.T) string {
	t.Helper()
	u, err := gitlab.MergeRequest(gitlab.MergeRequestParams{
		BaseURL: "https://gitlab.com/g/p", Title: "Refactor login",
		Description: "See design doc", SourceBranch: "feature/login",
		TargetBranch: "main", Squash: true, NavSource: "docs",
	})
	if err != nil {
		t.Fatalf("gitlab.MergeRequest: %v", err)
	}
	return u
}

func mustNewChat(t *testing.T) string {
	t.Helper()
	u, err := teams.NewChat(teams.ChatParams{
		TenantID:  "f4c2b8e1-7a3d-4e6f-9b12-8c5d0e3a1f47",
		Users:     []string{"joe@contoso.com", "bob@contoso.com"},
		TopicName: "Prep For Meeting Tomorrow",
		Message:   "Hi folks, kicking off a chat about our meeting tomorrow",
	})
	if err != nil {
		t.Fatalf("teams.NewChat: %v", err)
	}
	return u
}

func mustChannel(t *testing.T) string {
	t.Helper()
	u, err := teams.Channel(teams.ChannelParams{
		ChannelID:   "19:9be3de4e70874c71a608dee9ba803ed3@thread.tacv2",
		ChannelName: "My example channel",
		GroupID:     "72602e12-78ac-474c-99d6-f619710353a9",
		TenantID:    "aaaabbbb-0000-cccc-1111-dddd2222eeee",
	})
	if err != nil {
		t.Fatalf("teams.Channel: %v", err)
	}
	return u
}

func mustEvent(t *testing.T) string {
	t.Helper()
	u, err := outlook.Event(outlook.EventParams{
		Subject: "Sprint review", Body: "Agenda: 1. intro",
		StartDate: "2024-01-15T10:00:00.000", EndDate: "2024-01-15T11:00:00.000",
		Location: "Conf room 4 / Teams",
	})
	if err != nil {
		t.Fatalf("outlook.Event: %v", err)
	}
	return u
}

func mustEmail(t *testing.T) string {
	t.Helper()
	u, err := outlook.Email(outlook.EmailParams{
		To: []string{"alice@example.com"}, Subject: "Question about the rollout",
		Body: "Hi Alice,", CC: []string{"c@x.com"},
	})
	if err != nil {
		t.Fatalf("outlook.Email: %v", err)
	}
	return u
}
