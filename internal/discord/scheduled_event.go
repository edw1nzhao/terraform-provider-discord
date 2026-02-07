package discord

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// CreateScheduledEventParams are the parameters for creating a guild scheduled event.
type CreateScheduledEventParams struct {
	ChannelID          *Snowflake                    `json:"channel_id,omitempty"`
	EntityMetadata     *ScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               string                        `json:"name"`
	PrivacyLevel       int                           `json:"privacy_level"`
	ScheduledStartTime time.Time                     `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time                    `json:"scheduled_end_time,omitempty"`
	Description        *string                       `json:"description,omitempty"`
	EntityType         int                           `json:"entity_type"`
	Image              *string                       `json:"image,omitempty"`
}

// ModifyScheduledEventParams are the parameters for modifying a guild scheduled event.
type ModifyScheduledEventParams struct {
	ChannelID          *Snowflake                    `json:"channel_id,omitempty"`
	EntityMetadata     *ScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               *string                       `json:"name,omitempty"`
	PrivacyLevel       *int                          `json:"privacy_level,omitempty"`
	ScheduledStartTime *time.Time                    `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   *time.Time                    `json:"scheduled_end_time,omitempty"`
	Description        *string                       `json:"description,omitempty"`
	EntityType         *int                          `json:"entity_type,omitempty"`
	Status             *int                          `json:"status,omitempty"`
	Image              *string                       `json:"image,omitempty"`
}

// GetGuildScheduledEvent returns a single scheduled event.
func (c *Client) GetGuildScheduledEvent(ctx context.Context, guildID Snowflake, eventID Snowflake) (*GuildScheduledEvent, error) {
	event := new(GuildScheduledEvent)
	route := fmt.Sprintf("/guilds/%s/scheduled-events/%s", guildID, eventID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// CreateGuildScheduledEvent creates a new scheduled event for the guild.
func (c *Client) CreateGuildScheduledEvent(ctx context.Context, guildID Snowflake, params *CreateScheduledEventParams) (*GuildScheduledEvent, error) {
	event := new(GuildScheduledEvent)
	route := fmt.Sprintf("/guilds/%s/scheduled-events", guildID)
	err := c.doRequest(ctx, http.MethodPost, route, params, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// ModifyGuildScheduledEvent modifies a guild scheduled event.
func (c *Client) ModifyGuildScheduledEvent(ctx context.Context, guildID Snowflake, eventID Snowflake, params *ModifyScheduledEventParams) (*GuildScheduledEvent, error) {
	event := new(GuildScheduledEvent)
	route := fmt.Sprintf("/guilds/%s/scheduled-events/%s", guildID, eventID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// DeleteGuildScheduledEvent deletes a guild scheduled event.
func (c *Client) DeleteGuildScheduledEvent(ctx context.Context, guildID Snowflake, eventID Snowflake) error {
	route := fmt.Sprintf("/guilds/%s/scheduled-events/%s", guildID, eventID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}
