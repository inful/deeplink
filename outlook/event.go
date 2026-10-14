// Package outlook builds deep links into Microsoft Outlook: creating
// calendar events, composing email, and so on.
//
// All functions in this package are safe for concurrent use and have no
// state. Errors are returned only for invalid input that prevents a URL
// from being formed (e.g. a calendar event with no subject); other
// inputs are passed through to Outlook as-is, since the authoritative
// validation happens when the user actually opens the link.
package outlook

import (
	"errors"
	"strings"

	"github.com/inful/deeplink/internal/q"
)

// EventParams holds the fields for an Outlook "compose new calendar
// event" deep link.
//
// Date/time values are passed through as strings. Outlook accepts a wide
// range of formats; the most reliable is the ISO 8601 form
// "YYYY-MM-DDTHH:MM:SS.000" (e.g. "2024-01-15T10:00:00.000") used in
// Microsoft's own examples. For all-day events, a bare "YYYY-MM-DD"
// works for both StartDate and EndDate (where EndDate is the day AFTER
// the event — Outlook convention).
type EventParams struct {
	// Subject is the event title. Required.
	Subject string

	// Body is the event description. Plain text; newlines are preserved.
	Body string

	// Location is the event location (e.g. "Conf room 4" or a Teams link).
	Location string

	// StartDate is the event start in a format Outlook accepts. See the
	// type doc for the recommended form.
	StartDate string

	// EndDate is the event end in a format Outlook accepts.
	EndDate string
}

// Event returns a URL that opens Outlook on the web's "compose new
// calendar event" page with the given fields pre-filled.
//
// Example:
//
//	url, err := outlook.Event(outlook.EventParams{
//	    Subject:   "Sprint review",
//	    Body:      "Agenda: ...",
//	    StartDate: "2024-01-15T10:00:00.000",
//	    EndDate:   "2024-01-15T11:00:00.000",
//	    Location:  "Teams",
//	})
func Event(p EventParams) (string, error) {
	if strings.TrimSpace(p.Subject) == "" {
		return "", errors.New("outlook.Event: Subject is required")
	}

	const base = "https://outlook.office.com/calendar/deeplink/compose"
	// Order matches the canonical Microsoft example so output is
	// predictable and matches the documentation's URL layout.
	pairs := []q.Pair{
		{Key: "body", Value: p.Body},
		{Key: "enddt", Value: p.EndDate},
		{Key: "startdt", Value: p.StartDate},
		{Key: "subject", Value: p.Subject},
		{Key: "location", Value: p.Location},
	}
	qs := q.Encode(pairs)
	if qs == "" {
		return base, nil
	}
	return base + "?" + qs, nil
}
