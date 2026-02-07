package stage

import (
	"context"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &stageInstanceResource{}
	_ resource.ResourceWithConfigure   = &stageInstanceResource{}
	_ resource.ResourceWithImportState = &stageInstanceResource{}
)

// NewStageInstanceResource is a constructor that returns a new stage instance resource.
func NewStageInstanceResource() resource.Resource {
	return &stageInstanceResource{}
}

// stageInstanceResource is the resource implementation.
type stageInstanceResource struct {
	client *discord.Client
}

// stageInstanceResourceModel maps the resource schema data.
type stageInstanceResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	ChannelID             types.String `tfsdk:"channel_id"`
	Topic                 types.String `tfsdk:"topic"`
	PrivacyLevel          types.Int64  `tfsdk:"privacy_level"`
	GuildID               types.String `tfsdk:"guild_id"`
	GuildScheduledEventID types.String `tfsdk:"guild_scheduled_event_id"`
}

// Metadata returns the resource type name.
func (r *stageInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_instance"
}

// Schema defines the schema for the resource.
func (r *stageInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord stage instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the stage instance.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"channel_id": schema.StringAttribute{
				Description: "The ID of the stage channel.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"topic": schema.StringAttribute{
				Description: "The topic of the stage instance (1-120 characters).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 120),
				},
			},
			"privacy_level": schema.Int64Attribute{
				Description: "The privacy level of the stage instance (2 = GUILD_ONLY). Default: 2.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(2),
			},
			"guild_id": schema.StringAttribute{
				Description: "The guild ID of the associated stage channel.",
				Computed:    true,
			},
			"guild_scheduled_event_id": schema.StringAttribute{
				Description: "The ID of the guild scheduled event associated with this stage instance.",
				Optional:    true,
			},
		},
	}
}

// Configure sets the provider data on the resource.
func (r *stageInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource and sets the initial Terraform state.
func (r *stageInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan stageInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	privacyLevel := int(plan.PrivacyLevel.ValueInt64())
	params := &discord.CreateStageInstanceParams{
		ChannelID:    discord.Snowflake(plan.ChannelID.ValueString()),
		Topic:        plan.Topic.ValueString(),
		PrivacyLevel: &privacyLevel,
	}

	if !plan.GuildScheduledEventID.IsNull() && !plan.GuildScheduledEventID.IsUnknown() {
		eventID := discord.Snowflake(plan.GuildScheduledEventID.ValueString())
		params.GuildScheduledEventID = &eventID
	}

	stage, err := r.client.CreateStageInstance(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating stage instance", err.Error())
		return
	}

	r.flattenStageInstance(stage, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *stageInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state stageInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Stage instances are fetched by channel ID.
	stage, err := r.client.GetStageInstance(ctx, discord.Snowflake(state.ChannelID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading stage instance", err.Error())
		return
	}

	r.flattenStageInstance(stage, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *stageInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan stageInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state stageInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	topic := plan.Topic.ValueString()
	privacyLevel := int(plan.PrivacyLevel.ValueInt64())
	params := &discord.ModifyStageInstanceParams{
		Topic:        &topic,
		PrivacyLevel: &privacyLevel,
	}

	stage, err := r.client.ModifyStageInstance(ctx, discord.Snowflake(state.ChannelID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating stage instance", err.Error())
		return
	}

	r.flattenStageInstance(stage, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state.
func (r *stageInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state stageInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteStageInstance(ctx, discord.Snowflake(state.ChannelID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting stage instance", err.Error())
	}
}

// ImportState imports the resource state using the channel ID.
func (r *stageInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Stage instances are keyed by channel ID.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("channel_id"), req.ID)...)
}

// flattenStageInstance maps the API response to the Terraform state model.
func (r *stageInstanceResource) flattenStageInstance(stage *discord.StageInstance, model *stageInstanceResourceModel) {
	model.ID = types.StringValue(stage.ID.String())
	model.ChannelID = types.StringValue(stage.ChannelID.String())
	model.Topic = types.StringValue(stage.Topic)
	model.PrivacyLevel = types.Int64Value(int64(stage.PrivacyLevel))
	model.GuildID = types.StringValue(stage.GuildID.String())

	if stage.GuildScheduledEventID != nil {
		model.GuildScheduledEventID = types.StringValue(stage.GuildScheduledEventID.String())
	} else {
		model.GuildScheduledEventID = types.StringNull()
	}
}
