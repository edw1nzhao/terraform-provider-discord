package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateStageInstanceParams are the parameters for creating a stage instance.
type CreateStageInstanceParams struct {
	ChannelID             Snowflake  `json:"channel_id"`
	Topic                 string     `json:"topic"`
	PrivacyLevel          *int       `json:"privacy_level,omitempty"`
	SendStartNotification *bool      `json:"send_start_notification,omitempty"`
	GuildScheduledEventID *Snowflake `json:"guild_scheduled_event_id,omitempty"`
}

// ModifyStageInstanceParams are the parameters for modifying a stage instance.
type ModifyStageInstanceParams struct {
	Topic        *string `json:"topic,omitempty"`
	PrivacyLevel *int    `json:"privacy_level,omitempty"`
}

// CreateStageInstance creates a new stage instance associated with a stage channel.
func (c *Client) CreateStageInstance(ctx context.Context, params *CreateStageInstanceParams) (*StageInstance, error) {
	stage := new(StageInstance)
	err := c.doRequest(ctx, http.MethodPost, "/stage-instances", params, stage)
	if err != nil {
		return nil, err
	}
	return stage, nil
}

// GetStageInstance returns the stage instance for the given channel ID.
func (c *Client) GetStageInstance(ctx context.Context, channelID Snowflake) (*StageInstance, error) {
	stage := new(StageInstance)
	route := fmt.Sprintf("/stage-instances/%s", channelID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, stage)
	if err != nil {
		return nil, err
	}
	return stage, nil
}

// ModifyStageInstance modifies an existing stage instance.
func (c *Client) ModifyStageInstance(ctx context.Context, channelID Snowflake, params *ModifyStageInstanceParams) (*StageInstance, error) {
	stage := new(StageInstance)
	route := fmt.Sprintf("/stage-instances/%s", channelID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, stage)
	if err != nil {
		return nil, err
	}
	return stage, nil
}

// DeleteStageInstance deletes the stage instance for the given channel ID.
func (c *Client) DeleteStageInstance(ctx context.Context, channelID Snowflake) error {
	route := fmt.Sprintf("/stage-instances/%s", channelID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}
