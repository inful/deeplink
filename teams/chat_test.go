package teams

import "testing"

func TestNewChat(t *testing.T) {
	tests := []struct {
		name    string
		params  ChatParams
		want    string
		wantErr bool
	}{
		{
			name: "canonical example from Microsoft docs",
			params: ChatParams{
				TenantID:  "f4c2b8e1-7a3d-4e6f-9b12-8c5d0e3a1f47",
				Users:     []string{"joe@contoso.com", "bob@contoso.com"},
				TopicName: "Prep For Meeting Tomorrow",
				Message:   "Hi folks, kicking off a chat about our meeting tomorrow",
			},
			want: "https://teams.microsoft.com/l/chat/0/0?tenantId=f4c2b8e1-7a3d-4e6f-9b12-8c5d0e3a1f47&users=joe@contoso.com,bob@contoso.com&topicName=Prep%20For%20Meeting%20Tomorrow&message=Hi%20folks%2C%20kicking%20off%20a%20chat%20about%20our%20meeting%20tomorrow",
		},
		{
			name: "one user, no topic (1:1 chat)",
			params: ChatParams{
				Users: []string{"alice@contoso.com"},
			},
			want: "https://teams.microsoft.com/l/chat/0/0?users=alice@contoso.com",
		},
		{
			name: "with prepopulated message",
			params: ChatParams{
				Users:   []string{"alice@contoso.com"},
				Message: "Quick question",
			},
			want: "https://teams.microsoft.com/l/chat/0/0?users=alice@contoso.com&message=Quick%20question",
		},
		{
			name: "bot user with 28: prefix",
			params: ChatParams{
				Users:   []string{"28:47345678-2134-6534-9143-65146789012"},
				Message: "This message was triggered by a link!",
			},
			// RFC 3986-compliant encoding: '!' becomes %21. The Microsoft
			// docs show a literal '!' at the end of the value, which is
			// a stylistic choice in the documentation; both forms are
			// valid and Teams accepts either.
			want: "https://teams.microsoft.com/l/chat/0/0?users=28:47345678-2134-6534-9143-65146789012&message=This%20message%20was%20triggered%20by%20a%20link%21",
		},
		{
			name: "only tenant and users",
			params: ChatParams{
				TenantID: "0d9b645f-597b-41f0-a2a3-ef103fbd91bb",
				Users:    []string{"a@x.com", "b@x.com", "c@x.com"},
			},
			want: "https://teams.microsoft.com/l/chat/0/0?tenantId=0d9b645f-597b-41f0-a2a3-ef103fbd91bb&users=a@x.com,b@x.com,c@x.com",
		},
		{
			name:    "no users is an error (cannot chat with nobody)",
			params:  ChatParams{},
			wantErr: true,
		},
		{
			name:    "users slice with only empty strings is an error",
			params:  ChatParams{Users: []string{"", ""}},
			wantErr: true,
		},
		{
			name:    "explicit nil users is an error",
			params:  ChatParams{Users: nil},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChat(tt.params)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewChat() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("NewChat() =\n  %q\nwant:\n  %q", got, tt.want)
			}
		})
	}
}
