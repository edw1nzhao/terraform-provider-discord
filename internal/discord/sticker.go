package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateStickerParams are the parameters for creating a guild sticker.
// Note: Sticker creation requires multipart/form-data. This struct is used
// for the JSON fields; the actual file upload needs special handling.
type CreateStickerParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Tags        string  `json:"tags"`
}

// ModifyStickerParams are the parameters for modifying a guild sticker.
type ModifyStickerParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Tags        *string `json:"tags,omitempty"`
}

// GetGuildSticker returns a sticker object for the guild and sticker ID.
func (c *Client) GetGuildSticker(ctx context.Context, guildID Snowflake, stickerID Snowflake) (*Sticker, error) {
	sticker := new(Sticker)
	route := fmt.Sprintf("/guilds/%s/stickers/%s", guildID, stickerID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, sticker)
	if err != nil {
		return nil, err
	}
	return sticker, nil
}

// CreateGuildSticker creates a new sticker for the guild.
// Note: This method sends JSON. For actual sticker creation with file upload,
// a multipart form handler would be needed. This provides the JSON-only path.
func (c *Client) CreateGuildSticker(ctx context.Context, guildID Snowflake, params *CreateStickerParams) (*Sticker, error) {
	sticker := new(Sticker)
	route := fmt.Sprintf("/guilds/%s/stickers", guildID)
	err := c.doRequest(ctx, http.MethodPost, route, params, sticker)
	if err != nil {
		return nil, err
	}
	return sticker, nil
}

// ModifyGuildSticker modifies the given sticker.
func (c *Client) ModifyGuildSticker(ctx context.Context, guildID Snowflake, stickerID Snowflake, params *ModifyStickerParams) (*Sticker, error) {
	sticker := new(Sticker)
	route := fmt.Sprintf("/guilds/%s/stickers/%s", guildID, stickerID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, sticker)
	if err != nil {
		return nil, err
	}
	return sticker, nil
}

// DeleteGuildSticker deletes the given sticker.
func (c *Client) DeleteGuildSticker(ctx context.Context, guildID Snowflake, stickerID Snowflake) error {
	route := fmt.Sprintf("/guilds/%s/stickers/%s", guildID, stickerID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}
