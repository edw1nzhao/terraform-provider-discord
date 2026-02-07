package automod

import (
	"context"
	"strings"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure interface compliance.
var (
	_ resource.Resource                = &autoModerationRuleResource{}
	_ resource.ResourceWithConfigure   = &autoModerationRuleResource{}
	_ resource.ResourceWithImportState = &autoModerationRuleResource{}
)

// NewAutoModerationRuleResource returns a new resource for discord_auto_moderation_rule.
func NewAutoModerationRuleResource() resource.Resource {
	return &autoModerationRuleResource{}
}

type autoModerationRuleResource struct {
	client *discord.Client
}

// autoModerationRuleModel maps the Terraform schema to Go types.
type autoModerationRuleModel struct {
	ID              types.String `tfsdk:"id"`
	GuildID         types.String `tfsdk:"guild_id"`
	Name            types.String `tfsdk:"name"`
	EventType       types.Int64  `tfsdk:"event_type"`
	TriggerType     types.Int64  `tfsdk:"trigger_type"`
	TriggerMetadata types.Object `tfsdk:"trigger_metadata"`
	Actions         types.List   `tfsdk:"actions"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	ExemptRoles     types.Set    `tfsdk:"exempt_roles"`
	ExemptChannels  types.Set    `tfsdk:"exempt_channels"`
}

// triggerMetadataModel maps the trigger_metadata nested block.
type triggerMetadataModel struct {
	KeywordFilter                types.List  `tfsdk:"keyword_filter"`
	RegexPatterns                types.List  `tfsdk:"regex_patterns"`
	Presets                      types.List  `tfsdk:"presets"`
	AllowList                    types.List  `tfsdk:"allow_list"`
	MentionTotalLimit            types.Int64 `tfsdk:"mention_total_limit"`
	MentionRaidProtectionEnabled types.Bool  `tfsdk:"mention_raid_protection_enabled"`
}

// actionModel maps a single action element.
type actionModel struct {
	Type     types.Int64  `tfsdk:"type"`
	Metadata types.Object `tfsdk:"metadata"`
}

// actionMetadataModel maps the metadata nested block inside an action.
type actionMetadataModel struct {
	ChannelID       types.String `tfsdk:"channel_id"`
	DurationSeconds types.Int64  `tfsdk:"duration_seconds"`
	CustomMessage   types.String `tfsdk:"custom_message"`
}

// Attribute types for the nested objects -----------------------------------------------

func triggerMetadataAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"keyword_filter":                  types.ListType{ElemType: types.StringType},
		"regex_patterns":                  types.ListType{ElemType: types.StringType},
		"presets":                         types.ListType{ElemType: types.Int64Type},
		"allow_list":                      types.ListType{ElemType: types.StringType},
		"mention_total_limit":             types.Int64Type,
		"mention_raid_protection_enabled": types.BoolType,
	}
}

func actionMetadataAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"channel_id":       types.StringType,
		"duration_seconds": types.Int64Type,
		"custom_message":   types.StringType,
	}
}

func actionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":     types.Int64Type,
		"metadata": types.ObjectType{AttrTypes: actionMetadataAttrTypes()},
	}
}

// Metadata sets the type name.
func (r *autoModerationRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auto_moderation_rule"
}

// Schema defines the schema.
func (r *autoModerationRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord auto-moderation rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the auto-moderation rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the auto-moderation rule.",
				Required:    true,
			},
			"event_type": schema.Int64Attribute{
				Description: "The event type that triggers the rule (1 = MESSAGE_SEND, 2 = MEMBER_UPDATE).",
				Required:    true,
			},
			"trigger_type": schema.Int64Attribute{
				Description: "The trigger type (1 = KEYWORD, 3 = SPAM, 4 = KEYWORD_PRESET, 5 = MENTION_SPAM, 6 = MEMBER_PROFILE).",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"trigger_metadata": schema.SingleNestedAttribute{
				Description: "Additional metadata for the trigger.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"keyword_filter": schema.ListAttribute{
						Description: "Substrings which will be searched for in content.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"regex_patterns": schema.ListAttribute{
						Description: "Regular expression patterns which will be matched against content.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"presets": schema.ListAttribute{
						Description: "Preset keyword lists (1 = Profanity, 2 = Sexual Content, 3 = Slurs).",
						Optional:    true,
						ElementType: types.Int64Type,
					},
					"allow_list": schema.ListAttribute{
						Description: "Substrings which should not trigger the rule.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"mention_total_limit": schema.Int64Attribute{
						Description: "Total number of unique role and user mentions allowed per message.",
						Optional:    true,
					},
					"mention_raid_protection_enabled": schema.BoolAttribute{
						Description: "Whether to automatically detect mention raids.",
						Optional:    true,
					},
				},
			},
			"actions": schema.ListNestedAttribute{
				Description: "The actions to take when the rule is triggered.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.Int64Attribute{
							Description: "The type of action (1 = BLOCK_MESSAGE, 2 = SEND_ALERT_MESSAGE, 3 = TIMEOUT, 4 = BLOCK_MEMBER_INTERACTION).",
							Required:    true,
						},
						"metadata": schema.SingleNestedAttribute{
							Description: "Additional metadata for the action.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"channel_id": schema.StringAttribute{
									Description: "Channel to which user content should be logged.",
									Optional:    true,
								},
								"duration_seconds": schema.Int64Attribute{
									Description: "Timeout duration in seconds.",
									Optional:    true,
								},
								"custom_message": schema.StringAttribute{
									Description: "Additional explanation shown to members when their message is blocked.",
									Optional:    true,
								},
							},
						},
					},
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the rule is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"exempt_roles": schema.SetAttribute{
				Description: "Role IDs that are exempt from the rule (max 20).",
				Optional:    true,
				ElementType: types.StringType,
			},
			"exempt_channels": schema.SetAttribute{
				Description: "Channel IDs that are exempt from the rule (max 50).",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure stores the provider data.
func (r *autoModerationRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create creates a new auto-moderation rule.
func (r *autoModerationRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan autoModerationRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.CreateAutoModRuleParams{
		Name:        plan.Name.ValueString(),
		EventType:   int(plan.EventType.ValueInt64()),
		TriggerType: int(plan.TriggerType.ValueInt64()),
	}

	enabled := plan.Enabled.ValueBool()
	params.Enabled = &enabled

	// Trigger metadata.
	if !plan.TriggerMetadata.IsNull() && !plan.TriggerMetadata.IsUnknown() {
		tm, diags := expandTriggerMetadata(ctx, plan.TriggerMetadata)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		params.TriggerMetadata = tm
	}

	// Actions.
	actions, diags := expandActions(ctx, plan.Actions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	params.Actions = actions

	// Exempt roles.
	params.ExemptRoles = expandSnowflakeSet(ctx, plan.ExemptRoles, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Exempt channels.
	params.ExemptChannels = expandSnowflakeSet(ctx, plan.ExemptChannels, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.CreateAutoModerationRule(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating auto-moderation rule", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenAutoModRule(ctx, rule, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state from the API.
func (r *autoModerationRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state autoModerationRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetAutoModerationRule(
		ctx,
		discord.Snowflake(state.GuildID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
	)
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading auto-moderation rule", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenAutoModRule(ctx, rule, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update modifies an existing auto-moderation rule.
func (r *autoModerationRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan autoModerationRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	eventType := int(plan.EventType.ValueInt64())
	enabled := plan.Enabled.ValueBool()

	params := &discord.ModifyAutoModRuleParams{
		Name:      &name,
		EventType: &eventType,
		Enabled:   &enabled,
	}

	// Trigger metadata.
	if !plan.TriggerMetadata.IsNull() && !plan.TriggerMetadata.IsUnknown() {
		tm, diags := expandTriggerMetadata(ctx, plan.TriggerMetadata)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		params.TriggerMetadata = tm
	}

	// Actions.
	actions, diags := expandActions(ctx, plan.Actions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	params.Actions = actions

	// Exempt roles.
	params.ExemptRoles = expandSnowflakeSet(ctx, plan.ExemptRoles, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Exempt channels.
	params.ExemptChannels = expandSnowflakeSet(ctx, plan.ExemptChannels, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.ModifyAutoModerationRule(
		ctx,
		discord.Snowflake(plan.GuildID.ValueString()),
		discord.Snowflake(plan.ID.ValueString()),
		params,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating auto-moderation rule", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenAutoModRule(ctx, rule, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the auto-moderation rule.
func (r *autoModerationRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state autoModerationRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAutoModerationRule(
		ctx,
		discord.Snowflake(state.GuildID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
	)
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting auto-moderation rule", err.Error())
	}
}

// ImportState supports importing by guild_id/rule_id.
func (r *autoModerationRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: guild_id/rule_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// --- Helper functions ---

func expandTriggerMetadata(ctx context.Context, obj types.Object) (*discord.TriggerMetadata, diag.Diagnostics) {
	var diags diag.Diagnostics
	var tm triggerMetadataModel
	diags.Append(obj.As(ctx, &tm, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	result := &discord.TriggerMetadata{}

	if !tm.KeywordFilter.IsNull() && !tm.KeywordFilter.IsUnknown() {
		var kf []string
		diags.Append(tm.KeywordFilter.ElementsAs(ctx, &kf, false)...)
		result.KeywordFilter = kf
	}

	if !tm.RegexPatterns.IsNull() && !tm.RegexPatterns.IsUnknown() {
		var rp []string
		diags.Append(tm.RegexPatterns.ElementsAs(ctx, &rp, false)...)
		result.RegexPatterns = rp
	}

	if !tm.Presets.IsNull() && !tm.Presets.IsUnknown() {
		var presets []int
		// Framework gives us int64, we need to convert
		var presets64 []int64
		diags.Append(tm.Presets.ElementsAs(ctx, &presets64, false)...)
		for _, p := range presets64 {
			presets = append(presets, int(p))
		}
		result.Presets = presets
	}

	if !tm.AllowList.IsNull() && !tm.AllowList.IsUnknown() {
		var al []string
		diags.Append(tm.AllowList.ElementsAs(ctx, &al, false)...)
		result.AllowList = al
	}

	if !tm.MentionTotalLimit.IsNull() && !tm.MentionTotalLimit.IsUnknown() {
		v := int(tm.MentionTotalLimit.ValueInt64())
		result.MentionTotalLimit = &v
	}

	if !tm.MentionRaidProtectionEnabled.IsNull() && !tm.MentionRaidProtectionEnabled.IsUnknown() {
		result.MentionRaidProtectionEnabled = tm.MentionRaidProtectionEnabled.ValueBool()
	}

	return result, diags
}

func expandActions(ctx context.Context, list types.List) ([]*discord.AutoModAction, diag.Diagnostics) {
	var diags diag.Diagnostics

	var actionModels []actionModel
	diags.Append(list.ElementsAs(ctx, &actionModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	actions := make([]*discord.AutoModAction, 0, len(actionModels))
	for _, am := range actionModels {
		action := &discord.AutoModAction{
			Type: int(am.Type.ValueInt64()),
		}

		if !am.Metadata.IsNull() && !am.Metadata.IsUnknown() {
			var meta actionMetadataModel
			diags.Append(am.Metadata.As(ctx, &meta, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return nil, diags
			}
			actionMeta := &discord.AutoModActionMeta{}
			if !meta.ChannelID.IsNull() && !meta.ChannelID.IsUnknown() {
				s := discord.Snowflake(meta.ChannelID.ValueString())
				actionMeta.ChannelID = &s
			}
			if !meta.DurationSeconds.IsNull() && !meta.DurationSeconds.IsUnknown() {
				v := int(meta.DurationSeconds.ValueInt64())
				actionMeta.DurationSeconds = &v
			}
			if !meta.CustomMessage.IsNull() && !meta.CustomMessage.IsUnknown() {
				v := meta.CustomMessage.ValueString()
				actionMeta.CustomMessage = &v
			}
			action.Metadata = actionMeta
		}

		actions = append(actions, action)
	}

	return actions, diags
}

func expandSnowflakeSet(ctx context.Context, set types.Set, diags *diag.Diagnostics) []discord.Snowflake {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	var vals []string
	diags.Append(set.ElementsAs(ctx, &vals, false)...)
	result := make([]discord.Snowflake, len(vals))
	for i, v := range vals {
		result[i] = discord.Snowflake(v)
	}
	return result
}

func flattenAutoModRule(ctx context.Context, rule *discord.AutoModerationRule, model *autoModerationRuleModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.ID = types.StringValue(rule.ID.String())
	model.GuildID = types.StringValue(rule.GuildID.String())
	model.Name = types.StringValue(rule.Name)
	model.EventType = types.Int64Value(int64(rule.EventType))
	model.TriggerType = types.Int64Value(int64(rule.TriggerType))
	model.Enabled = types.BoolValue(rule.Enabled)

	// Trigger metadata.
	if rule.TriggerMetadata != nil {
		tm := rule.TriggerMetadata

		keywordFilter, d := types.ListValueFrom(ctx, types.StringType, stringSliceOrEmpty(tm.KeywordFilter))
		diags.Append(d...)
		regexPatterns, d := types.ListValueFrom(ctx, types.StringType, stringSliceOrEmpty(tm.RegexPatterns))
		diags.Append(d...)
		presets, d := flattenIntSlice(ctx, tm.Presets)
		diags.Append(d...)
		allowList, d := types.ListValueFrom(ctx, types.StringType, stringSliceOrEmpty(tm.AllowList))
		diags.Append(d...)

		var mentionTotalLimit types.Int64
		if tm.MentionTotalLimit != nil {
			mentionTotalLimit = types.Int64Value(int64(*tm.MentionTotalLimit))
		} else {
			mentionTotalLimit = types.Int64Null()
		}

		mentionRaidProtection := types.BoolValue(tm.MentionRaidProtectionEnabled)

		tmObj, d := types.ObjectValue(triggerMetadataAttrTypes(), map[string]attr.Value{
			"keyword_filter":                  keywordFilter,
			"regex_patterns":                  regexPatterns,
			"presets":                         presets,
			"allow_list":                      allowList,
			"mention_total_limit":             mentionTotalLimit,
			"mention_raid_protection_enabled": mentionRaidProtection,
		})
		diags.Append(d...)
		model.TriggerMetadata = tmObj
	} else {
		model.TriggerMetadata = types.ObjectNull(triggerMetadataAttrTypes())
	}

	// Actions.
	actionVals := make([]attr.Value, 0, len(rule.Actions))
	for _, a := range rule.Actions {
		var metaObj types.Object
		if a.Metadata != nil {
			var chID types.String
			if a.Metadata.ChannelID != nil {
				chID = types.StringValue(a.Metadata.ChannelID.String())
			} else {
				chID = types.StringNull()
			}

			var durSec types.Int64
			if a.Metadata.DurationSeconds != nil {
				durSec = types.Int64Value(int64(*a.Metadata.DurationSeconds))
			} else {
				durSec = types.Int64Null()
			}

			var customMsg types.String
			if a.Metadata.CustomMessage != nil {
				customMsg = types.StringValue(*a.Metadata.CustomMessage)
			} else {
				customMsg = types.StringNull()
			}

			obj, d := types.ObjectValue(actionMetadataAttrTypes(), map[string]attr.Value{
				"channel_id":       chID,
				"duration_seconds": durSec,
				"custom_message":   customMsg,
			})
			diags.Append(d...)
			metaObj = obj
		} else {
			metaObj = types.ObjectNull(actionMetadataAttrTypes())
		}

		actionObj, d := types.ObjectValue(actionAttrTypes(), map[string]attr.Value{
			"type":     types.Int64Value(int64(a.Type)),
			"metadata": metaObj,
		})
		diags.Append(d...)
		actionVals = append(actionVals, actionObj)
	}
	actionsList, d := types.ListValue(types.ObjectType{AttrTypes: actionAttrTypes()}, actionVals)
	diags.Append(d...)
	model.Actions = actionsList

	// Exempt roles.
	model.ExemptRoles = flattenSnowflakeSet(ctx, rule.ExemptRoles, &diags)

	// Exempt channels.
	model.ExemptChannels = flattenSnowflakeSet(ctx, rule.ExemptChannels, &diags)

	return diags
}

func flattenSnowflakeSet(ctx context.Context, snowflakes []discord.Snowflake, diags *diag.Diagnostics) types.Set {
	if len(snowflakes) == 0 {
		return types.SetNull(types.StringType)
	}
	vals := make([]attr.Value, len(snowflakes))
	for i, s := range snowflakes {
		vals[i] = types.StringValue(s.String())
	}
	set, d := types.SetValue(types.StringType, vals)
	diags.Append(d...)
	return set
}

func stringSliceOrEmpty(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func flattenIntSlice(ctx context.Context, ints []int) (types.List, diag.Diagnostics) {
	vals := make([]attr.Value, len(ints))
	for i, v := range ints {
		vals[i] = types.Int64Value(int64(v))
	}
	return types.ListValue(types.Int64Type, vals)
}

