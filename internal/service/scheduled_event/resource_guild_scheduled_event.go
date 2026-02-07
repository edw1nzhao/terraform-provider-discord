package scheduled_event

import (
	"context"
	"strings"
	"time"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure interface compliance.
var (
	_ resource.Resource                = &guildScheduledEventResource{}
	_ resource.ResourceWithConfigure   = &guildScheduledEventResource{}
	_ resource.ResourceWithImportState = &guildScheduledEventResource{}
)

// NewGuildScheduledEventResource returns a new resource for discord_guild_scheduled_event.
func NewGuildScheduledEventResource() resource.Resource {
	return &guildScheduledEventResource{}
}

type guildScheduledEventResource struct {
	client *discord.Client
}

// guildScheduledEventModel maps the Terraform schema to Go types.
type guildScheduledEventModel struct {
	ID                     types.String `tfsdk:"id"`
	GuildID                types.String `tfsdk:"guild_id"`
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	ScheduledStartTime     types.String `tfsdk:"scheduled_start_time"`
	ScheduledEndTime       types.String `tfsdk:"scheduled_end_time"`
	EntityType             types.Int64  `tfsdk:"entity_type"`
	ChannelID              types.String `tfsdk:"channel_id"`
	EntityMetadataLocation types.String `tfsdk:"entity_metadata_location"`
	PrivacyLevel           types.Int64  `tfsdk:"privacy_level"`
	Status                 types.Int64  `tfsdk:"status"`
	Image                  types.String `tfsdk:"image"`
}

// Metadata sets the type name.
func (r *guildScheduledEventResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guild_scheduled_event"
}

// Schema defines the schema.
func (r *guildScheduledEventResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord guild scheduled event.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the scheduled event.",
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
				Description: "The name of the scheduled event (1-100 characters).",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the scheduled event (1-1000 characters).",
				Optional:    true,
			},
			"scheduled_start_time": schema.StringAttribute{
				Description: "The scheduled start time in ISO8601 format.",
				Required:    true,
			},
			"scheduled_end_time": schema.StringAttribute{
				Description: "The scheduled end time in ISO8601 format. Required for EXTERNAL events.",
				Optional:    true,
			},
			"entity_type": schema.Int64Attribute{
				Description: "The entity type (1 = STAGE_INSTANCE, 2 = VOICE, 3 = EXTERNAL).",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"channel_id": schema.StringAttribute{
				Description: "The channel ID. Required for STAGE_INSTANCE and VOICE entity types.",
				Optional:    true,
			},
			"entity_metadata_location": schema.StringAttribute{
				Description: "The location of the event. Required for EXTERNAL entity type.",
				Optional:    true,
			},
			"privacy_level": schema.Int64Attribute{
				Description: "The privacy level (2 = GUILD_ONLY).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(discord.ScheduledEventPrivacyGuildOnly),
			},
			"status": schema.Int64Attribute{
				Description: "The status of the scheduled event (1 = SCHEDULED, 2 = ACTIVE, 3 = COMPLETED, 4 = CANCELED).",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"image": schema.StringAttribute{
				Description: "The cover image for the event (data URI or URL).",
				Optional:    true,
			},
		},
	}
}

