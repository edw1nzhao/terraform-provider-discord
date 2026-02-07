package guild

import (
	"context"
	"fmt"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &guildResource{}
	_ resource.ResourceWithConfigure   = &guildResource{}
	_ resource.ResourceWithImportState = &guildResource{}
)

// guildResource is the resource implementation.
type guildResource struct {
	client *discord.Client
}

// guildResourceModel maps the resource schema to a Go struct.
type guildResourceModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Icon                        types.String `tfsdk:"icon"`
	Splash                      types.String `tfsdk:"splash"`
	Banner                      types.String `tfsdk:"banner"`
	Description                 types.String `tfsdk:"description"`
	AFKChannelID                types.String `tfsdk:"afk_channel_id"`
	AFKTimeout                  types.Int64  `tfsdk:"afk_timeout"`
	VerificationLevel           types.Int64  `tfsdk:"verification_level"`
	DefaultMessageNotifications types.Int64  `tfsdk:"default_message_notifications"`
	ExplicitContentFilter       types.Int64  `tfsdk:"explicit_content_filter"`
	SystemChannelID             types.String `tfsdk:"system_channel_id"`
	SystemChannelFlags          types.Int64  `tfsdk:"system_channel_flags"`
	RulesChannelID              types.String `tfsdk:"rules_channel_id"`
	PublicUpdatesChannelID      types.String `tfsdk:"public_updates_channel_id"`
	PreferredLocale             types.String `tfsdk:"preferred_locale"`
	PremiumProgressBarEnabled   types.Bool   `tfsdk:"premium_progress_bar_enabled"`
	SafetyAlertsChannelID       types.String `tfsdk:"safety_alerts_channel_id"`
	OwnerID                     types.String `tfsdk:"owner_id"`
	PremiumTier                 types.Int64  `tfsdk:"premium_tier"`
	PremiumSubscriptionCount    types.Int64  `tfsdk:"premium_subscription_count"`
	Features                    types.List   `tfsdk:"features"`
	MFALevel                    types.Int64  `tfsdk:"mfa_level"`
}

// afkTimeoutValidator validates that the AFK timeout is one of the allowed values.
type afkTimeoutValidator struct{}

func (v afkTimeoutValidator) Description(_ context.Context) string {
	return "value must be one of: 60, 300, 900, 1800, 3600"
}

func (v afkTimeoutValidator) MarkdownDescription(_ context.Context) string {
	return "value must be one of: `60`, `300`, `900`, `1800`, `3600`"
}

