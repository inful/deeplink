package teams

import "testing"

func TestChannel(t *testing.T) {
	tests := []struct {
		name    string
		params  ChannelParams
		want    string
		wantErr bool
	}{
		{
			name: "standard channel (canonical Microsoft example)",
			params: ChannelParams{
				ChannelID:   "19:9be3de4e70874c71a608dee9ba803ed3@thread.tacv2",
				ChannelName: "My example channel",
				GroupID:     "72602e12-78ac-474c-99d6-f619710353a9",
				TenantID:    "aaaabbbb-0000-cccc-1111-dddd2222eeee",
			},
			want: "https://teams.microsoft.com/l/channel/19%3A9be3de4e70874c71a608dee9ba803ed3%40thread.tacv2/My%20example%20channel?groupId=72602e12-78ac-474c-99d6-f619710353a9&tenantId=aaaabbbb-0000-cccc-1111-dddd2222eeee",
		},
		{
			name: "private channel (ngc=true)",
			params: ChannelParams{
				ChannelID:   "19:abc@thread.tacv2",
				ChannelName: "Secret",
				GroupID:     "g-1",
				TenantID:    "t-1",
				Private:     true,
			},
			want: "https://teams.microsoft.com/l/channel/19%3Aabc%40thread.tacv2/Secret?groupId=g-1&tenantId=t-1&ngc=true",
		},
		{
			name: "shared channel (allowXTenantAccess=true)",
			params: ChannelParams{
				ChannelID:   "19:xyz@thread.tacv2",
				ChannelName: "Shared",
				GroupID:     "g-2",
				TenantID:    "t-2",
				Shared:      true,
			},
			want: "https://teams.microsoft.com/l/channel/19%3Axyz%40thread.tacv2/Shared?groupId=g-2&tenantId=t-2&allowXTenantAccess=true",
		},
		{
			name: "channel name with special characters",
			params: ChannelParams{
				ChannelID:   "19:id@thread.tacv2",
				ChannelName: "R&D / Engineering",
				GroupID:     "g-3",
				TenantID:    "t-3",
			},
			want: "https://teams.microsoft.com/l/channel/19%3Aid%40thread.tacv2/R%26D%20%2F%20Engineering?groupId=g-3&tenantId=t-3",
		},
		{
			name: "minimal: channelId and name only",
			params: ChannelParams{
				ChannelID:   "19:min@thread.tacv2",
				ChannelName: "Min",
			},
			want: "https://teams.microsoft.com/l/channel/19%3Amin%40thread.tacv2/Min",
		},
		{
			name: "empty channel name leaves no trailing slash",
			params: ChannelParams{
				ChannelID: "19:min@thread.tacv2",
			},
			want: "https://teams.microsoft.com/l/channel/19%3Amin%40thread.tacv2",
		},
		{
			name:    "missing channelId is an error",
			params:  ChannelParams{ChannelName: "Foo"},
			wantErr: true,
		},
		{
			name:    "whitespace-only channelId is an error",
			params:  ChannelParams{ChannelID: "   ", ChannelName: "Foo"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Channel(tt.params)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Channel() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("Channel() =\n  %q\nwant:\n  %q", got, tt.want)
			}
		})
	}
}
