package teams

import (
	"text/template"

	"github.com/inful/deeplink/internal/tmpl"
)

// FuncMap returns a template.FuncMap containing every Teams deep-link
// builder in this package, namespaced under the "teams" prefix.
//
// Available functions:
//
//	teamsNewChat — teams.ChatParams    → URL for starting a new chat
//	teamsChannel — teams.ChannelParams → URL for opening an existing channel
//
// Example (Hugo):
//
//	{{ $u := teamsNewChat (dict
//	    "Users"   (slice "joe@contoso.com" "bob@contoso.com")
//	    "Message" "Quick question about the rollout"
//	) }}
//	<a href="{{ $u }}">Open in Teams</a>
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"teamsNewChat": newChatFromMap,
		"teamsChannel": channelFromMap,
	}
}

// newChatFromMap is the template-friendly adapter for NewChat.
func newChatFromMap(m map[string]any) (string, error) {
	return NewChat(ChatParams{
		TenantID:  tmpl.String(m, "TenantID"),
		Users:     tmpl.StringList(m["Users"]),
		TopicName: tmpl.String(m, "TopicName"),
		Message:   tmpl.String(m, "Message"),
	})
}

// channelFromMap is the template-friendly adapter for Channel.
func channelFromMap(m map[string]any) (string, error) {
	return Channel(ChannelParams{
		TenantID:    tmpl.String(m, "TenantID"),
		GroupID:     tmpl.String(m, "GroupID"),
		ChannelID:   tmpl.String(m, "ChannelID"),
		ChannelName: tmpl.String(m, "ChannelName"),
		Private:     tmpl.Bool(m, "Private"),
		Shared:      tmpl.Bool(m, "Shared"),
	})
}
