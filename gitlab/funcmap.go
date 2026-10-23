package gitlab

import (
	"text/template"

	"github.com/inful/deeplink/internal/tmpl"
)

// FuncMap returns a template.FuncMap containing every GitLab deep-link
// builder in this package, namespaced under the "gitlab" prefix.
//
// The map's values take a map[string]any (e.g. from Hugo's `dict` function
// or docbuilder's template data) and return the URL string. Conversion
// errors are returned as a template execution error so authors see
// problems immediately.
//
// Available functions:
//
//	gitlabIssue       — gitlab.IssueParams       → URL for a new work item
//	gitlabMergeRequest — gitlab.MergeRequestParams → URL for a new merge request
//
// Example (Hugo):
//
//	{{ $u := gitlabIssue (dict
//	    "BaseURL"     "https://gitlab.com/g/p"
//	    "Title"       "Login is broken"
//	    "Description" "Steps to reproduce ..."
//	) }}
//	<a href="{{ $u }}">File an issue</a>
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"gitlabIssue":        issueFromMap,
		"gitlabMergeRequest": mergeRequestFromMap,
	}
}

// issueFromMap is the template-friendly adapter for Issue. It accepts a
// map[string]any (typical of template data) and produces a URL string.
//
// Recognized keys match the fields of IssueParams (case-insensitive
// lookups are intentionally not performed — match the field name
// exactly).
func issueFromMap(m map[string]any) (string, error) {
	return Issue(IssueParams{
		BaseURL:             tmpl.String(m, "BaseURL"),
		Title:               tmpl.String(m, "Title"),
		Description:         tmpl.String(m, "Description"),
		IssueType:           IssueType(tmpl.String(m, "IssueType")),
		Confidential:        tmpl.Bool(m, "Confidential"),
		Template:            tmpl.String(m, "Template"),
		DescriptionTemplate: tmpl.String(m, "DescriptionTemplate"),
		AddRelatedIssue:     tmpl.Int(m, "AddRelatedIssue"),
	})
}

// mergeRequestFromMap is the template-friendly adapter for MergeRequest.
func mergeRequestFromMap(m map[string]any) (string, error) {
	return MergeRequest(MergeRequestParams{
		BaseURL:            tmpl.String(m, "BaseURL"),
		Title:              tmpl.String(m, "Title"),
		Description:        tmpl.String(m, "Description"),
		SourceBranch:       tmpl.String(m, "SourceBranch"),
		TargetBranch:       tmpl.String(m, "TargetBranch"),
		RemoveSourceBranch: tmpl.Bool(m, "RemoveSourceBranch"),
		Squash:             tmpl.Bool(m, "Squash"),
		NavSource:          tmpl.String(m, "NavSource"),
	})
}
