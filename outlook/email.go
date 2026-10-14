package outlook

import (
	"errors"
	"strings"

	"github.com/inful/deeplink/internal/q"
)

// EmailParams holds the fields for an Outlook "compose new email" deep
// link, expressed as a `mailto:` URL.
//
// Outlook's web client understands `mailto:` URLs just like every other
// mail client: it opens a compose window with the recipients, subject,
// and body pre-filled. The same URL works for users whose default mail
// handler is a desktop client; if you need to force the web client,
// prepend "https://outlook.office.com/mail/deeplink/compose?" and convert
// the recipient list to a "to" query parameter.
type EmailParams struct {
	// To is the list of primary recipients. At least one recipient is
	// required across To/CC/BCC combined.
	To []string

	// CC is the list of carbon-copy recipients.
	CC []string

	// BCC is the list of blind carbon-copy recipients.
	BCC []string

	// Subject is the email subject.
	Subject string

	// Body is the email body. Plain text; newlines are preserved.
	Body string
}

// Email returns a `mailto:` URL that opens the user's default mail
// client (Outlook, Apple Mail, Thunderbird, …) with the given fields
// pre-filled. Recipients in each category are joined with commas, as
// is standard for `mailto:` URLs.
//
// Example:
//
//	url := outlook.Email(outlook.EmailParams{
//	    To:      []string{"alice@example.com"},
//	    Subject: "Sprint planning",
//	    Body:    "Agenda: ...",
//	})
//	// mailto:alice@example.com?body=Agenda%3A%20...&subject=Sprint%20planning
func Email(p EmailParams) (string, error) {
	to := joinNonEmpty(p.To)
	cc := joinNonEmpty(p.CC)
	bcc := joinNonEmpty(p.BCC)
	if to == "" && cc == "" && bcc == "" {
		return "", errors.New("outlook.Email: at least one recipient is required (To, CC, or BCC)")
	}

	base := "mailto:" + to
	// The cc/bcc parameters, like the `to` field, are comma-separated
	// lists of addresses; suppress encoding so the commas survive.
	pairs := []q.Pair{
		{Key: "subject", Value: p.Subject},
		{Key: "body", Value: p.Body},
		{Key: "cc", Value: cc, Raw: true},
		{Key: "bcc", Value: bcc, Raw: true},
	}
	qs := q.Encode(pairs)
	if qs == "" {
		return base, nil
	}
	return base + "?" + qs, nil
}

// joinNonEmpty joins non-empty strings with commas.
func joinNonEmpty(parts []string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, ",")
}
