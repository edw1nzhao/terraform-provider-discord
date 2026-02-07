package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateBanParams are the parameters for creating a guild ban.
type CreateBanParams struct {
	DeleteMessageDays    *int    `json:"delete_message_days,omitempty"`
	DeleteMessageSeconds *int    `json:"delete_message_seconds,omitempty"`
}

// GetGuildBans returns a list of ban objects for the guild.
func (c *Client) GetGuildBans(ctx context.Context, guildID Snowflake) ([]*Ban, error) {
	var bans []*Ban
	route := fmt.Sprintf("/guilds/%s/bans", guildID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, &bans)
	if err != nil {
		return nil, err
	}
	return bans, nil
}

// GetGuildBan returns the ban object for the given user.
func (c *Client) GetGuildBan(ctx context.Context, guildID Snowflake, userID Snowflake) (*Ban, error) {
	ban := new(Ban)
	route := fmt.Sprintf("/guilds/%s/bans/%s", guildID, userID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, ban)
	if err != nil {
		return nil, err
	}
	return ban, nil
}

// CreateGuildBan creates a guild ban and optionally deletes previous messages.
func (c *Client) CreateGuildBan(ctx context.Context, guildID Snowflake, userID Snowflake, params *CreateBanParams) error {
	route := fmt.Sprintf("/guilds/%s/bans/%s", guildID, userID)
	return c.doRequestNoContent(ctx, http.MethodPut, route, params)
}

// RemoveGuildBan removes the ban for a user.
func (c *Client) RemoveGuildBan(ctx context.Context, guildID Snowflake, userID Snowflake) error {
	route := fmt.Sprintf("/guilds/%s/bans/%s", guildID, userID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}
