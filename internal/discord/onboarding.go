package discord

import (
	"context"
	"fmt"
	"net/http"
)

// ModifyOnboardingParams are the parameters for modifying guild onboarding.
type ModifyOnboardingParams struct {
	Prompts           []*OnboardingPrompt `json:"prompts,omitempty"`
	DefaultChannelIDs []Snowflake         `json:"default_channel_ids,omitempty"`
	Enabled           *bool               `json:"enabled,omitempty"`
	Mode              *int                `json:"mode,omitempty"`
}

// GetGuildOnboarding returns the onboarding configuration for the guild.
func (c *Client) GetGuildOnboarding(ctx context.Context, guildID Snowflake) (*GuildOnboarding, error) {
	ob := new(GuildOnboarding)
	route := fmt.Sprintf("/guilds/%s/onboarding", guildID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, ob)
	if err != nil {
		return nil, err
	}
	return ob, nil
}

// ModifyGuildOnboarding modifies the guild's onboarding configuration.
func (c *Client) ModifyGuildOnboarding(ctx context.Context, guildID Snowflake, params *ModifyOnboardingParams) (*GuildOnboarding, error) {
	ob := new(GuildOnboarding)
	route := fmt.Sprintf("/guilds/%s/onboarding", guildID)
	err := c.doRequest(ctx, http.MethodPut, route, params, ob)
	if err != nil {
		return nil, err
	}
	return ob, nil
}
