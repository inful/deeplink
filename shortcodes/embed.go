// Package shortcodes holds the Hugo shortcode HTML templates shipped
// with this library. The files are embedded into the Go module via
// the embed directive so a build tool (e.g. docbuilder) can copy them
// into a generated Hugo site's layouts/shortcodes/ directory without
// having to track the shortcode files separately.
package shortcodes

import "embed"

// Files holds the raw contents of every Hugo shortcode shipped with
// this library. The keys are filenames (e.g. "gitlab-issue.html"); the
// values are the file contents, ready to be written into a Hugo site's
// layouts/shortcodes/<name>.html.
//
// Example: copy every shortcode into a Hugo site:
//
//	for name, content := range shortcodes.Files() {
//	    dst := filepath.Join(siteRoot, "layouts/shortcodes", name)
//	    os.WriteFile(dst, content, 0o644)
//	}
//
//go:embed *.html
var Files embed.FS

// Names returns the basenames of every embedded shortcode file
// (e.g. []string{"gitlab-issue.html", "teams-new-chat.html", ...}).
// Useful for build tools that need to enumerate the set.
func Names() []string {
	entries, err := Files.ReadDir(".")
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out
}
