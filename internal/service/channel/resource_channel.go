package channel

import (
	"context"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &channelResource{}
	_ resource.ResourceWithConfigure   = &channelResource{}
	_ resource.ResourceWithImportState = &channelResource{}
)

// channelResource is the resource implementation.
type channelResource struct {
	client *discord.Client
}

// channelResourceModel maps the resource schema to a Go struct.
type channelResourceModel struct {
	ID                            types.String `tfsdk:"id"`
	GuildID                       types.String `tfsdk:"guild_id"`
	Name                          types.String `tfsdk:"name"`
	Type                          types.Int64  `tfsdk:"type"`
	Position                      types.Int64  `tfsdk:"position"`
	Topic                         types.String `tfsdk:"topic"`
	NSFW                          types.Bool   `tfsdk:"nsfw"`
	RateLimitPerUser              types.Int64  `tfsdk:"rate_limit_per_user"`
	Bitrate                       types.Int64  `tfsdk:"bitrate"`
	UserLimit                     types.Int64  `tfsdk:"user_limit"`
	ParentID                      types.String `tfsdk:"parent_id"`
	RTCRegion                     types.String `tfsdk:"rtc_region"`
	VideoQualityMode              types.Int64  `tfsdk:"video_quality_mode"`
	DefaultAutoArchiveDuration    types.Int64  `tfsdk:"default_auto_archive_duration"`
	DefaultThreadRateLimitPerUser types.Int64  `tfsdk:"default_thread_rate_limit_per_user"`
	DefaultSortOrder              types.Int64  `tfsdk:"default_sort_order"`
	DefaultForumLayout            types.Int64  `tfsdk:"default_forum_layout"`
}

// NewChannelResource returns a new channel resource.
func NewChannelResource() resource.Resource {
	return &channelResource{}
}