func (v afkTimeoutValidator) ValidateInt64(_ context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	val := req.ConfigValue.ValueInt64()
	allowed := []int64{60, 300, 900, 1800, 3600}
	for _, a := range allowed {
		if val == a {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid AFK Timeout",
		fmt.Sprintf("AFK timeout must be one of 60, 300, 900, 1800, 3600, got: %d", val),
	)
}

// NewGuildResource returns a new guild resource.
func NewGuildResource() resource.Resource {
	return &guildResource{}
}

// Metadata returns the resource type name.
func (r *guildResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guild"
}

// Configure adds the provider configured client to the resource.
func (r *guildResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the resource.
func (r *guildResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord guild (server).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the guild (2-100 characters).",
				Required:    true,
			},
			"icon": schema.StringAttribute{
				Description: "The guild icon as a base64-encoded image data URI.",
				Optional:    true,
			},
			"splash": schema.StringAttribute{
				Description: "The guild splash image as a base64-encoded image data URI.",
				Optional:    true,
			},
			"banner": schema.StringAttribute{
				Description: "The guild banner image as a base64-encoded image data URI.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the guild.",
				Optional:    true,
			},
			"afk_channel_id": schema.StringAttribute{
				Description: "The ID of the AFK voice channel.",
				Optional:    true,
			},
			"afk_timeout": schema.Int64Attribute{
				Description: "AFK timeout in seconds. Must be one of: 60, 300, 900, 1800, 3600.",
				Optional:    true,
				Validators: []validator.Int64{
					afkTimeoutValidator{},
				},
			},
			"verification_level": schema.Int64Attribute{
				Description: "The verification level required for the guild (0-4).",
				Optional:    true,
			},
			"default_message_notifications": schema.Int64Attribute{
				Description: "The default message notification level (0 = all messages, 1 = only mentions).",
				Optional:    true,
			},
			"explicit_content_filter": schema.Int64Attribute{
				Description: "The explicit content filter level (0-2).",
				Optional:    true,
			},
			"system_channel_id": schema.StringAttribute{
				Description: "The ID of the channel where system messages are sent.",
				Optional:    true,
			},
			"system_channel_flags": schema.Int64Attribute{
				Description: "System channel flags.",
				Optional:    true,
			},
			"rules_channel_id": schema.StringAttribute{
				Description: "The ID of the channel where community rules are displayed.",
				Optional:    true,
			},
			"public_updates_channel_id": schema.StringAttribute{
				Description: "The ID of the channel where admins and moderators of community guilds receive notices from Discord.",
				Optional:    true,
			},
			"preferred_locale": schema.StringAttribute{
				Description: "The preferred locale of a community guild.",
				Optional:    true,
			},
			"premium_progress_bar_enabled": schema.BoolAttribute{
				Description: "Whether the guild has the boost progress bar enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"safety_alerts_channel_id": schema.StringAttribute{
				Description: "The ID of the channel where safety alerts are sent.",
				Optional:    true,
			},
			"owner_id": schema.StringAttribute{
				Description: "The ID of the guild owner.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"premium_tier": schema.Int64Attribute{
				Description: "The premium tier (Server Boost level).",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"premium_subscription_count": schema.Int64Attribute{
				Description: "The number of boosts this guild currently has.",
				Computed:    true,
			},
			"features": schema.ListAttribute{
				Description: "The list of enabled guild features.",
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"mfa_level": schema.Int64Attribute{
				Description: "The required MFA level for the guild.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the guild resource.
func (r *guildResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan guildResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.CreateGuildParams{
		Name: plan.Name.ValueString(),
	}

	if !plan.Icon.IsNull() && !plan.Icon.IsUnknown() {
		v := plan.Icon.ValueString()
		params.Icon = &v
	}
	if !plan.VerificationLevel.IsNull() && !plan.VerificationLevel.IsUnknown() {
		v := int(plan.VerificationLevel.ValueInt64())
		params.VerificationLevel = &v
	}
	if !plan.DefaultMessageNotifications.IsNull() && !plan.DefaultMessageNotifications.IsUnknown() {
		v := int(plan.DefaultMessageNotifications.ValueInt64())
		params.DefaultMessageNotifications = &v
	}
	if !plan.ExplicitContentFilter.IsNull() && !plan.ExplicitContentFilter.IsUnknown() {
		v := int(plan.ExplicitContentFilter.ValueInt64())
		params.ExplicitContentFilter = &v
	}
	if !plan.AFKChannelID.IsNull() && !plan.AFKChannelID.IsUnknown() {
		v := discord.Snowflake(plan.AFKChannelID.ValueString())
		params.AFKChannelID = &v
	}
	if !plan.AFKTimeout.IsNull() && !plan.AFKTimeout.IsUnknown() {
		v := int(plan.AFKTimeout.ValueInt64())
		params.AFKTimeout = &v
	}
	if !plan.SystemChannelID.IsNull() && !plan.SystemChannelID.IsUnknown() {
		v := discord.Snowflake(plan.SystemChannelID.ValueString())
		params.SystemChannelID = &v
	}
	if !plan.SystemChannelFlags.IsNull() && !plan.SystemChannelFlags.IsUnknown() {
		v := int(plan.SystemChannelFlags.ValueInt64())
		params.SystemChannelFlags = &v
	}

	guild, err := r.client.CreateGuild(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Discord Guild",
			"Could not create guild: "+err.Error(),
		)
		return
	}

	// Some fields require a follow-up Modify since CreateGuild doesn't accept them.
	plan.ID = types.StringValue(guild.ID.String())

	// If any modify-only fields are set, issue a ModifyGuild call.
	if needsPostCreateModify(plan) {
		modifyParams := buildModifyGuildParams(plan)
		guild, err = r.client.ModifyGuild(ctx, guild.ID, modifyParams)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Configuring Discord Guild",
				"Guild was created but some settings could not be applied: "+err.Error(),
			)
			return
		}
	}

	mapGuildToState(ctx, guild, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *guildResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state guildResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guild, err := r.client.GetGuild(ctx, discord.Snowflake(state.ID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Discord Guild",
			"Could not read guild ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	mapGuildToState(ctx, guild, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update modifies the guild resource.
func (r *guildResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan guildResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state guildResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := buildModifyGuildParams(plan)

	guild, err := r.client.ModifyGuild(ctx, discord.Snowflake(state.ID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Discord Guild",
			"Could not update guild ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.ID = state.ID
	mapGuildToState(ctx, guild, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the guild resource.
func (r *guildResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state guildResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGuild(ctx, discord.Snowflake(state.ID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Discord Guild",
			"Could not delete guild ID "+state.ID.ValueString()+": "+err.Error(),
		)
	}
}

// ImportState allows importing an existing guild by its ID.
func (r *guildResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// needsPostCreateModify returns true if any fields not supported by CreateGuild
// are set in the plan.
func needsPostCreateModify(plan guildResourceModel) bool {
	return !plan.Splash.IsNull() ||
		!plan.Banner.IsNull() ||
		!plan.Description.IsNull() ||
		!plan.RulesChannelID.IsNull() ||
		!plan.PublicUpdatesChannelID.IsNull() ||
		!plan.PreferredLocale.IsNull() ||
		!plan.PremiumProgressBarEnabled.IsNull() ||
		!plan.SafetyAlertsChannelID.IsNull()
}

// buildModifyGuildParams constructs the ModifyGuildParams from the plan.
func buildModifyGuildParams(plan guildResourceModel) *discord.ModifyGuildParams {
	params := &discord.ModifyGuildParams{}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		v := plan.Name.ValueString()
		params.Name = &v
	}
	if !plan.Icon.IsNull() {
		v := plan.Icon.ValueString()
		params.Icon = &v
	}
	if !plan.Splash.IsNull() {
		v := plan.Splash.ValueString()
		params.Splash = &v
	}
	if !plan.Banner.IsNull() {
		v := plan.Banner.ValueString()
		params.Banner = &v
	}
	if !plan.Description.IsNull() {
		v := plan.Description.ValueString()
		params.Description = &v
	}
	if !plan.AFKChannelID.IsNull() {
		v := discord.Snowflake(plan.AFKChannelID.ValueString())
		params.AFKChannelID = &v
	}
	if !plan.AFKTimeout.IsNull() && !plan.AFKTimeout.IsUnknown() {
		v := int(plan.AFKTimeout.ValueInt64())
		params.AFKTimeout = &v
	}
	if !plan.VerificationLevel.IsNull() && !plan.VerificationLevel.IsUnknown() {
		v := int(plan.VerificationLevel.ValueInt64())
		params.VerificationLevel = &v
	}
	if !plan.DefaultMessageNotifications.IsNull() && !plan.DefaultMessageNotifications.IsUnknown() {
		v := int(plan.DefaultMessageNotifications.ValueInt64())
		params.DefaultMessageNotifications = &v
	}
	if !plan.ExplicitContentFilter.IsNull() && !plan.ExplicitContentFilter.IsUnknown() {
		v := int(plan.ExplicitContentFilter.ValueInt64())
		params.ExplicitContentFilter = &v
	}
	if !plan.SystemChannelID.IsNull() {
		v := discord.Snowflake(plan.SystemChannelID.ValueString())
		params.SystemChannelID = &v
	}
	if !plan.SystemChannelFlags.IsNull() && !plan.SystemChannelFlags.IsUnknown() {
		v := int(plan.SystemChannelFlags.ValueInt64())
		params.SystemChannelFlags = &v
	}
	if !plan.RulesChannelID.IsNull() {
		v := discord.Snowflake(plan.RulesChannelID.ValueString())
		params.RulesChannelID = &v
	}
	if !plan.PublicUpdatesChannelID.IsNull() {
		v := discord.Snowflake(plan.PublicUpdatesChannelID.ValueString())
		params.PublicUpdatesChannelID = &v
	}
	if !plan.PreferredLocale.IsNull() {
		v := plan.PreferredLocale.ValueString()
		params.PreferredLocale = &v
	}
	if !plan.PremiumProgressBarEnabled.IsNull() && !plan.PremiumProgressBarEnabled.IsUnknown() {
		v := plan.PremiumProgressBarEnabled.ValueBool()
		params.PremiumProgressBarEnabled = &v
	}
	if !plan.SafetyAlertsChannelID.IsNull() {
		v := discord.Snowflake(plan.SafetyAlertsChannelID.ValueString())
		params.SafetyAlertsChannelID = &v
	}

	return params
}

// mapGuildToState maps a Discord Guild API response to the Terraform state model.
func mapGuildToState(ctx context.Context, guild *discord.Guild, state *guildResourceModel) {
	state.ID = types.StringValue(guild.ID.String())
	state.Name = types.StringValue(guild.Name)
	state.OwnerID = types.StringValue(guild.OwnerID.String())
	state.PremiumTier = types.Int64Value(int64(guild.PremiumTier))
	state.PremiumSubscriptionCount = types.Int64Value(int64(guild.PremiumSubscriptionCount))
	state.MFALevel = types.Int64Value(int64(guild.MFALevel))
	state.VerificationLevel = types.Int64Value(int64(guild.VerificationLevel))
	state.DefaultMessageNotifications = types.Int64Value(int64(guild.DefaultMessageNotifications))
	state.ExplicitContentFilter = types.Int64Value(int64(guild.ExplicitContentFilter))
	state.AFKTimeout = types.Int64Value(int64(guild.AFKTimeout))
	state.SystemChannelFlags = types.Int64Value(int64(guild.SystemChannelFlags))
	state.PremiumProgressBarEnabled = types.BoolValue(guild.PremiumProgressBarEnabled)
	state.PreferredLocale = types.StringValue(guild.PreferredLocale)

	// Nullable string fields.
	if guild.Icon != nil {
		state.Icon = types.StringValue(*guild.Icon)
	} else {
		state.Icon = types.StringNull()
	}
	if guild.Splash != nil {
		state.Splash = types.StringValue(*guild.Splash)
	} else {
		state.Splash = types.StringNull()
	}
	if guild.Banner != nil {
		state.Banner = types.StringValue(*guild.Banner)
	} else {
		state.Banner = types.StringNull()
	}
	if guild.Description != nil {
		state.Description = types.StringValue(*guild.Description)
	} else {
		state.Description = types.StringNull()
	}

	// Nullable snowflake fields.
	if guild.AFKChannelID != nil {
		state.AFKChannelID = types.StringValue(guild.AFKChannelID.String())
	} else {
		state.AFKChannelID = types.StringNull()
	}
	if guild.SystemChannelID != nil {
		state.SystemChannelID = types.StringValue(guild.SystemChannelID.String())
	} else {
		state.SystemChannelID = types.StringNull()
	}
	if guild.RulesChannelID != nil {
		state.RulesChannelID = types.StringValue(guild.RulesChannelID.String())
	} else {
		state.RulesChannelID = types.StringNull()
	}
	if guild.PublicUpdatesChannelID != nil {
		state.PublicUpdatesChannelID = types.StringValue(guild.PublicUpdatesChannelID.String())
	} else {
		state.PublicUpdatesChannelID = types.StringNull()
	}
	if guild.SafetyAlertsChannelID != nil {
		state.SafetyAlertsChannelID = types.StringValue(guild.SafetyAlertsChannelID.String())
	} else {
		state.SafetyAlertsChannelID = types.StringNull()
	}

	// Features list.
	features := make([]types.String, len(guild.Features))
	for i, f := range guild.Features {
		features[i] = types.StringValue(f)
	}
	var diags diag.Diagnostics
	state.Features, diags = types.ListValueFrom(ctx, types.StringType, features)
	_ = diags // TODO: propagate diagnostics to caller
}
