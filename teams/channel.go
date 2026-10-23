package teams

import (
	"errors"
	"fmt"
	"strings"

	"github.com/inful/deeplink/internal/q"
)

// ChannelParams holds the fields for a Teams "navigate to a channel" deep
// link. This is a navigation link (open an existing channel), not a
// channel-creation link — Teams has no documented deep link for creating
// channels.
type ChannelParams struct {
	// TenantID is the Microsoft Entra tenant ID.
	TenantID string

	// GroupID is the Microsoft 365 group ID for the team hosting the channel.
	GroupID string

	// ChannelID is the channel's identifier, typically in the form
	// "19:<GUID>@thread.tacv2". Required.
	ChannelID string

	// ChannelName is the human-readable channel name, used to label the
	// URL path (e.g. "My example channel").
	ChannelName string

	// Private marks the channel as a private (next-generation) channel;
	// appends "ngc=true" to the URL.
	Private bool

	// Shared marks the channel as a shared (cross-tenant) channel;
	// appends "allowXTenantAccess=true" to the URL.
	Shared bool
}

// Channel returns a URL that opens the given Teams channel. Both standard
// and private/shared channels are supported via the Private and Shared
// flags.
//
// The channelId and channelName are encoded using query-style escaping
// (matching the format shown in Microsoft's documentation, e.g.
// "19%3A...%40thread.tacv2"). Both encodings are valid and Teams
// accepts either, but the query-style form is the one Microsoft
// documents.
//
// Example:
//
//	url, err := teams.Channel(teams.ChannelParams{
//	    ChannelID:   "19:9be3de4e70874c71a608dee9ba803ed3@thread.tacv2",
//	    ChannelName: "General",
//	    GroupID:     "72602e12-78ac-474c-99d6-f619710353a9",
//	    TenantID:    "aaaabbbb-0000-cccc-1111-dddd2222eeee",
//	})
func Channel(p ChannelParams) (string, error) {
	if strings.TrimSpace(p.ChannelID) == "" {
		return "", errors.New("teams.Channel: ChannelID is required")
	}

	base := q.JoinPath("https://teams.microsoft.com/l/channel",
		q.PathSegment(p.ChannelID), q.PathSegment(p.ChannelName))

	pairs := []q.Pair{
		{Key: "groupId", Value: p.GroupID},
		{Key: "tenantId", Value: p.TenantID},
	}
	if p.Private {
		pairs = append(pairs, q.Pair{Key: "ngc", Value: "true"})
	}
	if p.Shared {
		pairs = append(pairs, q.Pair{Key: "allowXTenantAccess", Value: "true"})
	}
	qs, err := q.Encode(pairs)
	if err != nil {
		return "", fmt.Errorf("teams.Channel: %w", err)
	}
	if qs == "" {
		return base, nil
	}
	return base + "?" + qs, nil
}
