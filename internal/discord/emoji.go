package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateEmojiParams are the parameters for creating a guild emoji.
type CreateEmojiParams struct {
	Name  string      `json:"name"`
	Image string      `json:"image"` // Data URI scheme (data:image/png;base64,...)
	Roles []Snowflake `json:"roles,omitempty"`
}

// ModifyEmojiParams are the parameters for modifying a guild emoji.
type ModifyEmojiParams struct {
	Name  *string     `json:"name,omitempty"`
	Roles []Snowflake `json:"roles,omitempty"`
}

// GetGuildEmoji returns an emoji object for the guild and emoji ID.
func (c *Client) GetGuildEmoji(ctx context.Context, guildID Snowflake, emojiID Snowflake) (*Emoji, error) {
	emoji := new(Emoji)
	route := fmt.Sprintf("/guilds/%s/emojis/%s", guildID, emojiID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, emoji)
	if err != nil {
		return nil, err
	}
	return emoji, nil
}

// CreateGuildEmoji creates a new emoji for the guild.
func (c *Client) CreateGuildEmoji(ctx context.Context, guildID Snowflake, params *CreateEmojiParams) (*Emoji, error) {
	emoji := new(Emoji)
	route := fmt.Sprintf("/guilds/%s/emojis", guildID)
	err := c.doRequest(ctx, http.MethodPost, route, params, emoji)
	if err != nil {
		return nil, err
	}
	return emoji, nil
}

// ModifyGuildEmoji modifies the given emoji.
func (c *Client) ModifyGuildEmoji(ctx context.Context, guildID Snowflake, emojiID Snowflake, params *ModifyEmojiParams) (*Emoji, error) {
	emoji := new(Emoji)
	route := fmt.Sprintf("/guilds/%s/emojis/%s", guildID, emojiID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, emoji)
	if err != nil {
		return nil, err
	}
	return emoji, nil
}

// DeleteGuildEmoji deletes the given emoji.
func (c *Client) DeleteGuildEmoji(ctx context.Context, guildID Snowflake, emojiID Snowflake) error {
	route := fmt.Sprintf("/guilds/%s/emojis/%s", guildID, emojiID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}
