package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateChannelParams are the parameters for creating a guild channel.
type CreateChannelParams struct {
	Name                          string                 `json:"name"`
	Type                          *int                   `json:"type,omitempty"`
	Topic                         *string                `json:"topic,omitempty"`
	Bitrate                       *int                   `json:"bitrate,omitempty"`
	UserLimit                     *int                   `json:"user_limit,omitempty"`
	RateLimitPerUser              *int                   `json:"rate_limit_per_user,omitempty"`
	Position                      *int                   `json:"position,omitempty"`
	PermissionOverwrites          []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                      *Snowflake             `json:"parent_id,omitempty"`
	NSFW                          *bool                  `json:"nsfw,omitempty"`
	RTCRegion                     *string                `json:"rtc_region,omitempty"`
	VideoQualityMode              *int                   `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration    *int                   `json:"default_auto_archive_duration,omitempty"`
	DefaultReactionEmoji          *DefaultReaction       `json:"default_reaction_emoji,omitempty"`
	AvailableTags                 []*ForumTag            `json:"available_tags,omitempty"`
	DefaultSortOrder              *int                   `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *int                   `json:"default_forum_layout,omitempty"`
	DefaultThreadRateLimitPerUser *int                   `json:"default_thread_rate_limit_per_user,omitempty"`
}

// ModifyChannelParams are the parameters for modifying a channel.
type ModifyChannelParams struct {
	Name                          *string                `json:"name,omitempty"`
	Type                          *int                   `json:"type,omitempty"`
	Position                      *int                   `json:"position,omitempty"`
	Topic                         *string                `json:"topic,omitempty"`
	NSFW                          *bool                  `json:"nsfw,omitempty"`
	RateLimitPerUser              *int                   `json:"rate_limit_per_user,omitempty"`
	Bitrate                       *int                   `json:"bitrate,omitempty"`
	UserLimit                     *int                   `json:"user_limit,omitempty"`
	PermissionOverwrites          []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                      *Snowflake             `json:"parent_id,omitempty"`
	RTCRegion                     *string                `json:"rtc_region,omitempty"`
	VideoQualityMode              *int                   `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration    *int                   `json:"default_auto_archive_duration,omitempty"`
	Flags                         *int                   `json:"flags,omitempty"`
	AvailableTags                 []*ForumTag            `json:"available_tags,omitempty"`
	DefaultReactionEmoji          *DefaultReaction       `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser *int                   `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              *int                   `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *int                   `json:"default_forum_layout,omitempty"`
}

// EditPermissionsParams are the parameters for editing channel permissions.
type EditPermissionsParams struct {
	Allow *string `json:"allow,omitempty"`
	Deny  *string `json:"deny,omitempty"`
	Type  int     `json:"type"` // 0 = role, 1 = member
}

// CreateGuildChannel creates a new channel in a guild.
func (c *Client) CreateGuildChannel(ctx context.Context, guildID Snowflake, params *CreateChannelParams) (*Channel, error) {
	channel := new(Channel)
	route := fmt.Sprintf("/guilds/%s/channels", guildID)
	err := c.doRequest(ctx, http.MethodPost, route, params, channel)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

// GetChannel returns a channel object for the given ID.
func (c *Client) GetChannel(ctx context.Context, channelID Snowflake) (*Channel, error) {
	channel := new(Channel)
	route := fmt.Sprintf("/channels/%s", channelID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, channel)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

// ModifyChannel updates a channel's settings.
func (c *Client) ModifyChannel(ctx context.Context, channelID Snowflake, params *ModifyChannelParams) (*Channel, error) {
	channel := new(Channel)
	route := fmt.Sprintf("/channels/%s", channelID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, channel)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

// DeleteChannel deletes a channel.
func (c *Client) DeleteChannel(ctx context.Context, channelID Snowflake) error {
	route := fmt.Sprintf("/channels/%s", channelID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}

// EditChannelPermissions edits the channel permission overwrites for a user or role.
func (c *Client) EditChannelPermissions(ctx context.Context, channelID Snowflake, overwriteID Snowflake, params *EditPermissionsParams) error {
	route := fmt.Sprintf("/channels/%s/permissions/%s", channelID, overwriteID)
	return c.doRequestNoContent(ctx, http.MethodPut, route, params)
}

// DeleteChannelPermission deletes a channel permission overwrite for a user or role.
func (c *Client) DeleteChannelPermission(ctx context.Context, channelID Snowflake, overwriteID Snowflake) error {
	route := fmt.Sprintf("/channels/%s/permissions/%s", channelID, overwriteID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}
