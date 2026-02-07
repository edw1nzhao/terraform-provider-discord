package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateGuildParams are the parameters for creating a guild.
type CreateGuildParams struct {
	Name                        string                 `json:"name"`
	Region                      *string                `json:"region,omitempty"`
	Icon                        *string                `json:"icon,omitempty"`
	VerificationLevel           *int                   `json:"verification_level,omitempty"`
	DefaultMessageNotifications *int                   `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *int                   `json:"explicit_content_filter,omitempty"`
	Roles                       []*Role                `json:"roles,omitempty"`
	Channels                    []*Channel             `json:"channels,omitempty"`
	AFKChannelID                *Snowflake             `json:"afk_channel_id,omitempty"`
	AFKTimeout                  *int                   `json:"afk_timeout,omitempty"`
	SystemChannelID             *Snowflake             `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *int                   `json:"system_channel_flags,omitempty"`
}

// ModifyGuildParams are the parameters for modifying a guild.
type ModifyGuildParams struct {
	Name                        *string    `json:"name,omitempty"`
	Region                      *string    `json:"region,omitempty"`
	VerificationLevel           *int       `json:"verification_level,omitempty"`
	DefaultMessageNotifications *int       `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *int       `json:"explicit_content_filter,omitempty"`
	AFKChannelID                *Snowflake `json:"afk_channel_id,omitempty"`
	AFKTimeout                  *int       `json:"afk_timeout,omitempty"`
	Icon                        *string    `json:"icon,omitempty"`
	OwnerID                     *Snowflake `json:"owner_id,omitempty"`
	Splash                      *string    `json:"splash,omitempty"`
	DiscoverySplash             *string    `json:"discovery_splash,omitempty"`
	Banner                      *string    `json:"banner,omitempty"`
	SystemChannelID             *Snowflake `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *int       `json:"system_channel_flags,omitempty"`
	RulesChannelID              *Snowflake `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      *Snowflake `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             *string    `json:"preferred_locale,omitempty"`
	Features                    []string   `json:"features,omitempty"`
	Description                 *string    `json:"description,omitempty"`
	PremiumProgressBarEnabled   *bool      `json:"premium_progress_bar_enabled,omitempty"`
	SafetyAlertsChannelID       *Snowflake `json:"safety_alerts_channel_id,omitempty"`
}

// CreateGuild creates a new guild. Returns the created guild object.
func (c *Client) CreateGuild(ctx context.Context, params *CreateGuildParams) (*Guild, error) {
	guild := new(Guild)
	err := c.doRequest(ctx, http.MethodPost, "/guilds", params, guild)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

// GetGuild returns the guild object for the given ID.
func (c *Client) GetGuild(ctx context.Context, guildID Snowflake) (*Guild, error) {
	guild := new(Guild)
	route := fmt.Sprintf("/guilds/%s", guildID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, guild)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

// ModifyGuild modifies a guild's settings. Returns the updated guild object.
func (c *Client) ModifyGuild(ctx context.Context, guildID Snowflake, params *ModifyGuildParams) (*Guild, error) {
	guild := new(Guild)
	route := fmt.Sprintf("/guilds/%s", guildID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, guild)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

// DeleteGuild deletes a guild permanently. The bot must be the owner.
func (c *Client) DeleteGuild(ctx context.Context, guildID Snowflake) error {
	route := fmt.Sprintf("/guilds/%s", guildID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}

