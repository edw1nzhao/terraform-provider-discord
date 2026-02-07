package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateCommandParams are the parameters for creating an application command.
type CreateCommandParams struct {
	Name                     string                      `json:"name"`
	NameLocalizations        map[string]string           `json:"name_localizations,omitempty"`
	Description              *string                     `json:"description,omitempty"`
	DescriptionLocalizations map[string]string           `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                     `json:"default_member_permissions,omitempty"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"`
	Type                     *int                        `json:"type,omitempty"`
	NSFW                     *bool                       `json:"nsfw,omitempty"`
}

// EditCommandParams are the parameters for editing an application command.
type EditCommandParams struct {
	Name                     *string                     `json:"name,omitempty"`
	NameLocalizations        map[string]string           `json:"name_localizations,omitempty"`
	Description              *string                     `json:"description,omitempty"`
	DescriptionLocalizations map[string]string           `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                     `json:"default_member_permissions,omitempty"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"`
	NSFW                     *bool                       `json:"nsfw,omitempty"`
}

// --- Global Application Commands ---

// CreateGlobalApplicationCommand creates a new global application command.
func (c *Client) CreateGlobalApplicationCommand(ctx context.Context, appID Snowflake, params *CreateCommandParams) (*ApplicationCommand, error) {
	cmd := new(ApplicationCommand)
	route := fmt.Sprintf("/applications/%s/commands", appID)
	err := c.doRequest(ctx, http.MethodPost, route, params, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// GetGlobalApplicationCommand returns a specific global application command.
func (c *Client) GetGlobalApplicationCommand(ctx context.Context, appID Snowflake, cmdID Snowflake) (*ApplicationCommand, error) {
	cmd := new(ApplicationCommand)
	route := fmt.Sprintf("/applications/%s/commands/%s", appID, cmdID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// EditGlobalApplicationCommand edits a global application command.
func (c *Client) EditGlobalApplicationCommand(ctx context.Context, appID Snowflake, cmdID Snowflake, params *EditCommandParams) (*ApplicationCommand, error) {
	cmd := new(ApplicationCommand)
	route := fmt.Sprintf("/applications/%s/commands/%s", appID, cmdID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// DeleteGlobalApplicationCommand deletes a global application command.
func (c *Client) DeleteGlobalApplicationCommand(ctx context.Context, appID Snowflake, cmdID Snowflake) error {
	route := fmt.Sprintf("/applications/%s/commands/%s", appID, cmdID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}

// --- Guild Application Commands ---

// CreateGuildApplicationCommand creates a new guild application command.
func (c *Client) CreateGuildApplicationCommand(ctx context.Context, appID Snowflake, guildID Snowflake, params *CreateCommandParams) (*ApplicationCommand, error) {
	cmd := new(ApplicationCommand)
	route := fmt.Sprintf("/applications/%s/guilds/%s/commands", appID, guildID)
	err := c.doRequest(ctx, http.MethodPost, route, params, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// GetGuildApplicationCommand returns a specific guild application command.
func (c *Client) GetGuildApplicationCommand(ctx context.Context, appID Snowflake, guildID Snowflake, cmdID Snowflake) (*ApplicationCommand, error) {
	cmd := new(ApplicationCommand)
	route := fmt.Sprintf("/applications/%s/guilds/%s/commands/%s", appID, guildID, cmdID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// EditGuildApplicationCommand edits a guild application command.
func (c *Client) EditGuildApplicationCommand(ctx context.Context, appID Snowflake, guildID Snowflake, cmdID Snowflake, params *EditCommandParams) (*ApplicationCommand, error) {
	cmd := new(ApplicationCommand)
	route := fmt.Sprintf("/applications/%s/guilds/%s/commands/%s", appID, guildID, cmdID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// DeleteGuildApplicationCommand deletes a guild application command.
func (c *Client) DeleteGuildApplicationCommand(ctx context.Context, appID Snowflake, guildID Snowflake, cmdID Snowflake) error {
	route := fmt.Sprintf("/applications/%s/guilds/%s/commands/%s", appID, guildID, cmdID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}
