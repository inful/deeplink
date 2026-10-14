package gitlab

import "testing"

func TestIssue(t *testing.T) {
	tests := []struct {
		name    string
		params  IssueParams
		want    string
		wantErr bool
	}{
		{
			name: "minimal: only baseURL and title",
			params: IssueParams{
				BaseURL: "https://gitlab.com/my-group/my-project",
				Title:   "Login is broken",
			},
			want: "https://gitlab.com/my-group/my-project/-/work_items/new?issue[title]=Login%20is%20broken",
		},
		{
			name: "trailing slash on baseURL is normalised",
			params: IssueParams{
				BaseURL: "https://gitlab.com/my-group/my-project/",
				Title:   "Foo",
			},
			want: "https://gitlab.com/my-group/my-project/-/work_items/new?issue[title]=Foo",
		},
		{
			name: "all optional fields populated (canonical GitLab example)",
			params: IssueParams{
				BaseURL:     "https://gitlab.com/gitlab-org/gitlab",
				Title:       "Whoa, we're half-way there",
				Description: "Whoa, livin' in a URL",
				IssueType:   IssueTypeIssue,
			},
			want: "https://gitlab.com/gitlab-org/gitlab/-/work_items/new?issue[title]=Whoa%2C%20we%27re%20half-way%20there&issue[description]=Whoa%2C%20livin%27%20in%20a%20URL&issue[issue_type]=issue",
		},
		{
			name: "incident type",
			params: IssueParams{
				BaseURL:   "https://gitlab.com/g/p",
				Title:     "Production down",
				IssueType: IssueTypeIncident,
			},
			want: "https://gitlab.com/g/p/-/work_items/new?issue[title]=Production%20down&issue[issue_type]=incident",
		},
		{
			name: "confidential flag",
			params: IssueParams{
				BaseURL:      "https://gitlab.com/g/p",
				Title:        "Secret",
				Confidential: true,
			},
			want: "https://gitlab.com/g/p/-/work_items/new?issue[title]=Secret&issue[confidential]=true",
		},
		{
			name: "issuable_template",
			params: IssueParams{
				BaseURL:  "https://gitlab.com/g/p",
				Title:    "Validate new concept",
				Template: "Feature Proposal - basic",
			},
			want: "https://gitlab.com/g/p/-/work_items/new?issue[title]=Validate%20new%20concept&issuable_template=Feature%20Proposal%20-%20basic",
		},
		{
			name: "description_template (GitLab 17.9+)",
			params: IssueParams{
				BaseURL:             "https://gitlab.com/g/p",
				Title:               "Use the new template",
				DescriptionTemplate: "bug",
			},
			want: "https://gitlab.com/g/p/-/work_items/new?issue[title]=Use%20the%20new%20template&description_template=bug",
		},
		{
			name: "add_related_issue",
			params: IssueParams{
				BaseURL:         "https://gitlab.com/g/p",
				Title:           "Related to #123",
				AddRelatedIssue: 123,
			},
			want: "https://gitlab.com/g/p/-/work_items/new?issue[title]=Related%20to%20%23123&add_related_issue=123",
		},
		{
			name: "zero add_related_issue is omitted",
			params: IssueParams{
				BaseURL:         "https://gitlab.com/g/p",
				Title:           "No relation",
				AddRelatedIssue: 0,
			},
			want: "https://gitlab.com/g/p/-/work_items/new?issue[title]=No%20relation",
		},
		{
			name: "empty title still produces a valid URL (user might fill it in)",
			params: IssueParams{
				BaseURL: "https://gitlab.com/g/p",
			},
			want: "https://gitlab.com/g/p/-/work_items/new",
		},
		{
			name: "self-hosted GitLab (custom host)",
			params: IssueParams{
				BaseURL: "https://gitlab.example.com/team/board",
				Title:   "Custom host",
			},
			want: "https://gitlab.example.com/team/board/-/work_items/new?issue[title]=Custom%20host",
		},
		{
			name:    "missing baseURL is an error",
			params:  IssueParams{Title: "Foo"},
			wantErr: true,
		},
		{
			name: "empty baseURL (whitespace) is an error",
			params: IssueParams{
				BaseURL: "   ",
				Title:   "Foo",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Issue(tt.params)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Issue() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("Issue() =\n  %q\nwant:\n  %q", got, tt.want)
			}
		})
	}
}
