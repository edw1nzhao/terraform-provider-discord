package discord

import (
	"context"
	"fmt"
	"net/http"
)

// CreateAutoModRuleParams are the parameters for creating an auto-moderation rule.
type CreateAutoModRuleParams struct {
	Name            string           `json:"name"`
	EventType       int              `json:"event_type"`
	TriggerType     int              `json:"trigger_type"`
	TriggerMetadata *TriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         []*AutoModAction `json:"actions"`
	Enabled         *bool            `json:"enabled,omitempty"`
	ExemptRoles     []Snowflake      `json:"exempt_roles,omitempty"`
	ExemptChannels  []Snowflake      `json:"exempt_channels,omitempty"`
}

// ModifyAutoModRuleParams are the parameters for modifying an auto-moderation rule.
type ModifyAutoModRuleParams struct {
	Name            *string          `json:"name,omitempty"`
	EventType       *int             `json:"event_type,omitempty"`
	TriggerMetadata *TriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         []*AutoModAction `json:"actions,omitempty"`
	Enabled         *bool            `json:"enabled,omitempty"`
	ExemptRoles     []Snowflake      `json:"exempt_roles,omitempty"`
	ExemptChannels  []Snowflake      `json:"exempt_channels,omitempty"`
}

// GetAutoModerationRule returns a single auto-moderation rule.
func (c *Client) GetAutoModerationRule(ctx context.Context, guildID Snowflake, ruleID Snowflake) (*AutoModerationRule, error) {
	rule := new(AutoModerationRule)
	route := fmt.Sprintf("/guilds/%s/auto-moderation/rules/%s", guildID, ruleID)
	err := c.doRequest(ctx, http.MethodGet, route, nil, rule)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

// CreateAutoModerationRule creates a new auto-moderation rule.
func (c *Client) CreateAutoModerationRule(ctx context.Context, guildID Snowflake, params *CreateAutoModRuleParams) (*AutoModerationRule, error) {
	rule := new(AutoModerationRule)
	route := fmt.Sprintf("/guilds/%s/auto-moderation/rules", guildID)
	err := c.doRequest(ctx, http.MethodPost, route, params, rule)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

// ModifyAutoModerationRule modifies an existing auto-moderation rule.
func (c *Client) ModifyAutoModerationRule(ctx context.Context, guildID Snowflake, ruleID Snowflake, params *ModifyAutoModRuleParams) (*AutoModerationRule, error) {
	rule := new(AutoModerationRule)
	route := fmt.Sprintf("/guilds/%s/auto-moderation/rules/%s", guildID, ruleID)
	err := c.doRequest(ctx, http.MethodPatch, route, params, rule)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

// DeleteAutoModerationRule deletes an auto-moderation rule.
func (c *Client) DeleteAutoModerationRule(ctx context.Context, guildID Snowflake, ruleID Snowflake) error {
	route := fmt.Sprintf("/guilds/%s/auto-moderation/rules/%s", guildID, ruleID)
	return c.doRequestNoContent(ctx, http.MethodDelete, route, nil)
}
