package discord

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ModifyMemberParams are the parameters for modifying a guild member.
type ModifyMemberParams struct {
	Nick                       *string    `json:"nick,omitempty"`
	Roles                      []Snowflake `json:"roles,omitempty"`
	Mute                       *bool      `json:"mute,omitempty"`
	Deaf                       *bool      `json:"deaf,omitempty"`
	ChannelID                  *Snowflake `json:"channel_id,omitempty"`
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until,omitempty"`
	Flags                      *int       `json:"flags,omitempty"`
}

// GetGuildMember returns a guild member for the given user ID.
func (c *Client) GetGuildMember(ctx context.Context, guildID Snowflake, userID Snowflake) (*Member, error) {
	member := new(Member)
	route := fmt.Sprintf("/guilds/%s/members/%s", guildID, userID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, member)
	if err != nil {
		return nil, err
	}
	return member, nil
}

// ModifyGuildMember modifies attributes of a guild member.
func (c *Client) ModifyGuildMember(ctx context.Context, guildID Snowflake, userID Snowflake, params *ModifyMemberParams) (*Member, error) {
	member := new(Member)
	route := fmt.Sprintf("/guilds/%s/members/%s", guildID, userID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, member)
	if err != nil {
		return nil, err
	}
	return member, nil
}

// RemoveGuildMember removes a member from a guild (kick).
func (c *Client) RemoveGuildMember(ctx context.Context, guildID Snowflake, userID Snowflake) error {
	route := fmt.Sprintf("/guilds/%s/members/%s", guildID, userID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}

// AddGuildMemberRole adds a role to a guild member.
func (c *Client) AddGuildMemberRole(ctx context.Context, guildID Snowflake, userID Snowflake, roleID Snowflake) error {
	route := fmt.Sprintf("/guilds/%s/members/%s/roles/%s", guildID, userID, roleID)
	return c.doRequestNoContent(ctx, http.MethodPut, route, nil)
}

// RemoveGuildMemberRole removes a role from a guild member.
func (c *Client) RemoveGuildMemberRole(ctx context.Context, guildID Snowflake, userID Snowflake, roleID Snowflake) error {
	route := fmt.Sprintf("/guilds/%s/members/%s/roles/%s", guildID, userID, roleID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}
