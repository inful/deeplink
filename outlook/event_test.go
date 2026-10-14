package outlook

import "testing"

func TestEvent(t *testing.T) {
	tests := []struct {
		name    string
		params  EventParams
		want    string
		wantErr bool
	}{
		{
			name: "canonical example from Microsoft Q&A",
			params: EventParams{
				Subject:   "Test event",
				Body:      "Test description",
				StartDate: "2023-06-05T20:00:00.000",
				EndDate:   "2023-06-05T20:30:00.000",
			},
			want: "https://outlook.office.com/calendar/deeplink/compose?body=Test%20description&enddt=2023-06-05T20%3A30%3A00.000&startdt=2023-06-05T20%3A00%3A00.000&subject=Test%20event",
		},
		{
			name: "minimal: subject only",
			params: EventParams{
				Subject: "Quick sync",
			},
			want: "https://outlook.office.com/calendar/deeplink/compose?subject=Quick%20sync",
		},
		{
			name: "all-day event",
			params: EventParams{
				Subject:   "Company offsite",
				StartDate: "2024-09-15",
				EndDate:   "2024-09-16",
			},
			want: "https://outlook.office.com/calendar/deeplink/compose?enddt=2024-09-16&startdt=2024-09-15&subject=Company%20offsite",
		},
		{
			name: "with location",
			params: EventParams{
				Subject:   "Sprint review",
				StartDate: "2024-01-15T10:00:00.000",
				EndDate:   "2024-01-15T11:00:00.000",
				Location:  "Conf room 4 / Teams",
			},
			want: "https://outlook.office.com/calendar/deeplink/compose?enddt=2024-01-15T11%3A00%3A00.000&startdt=2024-01-15T10%3A00%3A00.000&subject=Sprint%20review&location=Conf%20room%204%20%2F%20Teams",
		},
		{
			name: "multi-line body preserved",
			params: EventParams{
				Subject:   "Standup",
				Body:      "Agenda:\n- Updates\n- Blockers",
				StartDate: "2024-01-15T10:00:00.000",
				EndDate:   "2024-01-15T10:15:00.000",
			},
			want: "https://outlook.office.com/calendar/deeplink/compose?body=Agenda%3A%0A-%20Updates%0A-%20Blockers&enddt=2024-01-15T10%3A15%3A00.000&startdt=2024-01-15T10%3A00%3A00.000&subject=Standup",
		},
		{
			name: "start only (Outlook may treat end as same time)",
			params: EventParams{
				Subject:   "Test",
				StartDate: "2024-01-15T10:00:00.000",
			},
			want: "https://outlook.office.com/calendar/deeplink/compose?startdt=2024-01-15T10%3A00%3A00.000&subject=Test",
		},
		{
			name:    "missing subject is an error",
			params:  EventParams{StartDate: "2024-01-15T10:00:00.000", EndDate: "2024-01-15T11:00:00.000"},
			wantErr: true,
		},
		{
			name:    "whitespace-only subject is an error",
			params:  EventParams{Subject: "   ", StartDate: "2024-01-15T10:00:00.000", EndDate: "2024-01-15T11:00:00.000"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Event(tt.params)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Event() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("Event() =\n  %q\nwant:\n  %q", got, tt.want)
			}
		})
	}
}
