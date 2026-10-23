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
	content, err := Files.ReadFile("deeplink-gitlab-issue.html")
	if err != nil {
		t.Fatalf("ReadFile(deeplink-gitlab-issue.html): %v", err)
	}
	if !strings.Contains(string(content), "deeplink-gitlab-issue") {
		t.Errorf("deeplink-gitlab-issue.html does not look like the expected shortcode")
	}
}

// TestShortcodeNameMatchesFilename guards the contract Hugo relies on:
// a shortcode invoked as {{< deeplink-foo >}} must live in a file named
// deeplink-foo.html, or Hugo silently fails to find it.
func TestShortcodeNameMatchesFilename(t *testing.T) {
	for _, name := range Names() {
		content, err := Files.ReadFile(name)
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", name, err)
		}
		want := "Shortcode: " + strings.TrimSuffix(name, ".html")
		if !strings.Contains(string(content), want) {
			t.Errorf("%s does not declare %q; Hugo looks the shortcode up by filename", name, want)
		}
	}
}
