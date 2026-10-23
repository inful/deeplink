// Package gitlab builds deep links into GitLab: pre-filled work items
// (issues, incidents, merge requests) and other project-scoped views.
//
// All functions in this package are safe for concurrent use and have no
// state. Errors are returned for invalid input (e.g. a missing base URL);
// all other inputs are passed through to GitLab as-is, since the
// authoritative validation happens when the user actually opens the link.
package gitlab

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/inful/deeplink/internal/q"
)

// IssueType identifies the kind of work item to create via a GitLab
// "new work item" deep link.
type IssueType string

const (
	// IssueTypeIssue creates a regular issue (the default when IssueType is empty).
	IssueTypeIssue IssueType = "issue"
	// IssueTypeIncident creates an incident.
	IssueTypeIncident IssueType = "incident"
)

// IssueParams holds the fields for a pre-filled GitLab issue deep link.
// See the GitLab documentation for the underlying URL parameters:
// https://docs.gitlab.com/user/project/issues/create_issues/#using-a-url-with-prefilled-values
type IssueParams struct {
	// BaseURL is the project (or group) URL, e.g.
	// "https://gitlab.com/my-group/my-project". Required. Trailing
	// slashes are tolerated and normalised.
	BaseURL string

	// Title is the issue title. Optional at the library level — GitLab
	// will let the user fill it in if blank.
	Title string

	// Description is the issue body. GitLab-flavored Markdown is supported.
	Description string

	// IssueType is "issue" or "incident". Empty defaults to "issue" on
	// the GitLab side, so we simply omit the parameter when empty.
	IssueType IssueType

	// Confidential marks the new issue as confidential.
	Confidential bool

	// Template is the name of a project description template (without
	// the ".md" extension). Maps to the `issuable_template` parameter.
	Template string

	// DescriptionTemplate is the newer `description_template` parameter
	// (introduced in GitLab 17.9 for tasks, OKRs, issues, and epics).
	// Use Template (issuable_template) unless you specifically need this.
	DescriptionTemplate string

	// AddRelatedIssue, if positive, links the new issue to the existing
	// issue with the given numeric ID. The user still has to confirm the
	// relationship in the GitLab UI.
	AddRelatedIssue int
}

// Issue returns a URL that opens GitLab's "new work item" page with the
// given fields pre-filled. The BaseURL must be a project or group URL;
// the function appends "/-/work_items/new" and the appropriate query
// parameters.
//
// Example:
//
//	url, err := gitlab.Issue(gitlab.IssueParams{
//	    BaseURL:     "https://gitlab.com/my-group/my-project",
//	    Title:       "Login is broken",
//	    Description: "Steps to reproduce:\n1. ...",
//	})
//	// https://gitlab.com/my-group/my-project/-/work_items/new?issue[title]=Login%20is%20broken&issue[description]=Steps%20to%20reproduce%3A%0A1.%20...
func Issue(p IssueParams) (string, error) {
	if strings.TrimSpace(p.BaseURL) == "" {
		return "", errors.New("gitlab.Issue: BaseURL is required")
	}
	switch p.IssueType {
	case "", IssueTypeIssue, IssueTypeIncident:
	default:
		return "", fmt.Errorf("gitlab.Issue: invalid IssueType %q (want %q or %q)",
			p.IssueType, IssueTypeIssue, IssueTypeIncident)
	}

	base := q.JoinPath(p.BaseURL, "/-/work_items/new")

	pairs := []q.Pair{
		{Key: "issue[title]", Value: p.Title},
		{Key: "issue[description]", Value: p.Description},
	}
	if p.IssueType != "" {
		pairs = append(pairs, q.Pair{
			Key: "issue[issue_type]", Value: string(p.IssueType),
		})
	}
	if p.Confidential {
		pairs = append(pairs, q.Pair{
			Key: "issue[confidential]", Value: "true",
		})
	}
	if p.Template != "" {
		pairs = append(pairs, q.Pair{
			Key: "issuable_template", Value: p.Template,
		})
	}
	if p.DescriptionTemplate != "" {
		pairs = append(pairs, q.Pair{
			Key: "description_template", Value: p.DescriptionTemplate,
		})
	}
	if p.AddRelatedIssue > 0 {
		pairs = append(pairs, q.Pair{
			Key: "add_related_issue", Value: strconv.Itoa(p.AddRelatedIssue),
		})
	}

	qs, err := q.Encode(pairs)
	if err != nil {
		return "", fmt.Errorf("gitlab.Issue: %w", err)
	}
	if qs == "" {
		return base, nil
	}
	return base + "?" + qs, nil
}
