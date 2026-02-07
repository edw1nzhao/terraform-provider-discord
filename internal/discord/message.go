package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateMessageParams are the parameters for creating a message.
type CreateMessageParams struct {
	Content *string  `json:"content,omitempty"`
	TTS     *bool    `json:"tts,omitempty"`
	Embeds  []*Embed `json:"embeds,omitempty"`
	Flags   *int     `json:"flags,omitempty"`
}

// EditMessageParams are the parameters for editing a message.
type EditMessageParams struct {
	Content *string  `json:"content,omitempty"`
	Embeds  []*Embed `json:"embeds,omitempty"`
	Flags   *int     `json:"flags,omitempty"`
}

// CreateMessage posts a message to a channel.
func (c *Client) CreateMessage(ctx context.Context, channelID Snowflake, params *CreateMessageParams) (*Message, error) {
	msg := new(Message)
	route := fmt.Sprintf("/channels/%s/messages", channelID)
	err := c.doRequest(ctx, http.MethodPost, route, params, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// GetChannelMessage returns a specific message in a channel.
func (c *Client) GetChannelMessage(ctx context.Context, channelID Snowflake, messageID Snowflake) (*Message, error) {
	msg := new(Message)
	route := fmt.Sprintf("/channels/%s/messages/%s", channelID, messageID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// EditMessage edits a previously sent message.
func (c *Client) EditMessage(ctx context.Context, channelID Snowflake, messageID Snowflake, params *EditMessageParams) (*Message, error) {
	msg := new(Message)
	route := fmt.Sprintf("/channels/%s/messages/%s", channelID, messageID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// DeleteMessage deletes a message.
func (c *Client) DeleteMessage(ctx context.Context, channelID Snowflake, messageID Snowflake) error {
	route := fmt.Sprintf("/channels/%s/messages/%s", channelID, messageID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}

