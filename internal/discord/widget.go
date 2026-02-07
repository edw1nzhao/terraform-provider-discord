package discord

import (
	"context"
	"fmt"
	"net/http"
)

// ModifyWidgetParams are the parameters for modifying a guild widget.
type ModifyWidgetParams struct {
	Enabled   *bool      `json:"enabled,omitempty"`
	ChannelID *Snowflake `json:"channel_id,omitempty"`
}

// GetGuildWidgetSettings returns the widget settings for the guild.
func (c *Client) GetGuildWidgetSettings(ctx context.Context, guildID Snowflake) (*GuildWidget, error) {
	widget := new(GuildWidget)
	route := fmt.Sprintf("/guilds/%s/widget", guildID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, widget)
	if err != nil {
		return nil, err
	}
	return widget, nil
}

// ModifyGuildWidget modifies the guild widget settings.
func (c *Client) ModifyGuildWidget(ctx context.Context, guildID Snowflake, params *ModifyWidgetParams) (*GuildWidget, error) {
	widget := new(GuildWidget)
	route := fmt.Sprintf("/guilds/%s/widget", guildID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, widget)
	if err != nil {
		return nil, err
	}
	return widget, nil
}
