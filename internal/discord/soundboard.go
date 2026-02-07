package discord

import (
	"context"
	"fmt"
	"net/http"
)

// soundboardSoundsResponse wraps the list response for soundboard sounds.
type soundboardSoundsResponse struct {
	Items []*SoundboardSound `json:"items"`
}

// CreateSoundboardSoundParams are the parameters for creating a soundboard sound.
type CreateSoundboardSoundParams struct {
	Name      string     `json:"name"`
	Sound     string     `json:"sound"` // Data URI (data:audio/ogg;base64,...)
	Volume    *float64   `json:"volume,omitempty"`
	EmojiID   *Snowflake `json:"emoji_id,omitempty"`
	EmojiName *string    `json:"emoji_name,omitempty"`
}

// ModifySoundboardSoundParams are the parameters for modifying a soundboard sound.
type ModifySoundboardSoundParams struct {
	Name      *string    `json:"name,omitempty"`
	Volume    *float64   `json:"volume,omitempty"`
	EmojiID   *Snowflake `json:"emoji_id,omitempty"`
	EmojiName *string    `json:"emoji_name,omitempty"`
}

// ListGuildSoundboardSounds returns a list of soundboard sounds for the guild.
func (c *Client) ListGuildSoundboardSounds(ctx context.Context, guildID Snowflake) ([]*SoundboardSound, error) {
	var resp soundboardSoundsResponse
	route := fmt.Sprintf("/guilds/%s/soundboard-sounds", guildID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// CreateGuildSoundboardSound creates a new soundboard sound for the guild.
func (c *Client) CreateGuildSoundboardSound(ctx context.Context, guildID Snowflake, params *CreateSoundboardSoundParams) (*SoundboardSound, error) {
	sound := new(SoundboardSound)
	route := fmt.Sprintf("/guilds/%s/soundboard-sounds", guildID)
	err := c.doRequest(ctx, http.MethodPost, route, params, sound)
	if err != nil {
		return nil, err
	}
	return sound, nil
}

// ModifyGuildSoundboardSound modifies a soundboard sound.
func (c *Client) ModifyGuildSoundboardSound(ctx context.Context, guildID Snowflake, soundID Snowflake, params *ModifySoundboardSoundParams) (*SoundboardSound, error) {
	sound := new(SoundboardSound)
	route := fmt.Sprintf("/guilds/%s/soundboard-sounds/%s", guildID, soundID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, sound)
	if err != nil {
		return nil, err
	}
	return sound, nil
}

// DeleteGuildSoundboardSound deletes a soundboard sound.
func (c *Client) DeleteGuildSoundboardSound(ctx context.Context, guildID Snowflake, soundID Snowflake) error {
	route := fmt.Sprintf("/guilds/%s/soundboard-sounds/%s", guildID, soundID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}

