package widget

import (
	"context"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &guildWidgetResource{}
	_ resource.ResourceWithConfigure   = &guildWidgetResource{}
	_ resource.ResourceWithImportState = &guildWidgetResource{}
)

// guildWidgetResource is the resource implementation.
type guildWidgetResource struct {
	client *discord.Client
}

// guildWidgetModel maps the resource schema data.
type guildWidgetModel struct {
	GuildID   types.String `tfsdk:"guild_id"`
	Enabled   types.Bool   `tfsdk:"enabled"`
	ChannelID types.String `tfsdk:"channel_id"`
}

// NewGuildWidgetResource is a helper function to simplify the provider implementation.
func NewGuildWidgetResource() resource.Resource {
	return &guildWidgetResource{}
}

// Metadata returns the resource type name.
func (r *guildWidgetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guild_widget"
}

// Schema defines the schema for the resource.
func (r *guildWidgetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord guild widget settings. " +
			"Creating this resource applies the widget settings. " +
			"Deleting this resource disables the widget.",
		Attributes: map[string]schema.Attribute{
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the widget is enabled.",
				Required:    true,
			},
			"channel_id": schema.StringAttribute{
				Description: "The widget channel ID. Set to the channel that the widget will generate an invite to.",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *guildWidgetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource by performing an update (widget settings always exist).
func (r *guildWidgetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan guildWidgetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.ModifyWidgetParams{}
	enabled := plan.Enabled.ValueBool()
	params.Enabled = &enabled

	if !plan.ChannelID.IsNull() && !plan.ChannelID.IsUnknown() {
		cid := discord.Snowflake(plan.ChannelID.ValueString())
		params.ChannelID = &cid
	}

	widget, err := r.client.ModifyGuildWidget(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Guild Widget",
			"Could not update guild widget settings: "+err.Error(),
		)
		return
	}

	plan.Enabled = types.BoolValue(widget.Enabled)
	if widget.ChannelID != nil {
		plan.ChannelID = types.StringValue(widget.ChannelID.String())
	} else {
		plan.ChannelID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *guildWidgetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state guildWidgetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	widget, err := r.client.GetGuildWidgetSettings(ctx, discord.Snowflake(state.GuildID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Guild Widget",
			"Could not read guild widget settings: "+err.Error(),
		)
		return
	}

	state.Enabled = types.BoolValue(widget.Enabled)
	if widget.ChannelID != nil {
		state.ChannelID = types.StringValue(widget.ChannelID.String())
	} else {
		state.ChannelID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *guildWidgetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan guildWidgetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.ModifyWidgetParams{}
	enabled := plan.Enabled.ValueBool()
	params.Enabled = &enabled

	if !plan.ChannelID.IsNull() && !plan.ChannelID.IsUnknown() {
		cid := discord.Snowflake(plan.ChannelID.ValueString())
		params.ChannelID = &cid
	}

	widget, err := r.client.ModifyGuildWidget(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Guild Widget",
			"Could not update guild widget settings: "+err.Error(),
		)
		return
	}

	plan.Enabled = types.BoolValue(widget.Enabled)
	if widget.ChannelID != nil {
		plan.ChannelID = types.StringValue(widget.ChannelID.String())
	} else {
		plan.ChannelID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete disables the widget instead of actually deleting it.
func (r *guildWidgetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state guildWidgetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disabled := false
	params := &discord.ModifyWidgetParams{
		Enabled: &disabled,
	}

	_, err := r.client.ModifyGuildWidget(ctx, discord.Snowflake(state.GuildID.ValueString()), params)
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Guild Widget",
			"Could not disable guild widget: "+err.Error(),
		)
	}
}

// ImportState implements the import by guild_id.
func (r *guildWidgetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), req.ID)...)
}
