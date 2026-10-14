package teams

import "text/template"

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
	users, _ := m["Users"].([]string) //nolint:errcheck // nil slice is the expected "absent" value, handled below
	// Hugo users may also pass a []any (from .Site.Params or frontmatter);
	// accept that too.
	if users == nil {
		if rawList, ok := m["Users"].([]any); ok {
			users = make([]string, 0, len(rawList))
			for _, v := range rawList {
				if s, ok := v.(string); ok {
					users = append(users, s)
				}
			}
		}
	}
	return NewChat(ChatParams{
		TenantID:  mapString(m, "TenantID"),
		Users:     users,
		TopicName: mapString(m, "TopicName"),
		Message:   mapString(m, "Message"),
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

func mapBool(m map[string]any, key string) bool {
	if v, ok := m[key]; ok {
		switch b := v.(type) {
		case bool:
			return b
		case string:
			return b == "true" || b == "1" || b == "yes"
		}
	}
	return false
}

// channelFromMap is the template-friendly adapter for Channel.
func channelFromMap(m map[string]any) (string, error) {
	return Channel(ChannelParams{
		TenantID:    mapString(m, "TenantID"),
		GroupID:     mapString(m, "GroupID"),
		ChannelID:   mapString(m, "ChannelID"),
		ChannelName: mapString(m, "ChannelName"),
		Private:     mapBool(m, "Private"),
		Shared:      mapBool(m, "Shared"),
	})
}
