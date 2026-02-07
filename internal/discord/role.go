package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateRoleParams are the parameters for creating a guild role.
type CreateRoleParams struct {
	Name         *string `json:"name,omitempty"`
	Permissions  *string `json:"permissions,omitempty"`
	Color        *int    `json:"color,omitempty"`
	Hoist        *bool   `json:"hoist,omitempty"`
	Icon         *string `json:"icon,omitempty"`
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	Mentionable  *bool   `json:"mentionable,omitempty"`
}

// ModifyRoleParams are the parameters for modifying a guild role.
type ModifyRoleParams struct {
	Name         *string `json:"name,omitempty"`
	Permissions  *string `json:"permissions,omitempty"`
	Color        *int    `json:"color,omitempty"`
	Hoist        *bool   `json:"hoist,omitempty"`
	Icon         *string `json:"icon,omitempty"`
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	Mentionable  *bool   `json:"mentionable,omitempty"`
}

// RolePosition represents a role position update for ModifyGuildRolePositions.
type RolePosition struct {
	ID       Snowflake `json:"id"`
	Position *int      `json:"position,omitempty"`
}

// CreateGuildRole creates a new role for the guild.
func (c *Client) CreateGuildRole(ctx context.Context, guildID Snowflake, params *CreateRoleParams) (*Role, error) {
	role := new(Role)
	route := fmt.Sprintf("/guilds/%s/roles", guildID)
	err := c.doRequest(ctx, http.MethodPost, route, params, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// GetGuildRoles returns a list of roles for the guild.
func (c *Client) GetGuildRoles(ctx context.Context, guildID Snowflake) ([]*Role, error) {
	var roles []*Role
	route := fmt.Sprintf("/guilds/%s/roles", guildID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// ModifyGuildRole modifies a guild role.
func (c *Client) ModifyGuildRole(ctx context.Context, guildID Snowflake, roleID Snowflake, params *ModifyRoleParams) (*Role, error) {
	role := new(Role)
	route := fmt.Sprintf("/guilds/%s/roles/%s", guildID, roleID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// DeleteGuildRole deletes a guild role.
func (c *Client) DeleteGuildRole(ctx context.Context, guildID Snowflake, roleID Snowflake) error {
	route := fmt.Sprintf("/guilds/%s/roles/%s", guildID, roleID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}

// ModifyGuildRolePositions modifies the positions of a set of role objects for the guild.
func (c *Client) ModifyGuildRolePositions(ctx context.Context, guildID Snowflake, positions []*RolePosition) ([]*Role, error) {
	var roles []*Role
	route := fmt.Sprintf("/guilds/%s/roles", guildID)
	err := c.doRequest(ctx, http.MethodPatch, route, positions, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}