// Metadata returns the resource type name.
func (r *channelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

// Configure adds the provider configured client to the resource.
func (r *channelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the resource.
func (r *channelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Discord channels (text, voice, category, announcement, stage, forum, media).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the channel.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild this channel belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the channel (1-100 characters).",
				Required:    true,
			},
			"type": schema.Int64Attribute{
				Description: "The type of channel (0=text, 2=voice, 4=category, 5=announcement, 13=stage, 15=forum, 16=media).",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"position": schema.Int64Attribute{
				Description: "The sorting position of the channel.",
				Optional:    true,
				Computed:    true,
			},
			"topic": schema.StringAttribute{
				Description: "The channel topic (0-4096 characters for forum channels, 0-1024 for others).",
				Optional:    true,
			},
			"nsfw": schema.BoolAttribute{
				Description: "Whether the channel is NSFW.",
				Optional:    true,
				Computed:    true,
			},
			"rate_limit_per_user": schema.Int64Attribute{
				Description: "Slowmode rate limit in seconds (0-21600). Users can send one message per this interval.",
				Optional:    true,
				Computed:    true,
			},
			"bitrate": schema.Int64Attribute{
				Description: "The bitrate (in bits) of the voice channel.",
				Optional:    true,
				Computed:    true,
			},
			"user_limit": schema.Int64Attribute{
				Description: "The user limit of the voice channel (0 for no limit).",
				Optional:    true,
				Computed:    true,
			},
			"parent_id": schema.StringAttribute{
				Description: "The ID of the parent category for a channel.",
				Optional:    true,
			},
			"rtc_region": schema.StringAttribute{
				Description: "Voice region ID for the voice channel. Automatic when set to null.",
				Optional:    true,
			},
			"video_quality_mode": schema.Int64Attribute{
				Description: "The camera video quality mode of the voice channel (1=auto, 2=720p).",
				Optional:    true,
				Computed:    true,
			},
			"default_auto_archive_duration": schema.Int64Attribute{
				Description: "Default duration in minutes for threads to auto-archive (60, 1440, 4320, 10080).",
				Optional:    true,
				Computed:    true,
			},
			"default_thread_rate_limit_per_user": schema.Int64Attribute{
				Description: "Default slowmode for threads created in this channel (0-21600 seconds).",
				Optional:    true,
				Computed:    true,
			},
			"default_sort_order": schema.Int64Attribute{
				Description: "Default sort order for forum channels (0=latest_activity, 1=creation_date).",
				Optional:    true,
			},
			"default_forum_layout": schema.Int64Attribute{
				Description: "Default layout for forum channels (0=not_set, 1=list_view, 2=gallery_view).",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// Create creates the channel resource.
func (r *channelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan channelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channelType := int(plan.Type.ValueInt64())
	params := &discord.CreateChannelParams{
		Name: plan.Name.ValueString(),
		Type: &channelType,
	}

	if !plan.Topic.IsNull() && !plan.Topic.IsUnknown() {
		v := plan.Topic.ValueString()
		params.Topic = &v
	}
	if !plan.Bitrate.IsNull() && !plan.Bitrate.IsUnknown() {
		v := int(plan.Bitrate.ValueInt64())
		params.Bitrate = &v
	}
	if !plan.UserLimit.IsNull() && !plan.UserLimit.IsUnknown() {
		v := int(plan.UserLimit.ValueInt64())
		params.UserLimit = &v
	}
	if !plan.RateLimitPerUser.IsNull() && !plan.RateLimitPerUser.IsUnknown() {
		v := int(plan.RateLimitPerUser.ValueInt64())
		params.RateLimitPerUser = &v
	}
	if !plan.Position.IsNull() && !plan.Position.IsUnknown() {
		v := int(plan.Position.ValueInt64())
		params.Position = &v
	}
	if !plan.ParentID.IsNull() && !plan.ParentID.IsUnknown() {
		v := discord.Snowflake(plan.ParentID.ValueString())
		params.ParentID = &v
	}
	if !plan.NSFW.IsNull() && !plan.NSFW.IsUnknown() {
		v := plan.NSFW.ValueBool()
		params.NSFW = &v
	}
	if !plan.RTCRegion.IsNull() && !plan.RTCRegion.IsUnknown() {
		v := plan.RTCRegion.ValueString()
		params.RTCRegion = &v
	}
	if !plan.VideoQualityMode.IsNull() && !plan.VideoQualityMode.IsUnknown() {
		v := int(plan.VideoQualityMode.ValueInt64())
		params.VideoQualityMode = &v
	}
	if !plan.DefaultAutoArchiveDuration.IsNull() && !plan.DefaultAutoArchiveDuration.IsUnknown() {
		v := int(plan.DefaultAutoArchiveDuration.ValueInt64())
		params.DefaultAutoArchiveDuration = &v
	}
	if !plan.DefaultThreadRateLimitPerUser.IsNull() && !plan.DefaultThreadRateLimitPerUser.IsUnknown() {
		v := int(plan.DefaultThreadRateLimitPerUser.ValueInt64())
		params.DefaultThreadRateLimitPerUser = &v
	}
	if !plan.DefaultSortOrder.IsNull() && !plan.DefaultSortOrder.IsUnknown() {
		v := int(plan.DefaultSortOrder.ValueInt64())
		params.DefaultSortOrder = &v
	}
	if !plan.DefaultForumLayout.IsNull() && !plan.DefaultForumLayout.IsUnknown() {
		v := int(plan.DefaultForumLayout.ValueInt64())
		params.DefaultForumLayout = &v
	}

	guildID := discord.Snowflake(plan.GuildID.ValueString())
	ch, err := r.client.CreateGuildChannel(ctx, guildID, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Discord Channel",
			"Could not create channel: "+err.Error(),
		)
		return
	}

	mapChannelToState(ch, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *channelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state channelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ch, err := r.client.GetChannel(ctx, discord.Snowflake(state.ID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Discord Channel",
			"Could not read channel ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	mapChannelToState(ch, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update modifies the channel resource.
func (r *channelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan channelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state channelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.ModifyChannelParams{}

	name := plan.Name.ValueString()
	params.Name = &name

	if !plan.Topic.IsNull() {
		v := plan.Topic.ValueString()
		params.Topic = &v
	}
	if !plan.Position.IsNull() && !plan.Position.IsUnknown() {
		v := int(plan.Position.ValueInt64())
		params.Position = &v
	}
	if !plan.NSFW.IsNull() && !plan.NSFW.IsUnknown() {
		v := plan.NSFW.ValueBool()
		params.NSFW = &v
	}
	if !plan.RateLimitPerUser.IsNull() && !plan.RateLimitPerUser.IsUnknown() {
		v := int(plan.RateLimitPerUser.ValueInt64())
		params.RateLimitPerUser = &v
	}
	if !plan.Bitrate.IsNull() && !plan.Bitrate.IsUnknown() {
		v := int(plan.Bitrate.ValueInt64())
		params.Bitrate = &v
	}
	if !plan.UserLimit.IsNull() && !plan.UserLimit.IsUnknown() {
		v := int(plan.UserLimit.ValueInt64())
		params.UserLimit = &v
	}
	if !plan.ParentID.IsNull() {
		v := discord.Snowflake(plan.ParentID.ValueString())
		params.ParentID = &v
	}
	if !plan.RTCRegion.IsNull() {
		v := plan.RTCRegion.ValueString()
		params.RTCRegion = &v
	}
	if !plan.VideoQualityMode.IsNull() && !plan.VideoQualityMode.IsUnknown() {
		v := int(plan.VideoQualityMode.ValueInt64())
		params.VideoQualityMode = &v
	}
	if !plan.DefaultAutoArchiveDuration.IsNull() && !plan.DefaultAutoArchiveDuration.IsUnknown() {
		v := int(plan.DefaultAutoArchiveDuration.ValueInt64())
		params.DefaultAutoArchiveDuration = &v
	}
	if !plan.DefaultThreadRateLimitPerUser.IsNull() && !plan.DefaultThreadRateLimitPerUser.IsUnknown() {
		v := int(plan.DefaultThreadRateLimitPerUser.ValueInt64())
		params.DefaultThreadRateLimitPerUser = &v
	}
	if !plan.DefaultSortOrder.IsNull() {
		v := int(plan.DefaultSortOrder.ValueInt64())
		params.DefaultSortOrder = &v
	}
	if !plan.DefaultForumLayout.IsNull() && !plan.DefaultForumLayout.IsUnknown() {
		v := int(plan.DefaultForumLayout.ValueInt64())
		params.DefaultForumLayout = &v
	}

	ch, err := r.client.ModifyChannel(ctx, discord.Snowflake(state.ID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Discord Channel",
			"Could not update channel ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	mapChannelToState(ch, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the channel resource.
func (r *channelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state channelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteChannel(ctx, discord.Snowflake(state.ID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Discord Channel",
			"Could not delete channel ID "+state.ID.ValueString()+": "+err.Error(),
		)
	}
}

// ImportState allows importing an existing channel by its ID.
func (r *channelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapChannelToState maps a Discord Channel API response to the Terraform state model.
func mapChannelToState(ch *discord.Channel, state *channelResourceModel) {
	state.ID = types.StringValue(ch.ID.String())
	state.Type = types.Int64Value(int64(ch.Type))
	state.NSFW = types.BoolValue(ch.NSFW)

	if ch.GuildID != nil {
		state.GuildID = types.StringValue(ch.GuildID.String())
	}
	if ch.Name != nil {
		state.Name = types.StringValue(*ch.Name)
	}

	// Position.
	if ch.Position != nil {
		state.Position = types.Int64Value(int64(*ch.Position))
	} else {
		state.Position = types.Int64Null()
	}

	// Topic.
	if ch.Topic != nil {
		state.Topic = types.StringValue(*ch.Topic)
	} else {
		state.Topic = types.StringNull()
	}

	// Rate limit per user.
	if ch.RateLimitPerUser != nil {
		state.RateLimitPerUser = types.Int64Value(int64(*ch.RateLimitPerUser))
	} else {
		state.RateLimitPerUser = types.Int64Value(0)
	}

	// Bitrate.
	if ch.Bitrate != nil {
		state.Bitrate = types.Int64Value(int64(*ch.Bitrate))
	} else {
		state.Bitrate = types.Int64Null()
	}

	// User limit.
	if ch.UserLimit != nil {
		state.UserLimit = types.Int64Value(int64(*ch.UserLimit))
	} else {
		state.UserLimit = types.Int64Null()
	}

	// Parent ID.
	if ch.ParentID != nil {
		state.ParentID = types.StringValue(ch.ParentID.String())
	} else {
		state.ParentID = types.StringNull()
	}

	// RTC Region.
	if ch.RTCRegion != nil {
		state.RTCRegion = types.StringValue(*ch.RTCRegion)
	} else {
		state.RTCRegion = types.StringNull()
	}

	// Video quality mode.
	if ch.VideoQualityMode != nil {
		state.VideoQualityMode = types.Int64Value(int64(*ch.VideoQualityMode))
	} else {
		state.VideoQualityMode = types.Int64Null()
	}

	// Default auto archive duration.
	if ch.DefaultAutoArchiveDuration != nil {
		state.DefaultAutoArchiveDuration = types.Int64Value(int64(*ch.DefaultAutoArchiveDuration))
	} else {
		state.DefaultAutoArchiveDuration = types.Int64Null()
	}

	// Default thread rate limit per user.
	if ch.DefaultThreadRateLimitPerUser != nil {
		state.DefaultThreadRateLimitPerUser = types.Int64Value(int64(*ch.DefaultThreadRateLimitPerUser))
	} else {
		state.DefaultThreadRateLimitPerUser = types.Int64Value(0)
	}

	// Default sort order.
	if ch.DefaultSortOrder != nil {
		state.DefaultSortOrder = types.Int64Value(int64(*ch.DefaultSortOrder))
	} else {
		state.DefaultSortOrder = types.Int64Null()
	}

	// Default forum layout.
	if ch.DefaultForumLayout != nil {
		state.DefaultForumLayout = types.Int64Value(int64(*ch.DefaultForumLayout))
	} else {
		state.DefaultForumLayout = types.Int64Value(0)
	}
}
