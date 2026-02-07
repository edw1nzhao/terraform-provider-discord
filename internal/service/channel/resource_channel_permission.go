package channel

import (
	"context"
	"fmt"
	"strings"

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
	_ resource.Resource                = &channelPermissionResource{}
	_ resource.ResourceWithConfigure   = &channelPermissionResource{}
	_ resource.ResourceWithImportState = &channelPermissionResource{}
)

// channelPermissionResource is the resource implementation.
type channelPermissionResource struct {
	client *discord.Client
}

// channelPermissionResourceModel maps the resource schema to a Go struct.
type channelPermissionResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ChannelID   types.String `tfsdk:"channel_id"`
	OverwriteID types.String `tfsdk:"overwrite_id"`
	Type        types.Int64  `tfsdk:"type"`
	Allow       types.String `tfsdk:"allow"`
	Deny        types.String `tfsdk:"deny"`
}

// NewChannelPermissionResource returns a new channel permission resource.
func NewChannelPermissionResource() resource.Resource {
	return &channelPermissionResource{}
}

// Metadata returns the resource type name.
func (r *channelPermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel_permission"
}

// Configure adds the provider configured client to the resource.
func (r *channelPermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the resource.
func (r *channelPermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord channel permission overwrite for a role or member.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The composite ID of the permission overwrite (channel_id/overwrite_id).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"channel_id": schema.StringAttribute{
				Description: "The ID of the channel.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"overwrite_id": schema.StringAttribute{
				Description: "The ID of the role or user for the permission overwrite.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.Int64Attribute{
				Description: "The type of the permission overwrite (0=role, 1=member).",
				Required:    true,
			},
			"allow": schema.StringAttribute{
				Description: "The bitwise value of all allowed permissions.",
				Optional:    true,
				Computed:    true,
			},
			"deny": schema.StringAttribute{
				Description: "The bitwise value of all denied permissions.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// Create creates the channel permission overwrite.
func (r *channelPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan channelPermissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.EditPermissionsParams{
		Type: int(plan.Type.ValueInt64()),
	}
	if !plan.Allow.IsNull() && !plan.Allow.IsUnknown() {
		v := plan.Allow.ValueString()
		params.Allow = &v
	}
	if !plan.Deny.IsNull() && !plan.Deny.IsUnknown() {
		v := plan.Deny.ValueString()
		params.Deny = &v
	}

	channelID := discord.Snowflake(plan.ChannelID.ValueString())
	overwriteID := discord.Snowflake(plan.OverwriteID.ValueString())

	err := r.client.EditChannelPermissions(ctx, channelID, overwriteID, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Discord Channel Permission",
			"Could not create channel permission overwrite: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(plan.ChannelID.ValueString() + "/" + plan.OverwriteID.ValueString())

	// Set defaults for computed fields if not provided.
	if plan.Allow.IsNull() || plan.Allow.IsUnknown() {
		plan.Allow = types.StringValue("0")
	}
	if plan.Deny.IsNull() || plan.Deny.IsUnknown() {
		plan.Deny = types.StringValue("0")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *channelPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state channelPermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ch, err := r.client.GetChannel(ctx, discord.Snowflake(state.ChannelID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Discord Channel Permission",
			"Could not read channel ID "+state.ChannelID.ValueString()+": "+err.Error(),
		)
		return
	}

	overwriteID := state.OverwriteID.ValueString()
	var found *discord.PermissionOverwrite
	for _, ow := range ch.PermissionOverwrites {
		if ow.ID.String() == overwriteID {
			found = ow
			break
		}
	}

	if found == nil {
		// The permission overwrite was deleted externally.
		resp.State.RemoveResource(ctx)
		return
	}

	state.Type = types.Int64Value(int64(found.Type))
	state.Allow = types.StringValue(found.Allow)
	state.Deny = types.StringValue(found.Deny)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update modifies the channel permission overwrite.
func (r *channelPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan channelPermissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state channelPermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.EditPermissionsParams{
		Type: int(plan.Type.ValueInt64()),
	}
	if !plan.Allow.IsNull() && !plan.Allow.IsUnknown() {
		v := plan.Allow.ValueString()
		params.Allow = &v
	}
	if !plan.Deny.IsNull() && !plan.Deny.IsUnknown() {
		v := plan.Deny.ValueString()
		params.Deny = &v
	}

	channelID := discord.Snowflake(state.ChannelID.ValueString())
	overwriteID := discord.Snowflake(state.OverwriteID.ValueString())

	err := r.client.EditChannelPermissions(ctx, channelID, overwriteID, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Discord Channel Permission",
			"Could not update channel permission overwrite: "+err.Error(),
		)
		return
	}

	plan.ID = state.ID

	if plan.Allow.IsNull() || plan.Allow.IsUnknown() {
		plan.Allow = types.StringValue("0")
	}
	if plan.Deny.IsNull() || plan.Deny.IsUnknown() {
		plan.Deny = types.StringValue("0")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the channel permission overwrite.
func (r *channelPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state channelPermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channelID := discord.Snowflake(state.ChannelID.ValueString())
	overwriteID := discord.Snowflake(state.OverwriteID.ValueString())

	err := r.client.DeleteChannelPermission(ctx, channelID, overwriteID)
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Discord Channel Permission",
			"Could not delete channel permission overwrite: "+err.Error(),
		)
	}
}

// ImportState allows importing an existing channel permission by channel_id/overwrite_id.
func (r *channelPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format 'channel_id/overwrite_id', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("channel_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("overwrite_id"), parts[1])...)
}
