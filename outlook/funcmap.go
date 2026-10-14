package outlook

import "text/template"

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
		Subject:   mapString(m, "Subject"),
		Body:      mapString(m, "Body"),
		Location:  mapString(m, "Location"),
		StartDate: mapString(m, "StartDate"),
		EndDate:   mapString(m, "EndDate"),
	})
}

func mapString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// stringList extracts a []string from a map value, accepting both
// []string and []any (the latter is common when reading from Hugo's
// .Site.Params or frontmatter).
func stringList(v any) []string {
	switch xs := v.(type) {
	case []string:
		return xs
	case []any:
		out := make([]string, 0, len(xs))
		for _, x := range xs {
			if s, ok := x.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// emailFromMap is the template-friendly adapter for Email.
func emailFromMap(m map[string]any) (string, error) {
	return Email(EmailParams{
		To:      stringList(m["To"]),
		CC:      stringList(m["CC"]),
		BCC:     stringList(m["BCC"]),
		Subject: mapString(m, "Subject"),
		Body:    mapString(m, "Body"),
	})
}
