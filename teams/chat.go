// Package teams builds deep links into Microsoft Teams: starting a new
// chat, navigating to a chat, team, or channel, and so on.
//
// All functions in this package are safe for concurrent use and have no
// state. Errors are returned only for invalid input that prevents a URL
// from being formed (e.g. a chat with no participants); other inputs are
// passed through to Teams as-is, since the authoritative validation
// happens when the user actually opens the link.
//
// Reference: https://learn.microsoft.com/en-us/microsoftteams/platform/concepts/build-and-test/deep-link-teams
package teams

import (
	"errors"
	"fmt"
	"strings"

	"github.com/inful/deeplink/internal/q"
)

// ChatParams holds the fields for a Teams "start a new chat" deep link.
//
// All fields are optional except Users, which must contain at least one
// non-empty UPN (typically an email address) or a bot ID prefixed with
// "28:". A one-on-one chat is created when exactly one user is supplied;
// a group chat is created with two or more.
type ChatParams struct {
	// TenantID is the Microsoft Entra tenant ID. Optional but recommended
	// — when omitted, the user may be prompted to pick a tenant.
	TenantID string

	// Users is the list of participants. Each entry should be a User
	// Principal Name (email address) or a bot app ID prefixed with "28:".
	// At least one non-empty entry is required.
	Users []string

	// TopicName is the chat's display name. Only meaningful for chats
	// with three or more participants; ignored otherwise.
	TopicName string

	// Message is pre-populated text inserted into the compose box. The
	// user still has to press send.
	Message string
}

// NewChat returns a URL that opens Teams to a "new chat" compose window
// with the given participants (and an optional topic and prepopulated
// message). If a chat with the same participants already exists, Teams
// opens that existing chat instead of creating a new one.
//
// Example:
//
//	url, err := teams.NewChat(teams.ChatParams{
//	    Users:   []string{"joe@contoso.com", "bob@contoso.com"},
//	    Message: "Quick question about the rollout",
//	})
func NewChat(p ChatParams) (string, error) {
	// Drop empty entries; require at least one real user.
	users := make([]string, 0, len(p.Users))
	for _, u := range p.Users {
		if trimmed := strings.TrimSpace(u); trimmed != "" {
			users = append(users, trimmed)
		}
	}
	if len(users) == 0 {
		return "", errors.New("teams.NewChat: at least one user is required")
	}

	const base = "https://teams.microsoft.com/l/chat/0/0"
	pairs := []q.Pair{
		{Key: "tenantId", Value: p.TenantID},
		// Teams expects the users value to be a literal comma-separated
		// list (per Microsoft's documentation), not a single value with
		// "%2C" between the addresses. Use Raw to suppress encoding;
		// q.Encode still rejects values that could escape this field.
		{Key: "users", Value: strings.Join(users, ","), Raw: true},
		{Key: "topicName", Value: p.TopicName},
		{Key: "message", Value: p.Message},
	}
	qs, err := q.Encode(pairs)
	if err != nil {
		return "", fmt.Errorf("teams.NewChat: %w", err)
	}
	if qs == "" {
		return base, nil
	}
	return base + "?" + qs, nil
}
