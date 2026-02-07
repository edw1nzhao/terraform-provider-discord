package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateWebhookParams are the parameters for creating a webhook.
type CreateWebhookParams struct {
	Name   string  `json:"name"`
	Avatar *string `json:"avatar,omitempty"`
}

// ModifyWebhookParams are the parameters for modifying a webhook.
type ModifyWebhookParams struct {
	Name      *string    `json:"name,omitempty"`
	Avatar    *string    `json:"avatar,omitempty"`
	ChannelID *Snowflake `json:"channel_id,omitempty"`
}

// CreateWebhook creates a new webhook for a channel.
func (c *Client) CreateWebhook(ctx context.Context, channelID Snowflake, params *CreateWebhookParams) (*Webhook, error) {
	webhook := new(Webhook)
	route := fmt.Sprintf("/channels/%s/webhooks", channelID)
	err := c.doRequest(ctx, http.MethodPost, route, params, webhook)
	if err != nil {
		return nil, err
	}
	return webhook, nil
}

// GetWebhook returns the webhook object for the given ID.
func (c *Client) GetWebhook(ctx context.Context, webhookID Snowflake) (*Webhook, error) {
	webhook := new(Webhook)
	route := fmt.Sprintf("/webhooks/%s", webhookID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, webhook)
	if err != nil {
		return nil, err
	}
	return webhook, nil
}

// ModifyWebhook modifies a webhook.
func (c *Client) ModifyWebhook(ctx context.Context, webhookID Snowflake, params *ModifyWebhookParams) (*Webhook, error) {
	webhook := new(Webhook)
	route := fmt.Sprintf("/webhooks/%s", webhookID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, webhook)
	if err != nil {
		return nil, err
	}
	return webhook, nil
}

// DeleteWebhook deletes a webhook.
func (c *Client) DeleteWebhook(ctx context.Context, webhookID Snowflake) error {
	route := fmt.Sprintf("/webhooks/%s", webhookID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}

