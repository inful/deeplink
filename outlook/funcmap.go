package outlook

import (
	"text/template"

	"github.com/inful/deeplink/internal/tmpl"
)

// FuncMap returns a template.FuncMap containing every Outlook deep-link
// builder in this package, namespaced under the "outlook" prefix.
//
// Available functions:
//
//	outlookEvent — outlook.EventParams → URL for composing a new calendar event
//	outlookEmail — outlook.EmailParams → mailto: URL for composing a new email
//
// Example (Hugo):
//
//	{{ $u := outlookEvent (dict
//	    "Subject"   "Sprint review"
//	    "StartDate" "2024-01-15T10:00:00.000"
//	    "EndDate"   "2024-01-15T11:00:00.000"
//	    "Location"  "Teams"
//	) }}
//	<a href="{{ $u }}">Add to calendar</a>
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"outlookEvent": eventFromMap,
		"outlookEmail": emailFromMap,
	}
}

// eventFromMap is the template-friendly adapter for Event.
func eventFromMap(m map[string]any) (string, error) {
	return Event(EventParams{
		Subject:   tmpl.String(m, "Subject"),
		Body:      tmpl.String(m, "Body"),
		Location:  tmpl.String(m, "Location"),
		StartDate: tmpl.String(m, "StartDate"),
		EndDate:   tmpl.String(m, "EndDate"),
	})
}

// emailFromMap is the template-friendly adapter for Email.
func emailFromMap(m map[string]any) (string, error) {
	return Email(EmailParams{
		To:      tmpl.StringList(m["To"]),
		CC:      tmpl.StringList(m["CC"]),
		BCC:     tmpl.StringList(m["BCC"]),
		Subject: tmpl.String(m, "Subject"),
		Body:    tmpl.String(m, "Body"),
	})
}
