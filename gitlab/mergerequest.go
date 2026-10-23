package gitlab

import (
	"errors"
	"fmt"
	"strings"

	"github.com/inful/deeplink/internal/q"
)

// MergeRequestParams holds the fields for a pre-filled GitLab merge
// request deep link.
//
// GitLab's "new merge request" UI requires the user to pick a source and
// target branch before showing the rest of the form, so a deep link
// without both branches opens an empty picker. To produce a useful
// pre-filled form, set SourceBranch and TargetBranch.
type MergeRequestParams struct {
	// BaseURL is the project URL, e.g. "https://gitlab.com/my-group/my-project".
	// Required.
	BaseURL string

	// Title is the merge request title.
	Title string

	// Description is the merge request body. GitLab-flavored Markdown is supported.
	Description string

	// SourceBranch is the branch containing the changes (e.g. "feature/login").
	// When set, the GitLab UI pre-selects this branch.
	SourceBranch string

	// TargetBranch is the branch the changes will be merged into (e.g. "main").
	// When set, the GitLab UI pre-selects this branch.
	TargetBranch string

	// RemoveSourceBranch tells GitLab to delete the source branch after the
	// merge request is merged.
	RemoveSourceBranch bool

	// Squash tells GitLab to squash commits on merge.
	Squash bool

	// NavSource is a free-form label that GitLab uses to attribute the
	// merge request to a navigation source (e.g. "docs" for links
	// originating from documentation).
	NavSource string
}

// MergeRequest returns a URL that opens GitLab's "new merge request" page
// with the given fields pre-filled. The function appends "/-/merge_requests/new"
// to BaseURL and the appropriate query parameters.
//
// Example:
//
//	url, err := gitlab.MergeRequest(gitlab.MergeRequestParams{
//	    BaseURL:      "https://gitlab.com/g/p",
//	    Title:        "Refactor login",
//	    SourceBranch: "feature/login",
//	    TargetBranch: "main",
//	})
func MergeRequest(p MergeRequestParams) (string, error) {
	if strings.TrimSpace(p.BaseURL) == "" {
		return "", errors.New("gitlab.MergeRequest: BaseURL is required")
	}

	base := q.JoinPath(p.BaseURL, "/-/merge_requests/new")

	pairs := []q.Pair{
		{Key: "merge_request[title]", Value: p.Title},
		{Key: "merge_request[description]", Value: p.Description},
		{Key: "merge_request[source_branch]", Value: p.SourceBranch},
		{Key: "merge_request[target_branch]", Value: p.TargetBranch},
	}
	if p.RemoveSourceBranch {
		pairs = append(pairs, q.Pair{
			Key: "merge_request[remove_source_branch]", Value: "true",
		})
	}
	if p.Squash {
		pairs = append(pairs, q.Pair{
			Key: "merge_request[squash]", Value: "true",
		})
	}
	if p.NavSource != "" {
		pairs = append(pairs, q.Pair{
			Key: "nav_source", Value: p.NavSource,
		})
	}

	qs, err := q.Encode(pairs)
	if err != nil {
		return "", fmt.Errorf("gitlab.MergeRequest: %w", err)
	}
	if qs == "" {
		return base, nil
	}
	return base + "?" + qs, nil
}
