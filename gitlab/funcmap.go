package gitlab

import (
	"strconv"
	"text/template"
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
		BaseURL:             mapString(m, "BaseURL"),
		Title:               mapString(m, "Title"),
		Description:         mapString(m, "Description"),
		IssueType:           IssueType(mapString(m, "IssueType")),
		Confidential:        mapBool(m, "Confidential"),
		Template:            mapString(m, "Template"),
		DescriptionTemplate: mapString(m, "DescriptionTemplate"),
		AddRelatedIssue:     mapInt(m, "AddRelatedIssue"),
	})
}

// mergeRequestFromMap is the template-friendly adapter for MergeRequest.
func mergeRequestFromMap(m map[string]any) (string, error) {
	return MergeRequest(MergeRequestParams{
		BaseURL:            mapString(m, "BaseURL"),
		Title:              mapString(m, "Title"),
		Description:        mapString(m, "Description"),
		SourceBranch:       mapString(m, "SourceBranch"),
		TargetBranch:       mapString(m, "TargetBranch"),
		RemoveSourceBranch: mapBool(m, "RemoveSourceBranch"),
		Squash:             mapBool(m, "Squash"),
		NavSource:          mapString(m, "NavSource"),
	})
}

func mapString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func mapBool(m map[string]any, key string) bool {
	if v, ok := m[key]; ok {
		switch b := v.(type) {
		case bool:
			return b
		case string:
			// tolerate string "true"/"false" coming from YAML/JSON/Hugo
			return b == "true" || b == "1" || b == "yes"
		}
	}
	return false
}

func mapInt(m map[string]any, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		case string:
			// tolerate numeric strings from YAML/JSON; Atoi failure
			// (non-numeric input) returns 0, which is the documented
			// "absent" sentinel for this helper.
			i, _ := strconv.Atoi(n) //nolint:errcheck // Atoi failure (non-numeric input) returns 0, the documented "absent" sentinel
			return i
		}
	}
	return 0
}
