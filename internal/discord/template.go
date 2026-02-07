package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateTemplateParams are the parameters for creating a guild template.
type CreateTemplateParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// ModifyTemplateParams are the parameters for modifying a guild template.
type ModifyTemplateParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// GetGuildTemplates returns a list of guild template objects for the guild.
func (c *Client) GetGuildTemplates(ctx context.Context, guildID Snowflake) ([]*GuildTemplate, error) {
	var templates []*GuildTemplate
	route := fmt.Sprintf("/guilds/%s/templates", guildID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, &templates)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// CreateGuildTemplate creates a guild template for the guild.
func (c *Client) CreateGuildTemplate(ctx context.Context, guildID Snowflake, params *CreateTemplateParams) (*GuildTemplate, error) {
	tmpl := new(GuildTemplate)
	route := fmt.Sprintf("/guilds/%s/templates", guildID)
	err := c.doRequest(ctx, http.MethodPost, route, params, tmpl)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// ModifyGuildTemplate modifies a guild template.
func (c *Client) ModifyGuildTemplate(ctx context.Context, guildID Snowflake, code string, params *ModifyTemplateParams) (*GuildTemplate, error) {
	tmpl := new(GuildTemplate)
	route := fmt.Sprintf("/guilds/%s/templates/%s", guildID, code)
	err := c.doRequest(ctx, http.MethodPatch, route, params, tmpl)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// DeleteGuildTemplate deletes a guild template.
func (c *Client) DeleteGuildTemplate(ctx context.Context, guildID Snowflake, code string) error {
	route := fmt.Sprintf("/guilds/%s/templates/%s", guildID, code)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}
