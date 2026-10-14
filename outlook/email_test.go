package outlook

import "testing"

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		params  EmailParams
		want    string
		wantErr bool
	}{
		{
			name: "single recipient, subject and body",
			params: EmailParams{
				To:      []string{"alice@example.com"},
				Subject: "Question about the rollout",
				Body:    "Hi Alice,\n\nQuick question...",
			},
			want: "mailto:alice@example.com?subject=Question%20about%20the%20rollout&body=Hi%20Alice%2C%0A%0AQuick%20question...",
		},
		{
			name: "multiple To recipients",
			params: EmailParams{
				To: []string{"a@x.com", "b@x.com", "c@x.com"},
			},
			want: "mailto:a@x.com,b@x.com,c@x.com",
		},
		{
			name: "To, CC, BCC, subject, body",
			params: EmailParams{
				To:      []string{"a@x.com", "b@x.com"},
				CC:      []string{"c@x.com"},
				BCC:     []string{"d@x.com", "e@x.com"},
				Subject: "Group email",
				Body:    "FYI",
			},
			want: "mailto:a@x.com,b@x.com?subject=Group%20email&body=FYI&cc=c@x.com&bcc=d@x.com,e@x.com",
		},
		{
			name: "only CC, no To (unusual but valid)",
			params: EmailParams{
				CC: []string{"only@x.com"},
			},
			want: "mailto:?cc=only@x.com",
		},
		{
			name: "only BCC",
			params: EmailParams{
				BCC: []string{"secret@x.com"},
			},
			want: "mailto:?bcc=secret@x.com",
		},
		{
			name: "multi-line body preserved",
			params: EmailParams{
				To:   []string{"a@x.com"},
				Body: "Line 1\nLine 2\n\nLine 4",
			},
			want: "mailto:a@x.com?body=Line%201%0ALine%202%0A%0ALine%204",
		},
		{
			name:    "no recipients at all is an error",
			params:  EmailParams{Subject: "Hello"},
			wantErr: true,
		},
		{
			name:    "all recipient slices contain only empty strings is an error",
			params:  EmailParams{To: []string{""}, CC: []string{""}, BCC: []string{""}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Email(tt.params)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Email() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("Email() =\n  %q\nwant:\n  %q", got, tt.want)
			}
		})
	}
}
