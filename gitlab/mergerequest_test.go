package gitlab

import "testing"

func TestMergeRequest(t *testing.T) {
	tests := []struct {
		name    string
		params  MergeRequestParams
		want    string
		wantErr bool
	}{
		{
			name: "minimal: only baseURL and title",
			params: MergeRequestParams{
				BaseURL: "https://gitlab.com/g/p",
				Title:   "Refactor login",
			},
			want: "https://gitlab.com/g/p/-/merge_requests/new?merge_request[title]=Refactor%20login",
		},
		{
			name: "with source and target branches",
			params: MergeRequestParams{
				BaseURL:      "https://gitlab.com/g/p",
				Title:        "Refactor login",
				SourceBranch: "feature/login",
				TargetBranch: "main",
			},
			want: "https://gitlab.com/g/p/-/merge_requests/new?merge_request[title]=Refactor%20login&merge_request[source_branch]=feature%2Flogin&merge_request[target_branch]=main",
		},
		{
			name: "all fields populated",
			params: MergeRequestParams{
				BaseURL:            "https://gitlab.com/g/p",
				Title:              "Refactor login",
				Description:        "See the design doc.",
				SourceBranch:       "feature/login",
				TargetBranch:       "main",
				RemoveSourceBranch: true,
				Squash:             true,
				NavSource:          "docs",
			},
			want: "https://gitlab.com/g/p/-/merge_requests/new?merge_request[title]=Refactor%20login&merge_request[description]=See%20the%20design%20doc.&merge_request[source_branch]=feature%2Flogin&merge_request[target_branch]=main&merge_request[remove_source_branch]=true&merge_request[squash]=true&nav_source=docs",
		},
		{
			name: "only source branch, no target",
			params: MergeRequestParams{
				BaseURL:      "https://gitlab.com/g/p",
				Title:        "WIP",
				SourceBranch: "feature/x",
			},
			want: "https://gitlab.com/g/p/-/merge_requests/new?merge_request[title]=WIP&merge_request[source_branch]=feature%2Fx",
		},
		{
			name: "trailing slash on baseURL normalised",
			params: MergeRequestParams{
				BaseURL: "https://gitlab.com/g/p/",
				Title:   "Foo",
			},
			want: "https://gitlab.com/g/p/-/merge_requests/new?merge_request[title]=Foo",
		},
		{
			name:    "missing baseURL is an error",
			params:  MergeRequestParams{Title: "Foo"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeRequest(tt.params)
			if (err != nil) != tt.wantErr {
				t.Fatalf("MergeRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("MergeRequest() =\n  %q\nwant:\n  %q", got, tt.want)
			}
		})
	}
}
