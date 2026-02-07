package welcome_screen

import (
	"context"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure interface compliance.
var (
	_ resource.Resource                = &welcomeScreenResource{}
	_ resource.ResourceWithConfigure   = &welcomeScreenResource{}
	_ resource.ResourceWithImportState = &welcomeScreenResource{}
)

// NewWelcomeScreenResource returns a new resource for discord_welcome_screen.
func NewWelcomeScreenResource() resource.Resource {
	return &welcomeScreenResource{}
}

type welcomeScreenResource struct {
	client *discord.Client
}

// welcomeScreenModel maps the Terraform schema to Go types.
type welcomeScreenModel struct {
	GuildID         types.String `tfsdk:"guild_id"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Description     types.String `tfsdk:"description"`
	WelcomeChannels types.List   `tfsdk:"welcome_channels"`
}

// welcomeChannelModel maps a single welcome channel entry.
type welcomeChannelModel struct {
	ChannelID   types.String `tfsdk:"channel_id"`
	Description types.String `tfsdk:"description"`
	EmojiID     types.String `tfsdk:"emoji_id"`
	EmojiName   types.String `tfsdk:"emoji_name"`
}

func welcomeChannelAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"channel_id":  types.StringType,
		"description": types.StringType,
		"emoji_id":    types.StringType,
		"emoji_name":  types.StringType,
	}
}

// Metadata sets the type name.
func (r *welcomeScreenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_welcome_screen"
}

// Schema defines the schema.
func (r *welcomeScreenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord guild welcome screen.",
		Attributes: map[string]schema.Attribute{
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild. Acts as the resource ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the welcome screen is enabled.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The server description shown in the welcome screen.",
				Optional:    true,
			},
			"welcome_channels": schema.ListNestedAttribute{
				Description: "Channels shown in the welcome screen (max 5).",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"channel_id": schema.StringAttribute{
							Description: "The channel ID.",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description shown for this channel.",
							Required:    true,
						},
						"emoji_id": schema.StringAttribute{
							Description: "The emoji ID, if the emoji is custom.",
							Optional:    true,
						},
						"emoji_name": schema.StringAttribute{
							Description: "The emoji name if custom, or the unicode character.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Configure stores the provider data.
func (r *welcomeScreenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create applies the welcome screen settings via PATCH (the welcome screen always exists as part of the guild).
func (r *welcomeScreenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan welcomeScreenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params, diags := expandWelcomeScreenParams(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ws, err := r.client.ModifyGuildWelcomeScreen(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating welcome screen", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenWelcomeScreen(ctx, ws, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state from the API.
func (r *welcomeScreenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state welcomeScreenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ws, err := r.client.GetGuildWelcomeScreen(ctx, discord.Snowflake(state.GuildID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading welcome screen", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenWelcomeScreen(ctx, ws, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update modifies the welcome screen.
func (r *welcomeScreenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan welcomeScreenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params, diags := expandWelcomeScreenParams(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ws, err := r.client.ModifyGuildWelcomeScreen(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating welcome screen", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenWelcomeScreen(ctx, ws, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete disables the welcome screen and clears channels.
func (r *welcomeScreenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state welcomeScreenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enabled := false
	params := &discord.ModifyWelcomeScreenParams{
		Enabled:         &enabled,
		WelcomeChannels: []*discord.WelcomeScreenChannel{},
	}

	_, err := r.client.ModifyGuildWelcomeScreen(ctx, discord.Snowflake(state.GuildID.ValueString()), params)
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting welcome screen", err.Error())
	}
}

// ImportState supports importing by guild_id.
func (r *welcomeScreenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), req.ID)...)
}

// --- Helper functions ---

func expandWelcomeScreenParams(ctx context.Context, model *welcomeScreenModel) (*discord.ModifyWelcomeScreenParams, diag.Diagnostics) {
	var diags diag.Diagnostics

	enabled := model.Enabled.ValueBool()
	params := &discord.ModifyWelcomeScreenParams{
		Enabled: &enabled,
	}

	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		v := model.Description.ValueString()
		params.Description = &v
	}

	var channelModels []welcomeChannelModel
	diags.Append(model.WelcomeChannels.ElementsAs(ctx, &channelModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	channels := make([]*discord.WelcomeScreenChannel, 0, len(channelModels))
	for _, cm := range channelModels {
		ch := &discord.WelcomeScreenChannel{
			ChannelID:   discord.Snowflake(cm.ChannelID.ValueString()),
			Description: cm.Description.ValueString(),
		}
		if !cm.EmojiID.IsNull() && !cm.EmojiID.IsUnknown() {
			s := discord.Snowflake(cm.EmojiID.ValueString())
			ch.EmojiID = &s
		}
		if !cm.EmojiName.IsNull() && !cm.EmojiName.IsUnknown() {
			v := cm.EmojiName.ValueString()
			ch.EmojiName = &v
		}
		channels = append(channels, ch)
	}
	params.WelcomeChannels = channels

	return params, diags
}

func flattenWelcomeScreen(ctx context.Context, ws *discord.WelcomeScreen, model *welcomeScreenModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// guild_id is not returned by the API so we preserve the existing value.

	if ws.Description != nil {
		model.Description = types.StringValue(*ws.Description)
	} else {
		model.Description = types.StringNull()
	}

	// We do not know the enabled state from the WelcomeScreen struct directly;
	// the caller sets it from context. For Read, we consider the screen enabled
	// if it has channels. On Create/Update we already set it.
	// Note: the API response for ModifyGuildWelcomeScreen does not include an "enabled" field,
	// so we preserve the plan/state value. For Read, we infer from the data available.

	channelVals := make([]attr.Value, 0, len(ws.WelcomeChannels))
	for _, ch := range ws.WelcomeChannels {
		var emojiID types.String
		if ch.EmojiID != nil {
			emojiID = types.StringValue(ch.EmojiID.String())
		} else {
			emojiID = types.StringNull()
		}

		var emojiName types.String
		if ch.EmojiName != nil {
			emojiName = types.StringValue(*ch.EmojiName)
		} else {
			emojiName = types.StringNull()
		}

		obj, d := types.ObjectValue(welcomeChannelAttrTypes(), map[string]attr.Value{
			"channel_id":  types.StringValue(ch.ChannelID.String()),
			"description": types.StringValue(ch.Description),
			"emoji_id":    emojiID,
			"emoji_name":  emojiName,
		})
		diags.Append(d...)
		channelVals = append(channelVals, obj)
	}

	channelsList, d := types.ListValue(types.ObjectType{AttrTypes: welcomeChannelAttrTypes()}, channelVals)
	diags.Append(d...)
	model.WelcomeChannels = channelsList

	return diags
}
