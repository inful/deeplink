package shortcodes

import (
	"strings"
	"testing"
)

func TestFilesAreEmbedded(t *testing.T) {
	entries, err := Files.ReadDir(".")
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no shortcode files embedded; expected at least one")
	}
}

func TestNames(t *testing.T) {
	names := Names()
	if len(names) == 0 {
		t.Fatal("Names() returned empty slice")
	}
	for _, n := range names {
		if !strings.HasSuffix(n, ".html") {
			t.Errorf("Names() contains non-html file %q", n)
		}
	}
}

func TestGitLabIssueShortcodeIsEmbedded(t *testing.T) {
	content, err := Files.ReadFile("gitlab-issue.html")
	if err != nil {
		t.Fatalf("ReadFile(gitlab-issue.html): %v", err)
	}
	if !strings.Contains(string(content), "deeplink-gitlab-issue") {
		t.Errorf("gitlab-issue.html does not look like the expected shortcode")
	}
}
