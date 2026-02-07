package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateInviteParams are the parameters for creating a channel invite.
type CreateInviteParams struct {
	MaxAge              *int       `json:"max_age,omitempty"`
	MaxUses             *int       `json:"max_uses,omitempty"`
	Temporary           *bool      `json:"temporary,omitempty"`
	Unique              *bool      `json:"unique,omitempty"`
	TargetType          *int       `json:"target_type,omitempty"`
	TargetUserID        *Snowflake `json:"target_user_id,omitempty"`
	TargetApplicationID *Snowflake `json:"target_application_id,omitempty"`
}

// CreateChannelInvite creates a new invite for a channel.
func (c *Client) CreateChannelInvite(ctx context.Context, channelID Snowflake, params *CreateInviteParams) (*Invite, error) {
	invite := new(Invite)
	route := fmt.Sprintf("/channels/%s/invites", channelID)
	err := c.doRequest(ctx, http.MethodPost, route, params, invite)
	if err != nil {
		return nil, err
	}
	return invite, nil
}

// GetInvite returns an invite object for the given code.
func (c *Client) GetInvite(ctx context.Context, code string) (*Invite, error) {
	invite := new(Invite)
	route := fmt.Sprintf("/invites/%s", code)
	err := c.doRequest(ctx, http.MethodGet, route, nil, invite)
	if err != nil {
		return nil, err
	}
	return invite, nil
}

// DeleteInvite deletes an invite. Returns the deleted invite object.
func (c *Client) DeleteInvite(ctx context.Context, code string) (*Invite, error) {
	invite := new(Invite)
	route := fmt.Sprintf("/invites/%s", code)
	err := c.doRequest(ctx, http.MethodDelete, route, nil, invite)
	if err != nil {
		return nil, err
	}
	return invite, nil
}

// GetGuildInvites returns a list of invites for the guild.
func (c *Client) GetGuildInvites(ctx context.Context, guildID Snowflake) ([]*Invite, error) {
	var invites []*Invite
	route := fmt.Sprintf("/guilds/%s/invites", guildID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, &invites)
	if err != nil {
		return nil, err
	}
	return invites, nil
}
