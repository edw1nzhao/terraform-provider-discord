package discord

import (
	"context"
	"fmt"
	"net/http"
)

// ModifyWelcomeScreenParams are the parameters for modifying a guild's welcome screen.
type ModifyWelcomeScreenParams struct {
	Enabled         *bool                   `json:"enabled,omitempty"`
	WelcomeChannels []*WelcomeScreenChannel `json:"welcome_channels,omitempty"`
	Description     *string                 `json:"description,omitempty"`
}

// GetGuildWelcomeScreen returns the welcome screen for the guild.
func (c *Client) GetGuildWelcomeScreen(ctx context.Context, guildID Snowflake) (*WelcomeScreen, error) {
	ws := new(WelcomeScreen)
	route := fmt.Sprintf("/guilds/%s/welcome-screen", guildID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, ws)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// ModifyGuildWelcomeScreen modifies the guild's welcome screen.
func (c *Client) ModifyGuildWelcomeScreen(ctx context.Context, guildID Snowflake, params *ModifyWelcomeScreenParams) (*WelcomeScreen, error) {
	ws := new(WelcomeScreen)
	route := fmt.Sprintf("/guilds/%s/welcome-screen", guildID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, ws)
	if err != nil {
		return nil, err
	}
	return ws, nil
}