// Configure stores the provider data.
func (r *guildScheduledEventResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create creates a new guild scheduled event.
func (r *guildScheduledEventResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan guildScheduledEventModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	startTime, err := time.Parse(time.RFC3339, plan.ScheduledStartTime.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid scheduled_start_time", "Must be a valid ISO8601/RFC3339 timestamp: "+err.Error())
		return
	}

	params := &discord.CreateScheduledEventParams{
		Name:               plan.Name.ValueString(),
		PrivacyLevel:       int(plan.PrivacyLevel.ValueInt64()),
		ScheduledStartTime: startTime,
		EntityType:         int(plan.EntityType.ValueInt64()),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		v := plan.Description.ValueString()
		params.Description = &v
	}

	if !plan.ScheduledEndTime.IsNull() && !plan.ScheduledEndTime.IsUnknown() {
		endTime, err := time.Parse(time.RFC3339, plan.ScheduledEndTime.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid scheduled_end_time", "Must be a valid ISO8601/RFC3339 timestamp: "+err.Error())
			return
		}
		params.ScheduledEndTime = &endTime
	}

	if !plan.ChannelID.IsNull() && !plan.ChannelID.IsUnknown() {
		s := discord.Snowflake(plan.ChannelID.ValueString())
		params.ChannelID = &s
	}

	if !plan.EntityMetadataLocation.IsNull() && !plan.EntityMetadataLocation.IsUnknown() {
		loc := plan.EntityMetadataLocation.ValueString()
		params.EntityMetadata = &discord.ScheduledEventEntityMetadata{
			Location: &loc,
		}
	}

	if !plan.Image.IsNull() && !plan.Image.IsUnknown() {
		v := plan.Image.ValueString()
		params.Image = &v
	}

	event, err := r.client.CreateGuildScheduledEvent(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating guild scheduled event", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenScheduledEvent(event, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state from the API.
func (r *guildScheduledEventResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state guildScheduledEventModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	event, err := r.client.GetGuildScheduledEvent(
		ctx,
		discord.Snowflake(state.GuildID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
	)
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading guild scheduled event", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenScheduledEvent(event, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update modifies an existing guild scheduled event.
func (r *guildScheduledEventResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan guildScheduledEventModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	privacyLevel := int(plan.PrivacyLevel.ValueInt64())

	params := &discord.ModifyScheduledEventParams{
		Name:         &name,
		PrivacyLevel: &privacyLevel,
	}

	startTime, err := time.Parse(time.RFC3339, plan.ScheduledStartTime.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid scheduled_start_time", "Must be a valid ISO8601/RFC3339 timestamp: "+err.Error())
		return
	}
	params.ScheduledStartTime = &startTime

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		v := plan.Description.ValueString()
		params.Description = &v
	}

	if !plan.ScheduledEndTime.IsNull() && !plan.ScheduledEndTime.IsUnknown() {
		endTime, err := time.Parse(time.RFC3339, plan.ScheduledEndTime.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid scheduled_end_time", "Must be a valid ISO8601/RFC3339 timestamp: "+err.Error())
			return
		}
		params.ScheduledEndTime = &endTime
	}

	if !plan.ChannelID.IsNull() && !plan.ChannelID.IsUnknown() {
		s := discord.Snowflake(plan.ChannelID.ValueString())
		params.ChannelID = &s
	}

	if !plan.EntityMetadataLocation.IsNull() && !plan.EntityMetadataLocation.IsUnknown() {
		loc := plan.EntityMetadataLocation.ValueString()
		params.EntityMetadata = &discord.ScheduledEventEntityMetadata{
			Location: &loc,
		}
	}

	if !plan.Image.IsNull() && !plan.Image.IsUnknown() {
		v := plan.Image.ValueString()
		params.Image = &v
	}

	event, err := r.client.ModifyGuildScheduledEvent(
		ctx,
		discord.Snowflake(plan.GuildID.ValueString()),
		discord.Snowflake(plan.ID.ValueString()),
		params,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating guild scheduled event", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenScheduledEvent(event, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the guild scheduled event.
func (r *guildScheduledEventResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state guildScheduledEventModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGuildScheduledEvent(
		ctx,
		discord.Snowflake(state.GuildID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
	)
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting guild scheduled event", err.Error())
	}
}

// ImportState supports importing by guild_id/event_id.
func (r *guildScheduledEventResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: guild_id/event_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// --- Helper functions ---

func flattenScheduledEvent(event *discord.GuildScheduledEvent, model *guildScheduledEventModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.ID = types.StringValue(event.ID.String())
	model.GuildID = types.StringValue(event.GuildID.String())
	model.Name = types.StringValue(event.Name)
	model.ScheduledStartTime = types.StringValue(event.ScheduledStartTime.Format(time.RFC3339))
	model.EntityType = types.Int64Value(int64(event.EntityType))
	model.PrivacyLevel = types.Int64Value(int64(event.PrivacyLevel))
	model.Status = types.Int64Value(int64(event.Status))

	if event.Description != nil {
		model.Description = types.StringValue(*event.Description)
	} else {
		model.Description = types.StringNull()
	}

	if event.ScheduledEndTime != nil {
		model.ScheduledEndTime = types.StringValue(event.ScheduledEndTime.Format(time.RFC3339))
	} else {
		model.ScheduledEndTime = types.StringNull()
	}

	if event.ChannelID != nil {
		model.ChannelID = types.StringValue(event.ChannelID.String())
	} else {
		model.ChannelID = types.StringNull()
	}

	if event.EntityMetadata != nil && event.EntityMetadata.Location != nil {
		model.EntityMetadataLocation = types.StringValue(*event.EntityMetadata.Location)
	} else {
		model.EntityMetadataLocation = types.StringNull()
	}

	if event.Image != nil {
		model.Image = types.StringValue(*event.Image)
	} else {
		model.Image = types.StringNull()
	}

	return diags
}
