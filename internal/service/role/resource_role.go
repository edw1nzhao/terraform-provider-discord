package role

import (
	"context"
	"fmt"
	"strings"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &roleResource{}
	_ resource.ResourceWithConfigure   = &roleResource{}
	_ resource.ResourceWithImportState = &roleResource{}
)

// roleResource is the resource implementation.
type roleResource struct {
	client *discord.Client
}

// roleResourceModel maps the resource schema to a Go struct.
type roleResourceModel struct {
	ID           types.String `tfsdk:"id"`
	GuildID      types.String `tfsdk:"guild_id"`
	Name         types.String `tfsdk:"name"`
	Permissions  types.String `tfsdk:"permissions"`
	Color        types.Int64  `tfsdk:"color"`
	Hoist        types.Bool   `tfsdk:"hoist"`
	Icon         types.String `tfsdk:"icon"`
	UnicodeEmoji types.String `tfsdk:"unicode_emoji"`
	Mentionable  types.Bool   `tfsdk:"mentionable"`
	Position     types.Int64  `tfsdk:"position"`
	Managed      types.Bool   `tfsdk:"managed"`
}

// NewRoleResource returns a new role resource.
func NewRoleResource() resource.Resource {
	return &roleResource{}
}

// Metadata returns the resource type name.
func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Configure adds the provider configured client to the resource.
func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the resource.
func (r *roleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord guild role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the role.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild this role belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the role.",
				Required:    true,
			},
			"permissions": schema.StringAttribute{
				Description: "The permission bitfield for the role.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"color": schema.Int64Attribute{
				Description: "The RGB color value for the role (integer).",
				Optional:    true,
				Computed:    true,
			},
			"hoist": schema.BoolAttribute{
				Description: "Whether the role should be displayed separately in the sidebar.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"icon": schema.StringAttribute{
				Description: "The role icon as a base64-encoded image data URI.",
				Optional:    true,
			},
			"unicode_emoji": schema.StringAttribute{
				Description: "The role unicode emoji.",
				Optional:    true,
			},
			"mentionable": schema.BoolAttribute{
				Description: "Whether the role can be mentioned by everyone.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"position": schema.Int64Attribute{
				Description: "The position of the role. Roles with the same position are sorted by ID.",
				Optional:    true,
				Computed:    true,
			},
			"managed": schema.BoolAttribute{
				Description: "Whether the role is managed by an integration.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the role resource.
func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan roleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	params := &discord.CreateRoleParams{
		Name: &name,
	}

	if !plan.Permissions.IsNull() && !plan.Permissions.IsUnknown() {
		v := plan.Permissions.ValueString()
		params.Permissions = &v
	}
	if !plan.Color.IsNull() && !plan.Color.IsUnknown() {
		v := int(plan.Color.ValueInt64())
		params.Color = &v
	}
	if !plan.Hoist.IsNull() && !plan.Hoist.IsUnknown() {
		v := plan.Hoist.ValueBool()
		params.Hoist = &v
	}
	if !plan.Icon.IsNull() && !plan.Icon.IsUnknown() {
		v := plan.Icon.ValueString()
		params.Icon = &v
	}
	if !plan.UnicodeEmoji.IsNull() && !plan.UnicodeEmoji.IsUnknown() {
		v := plan.UnicodeEmoji.ValueString()
		params.UnicodeEmoji = &v
	}
	if !plan.Mentionable.IsNull() && !plan.Mentionable.IsUnknown() {
		v := plan.Mentionable.ValueBool()
		params.Mentionable = &v
	}

	guildID := discord.Snowflake(plan.GuildID.ValueString())
	role, err := r.client.CreateGuildRole(ctx, guildID, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Discord Role",
			"Could not create role: "+err.Error(),
		)
		return
	}

	// If position is specified, move the role after creation.
	if !plan.Position.IsNull() && !plan.Position.IsUnknown() {
		pos := int(plan.Position.ValueInt64())
		_, err = r.client.ModifyGuildRolePositions(ctx, guildID, []*discord.RolePosition{
			{
				ID:       role.ID,
				Position: &pos,
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Setting Discord Role Position",
				"Role was created but position could not be set: "+err.Error(),
			)
			// Continue - the role was created successfully, just position failed.
		}
	}

	// Re-read to get the final state (position may have shifted).
	mapRoleToState(role, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state roleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guildID := discord.Snowflake(state.GuildID.ValueString())
	roleID := state.ID.ValueString()

	roles, err := r.client.GetGuildRoles(ctx, guildID)
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Discord Role",
			"Could not read roles for guild ID "+state.GuildID.ValueString()+": "+err.Error(),
		)
		return
	}

	var found *discord.Role
	for _, role := range roles {
		if role.ID.String() == roleID {
			found = role
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	mapRoleToState(found, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update modifies the role resource.
func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan roleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state roleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.ModifyRoleParams{}

	name := plan.Name.ValueString()
	params.Name = &name

	if !plan.Permissions.IsNull() && !plan.Permissions.IsUnknown() {
		v := plan.Permissions.ValueString()
		params.Permissions = &v
	}
	if !plan.Color.IsNull() && !plan.Color.IsUnknown() {
		v := int(plan.Color.ValueInt64())
		params.Color = &v
	}
	if !plan.Hoist.IsNull() && !plan.Hoist.IsUnknown() {
		v := plan.Hoist.ValueBool()
		params.Hoist = &v
	}
	if !plan.Icon.IsNull() {
		v := plan.Icon.ValueString()
		params.Icon = &v
	}
	if !plan.UnicodeEmoji.IsNull() {
		v := plan.UnicodeEmoji.ValueString()
		params.UnicodeEmoji = &v
	}
	if !plan.Mentionable.IsNull() && !plan.Mentionable.IsUnknown() {
		v := plan.Mentionable.ValueBool()
		params.Mentionable = &v
	}

	guildID := discord.Snowflake(state.GuildID.ValueString())
	roleID := discord.Snowflake(state.ID.ValueString())

	role, err := r.client.ModifyGuildRole(ctx, guildID, roleID, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Discord Role",
			"Could not update role ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update position if changed.
	if !plan.Position.IsNull() && !plan.Position.IsUnknown() {
		pos := int(plan.Position.ValueInt64())
		_, err = r.client.ModifyGuildRolePositions(ctx, guildID, []*discord.RolePosition{
			{
				ID:       roleID,
				Position: &pos,
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Discord Role Position",
				"Role was updated but position could not be set: "+err.Error(),
			)
		}
	}

	plan.ID = state.ID
	plan.GuildID = state.GuildID
	mapRoleToState(role, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the role resource.
func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guildID := discord.Snowflake(state.GuildID.ValueString())
	roleID := discord.Snowflake(state.ID.ValueString())

	err := r.client.DeleteGuildRole(ctx, guildID, roleID)
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Discord Role",
			"Could not delete role ID "+state.ID.ValueString()+": "+err.Error(),
		)
	}
}

// ImportState allows importing an existing role by guild_id/role_id.
func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format 'guild_id/role_id', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// mapRoleToState maps a Discord Role API response to the Terraform state model.
func mapRoleToState(role *discord.Role, state *roleResourceModel) {
	state.ID = types.StringValue(role.ID.String())
	state.Name = types.StringValue(role.Name)
	state.Permissions = types.StringValue(role.Permissions)
	state.Color = types.Int64Value(int64(role.Color))
	state.Hoist = types.BoolValue(role.Hoist)
	state.Mentionable = types.BoolValue(role.Mentionable)
	state.Position = types.Int64Value(int64(role.Position))
	state.Managed = types.BoolValue(role.Managed)

	if role.Icon != nil {
		state.Icon = types.StringValue(*role.Icon)
	} else {
		state.Icon = types.StringNull()
	}
	if role.UnicodeEmoji != nil {
		state.UnicodeEmoji = types.StringValue(*role.UnicodeEmoji)
	} else {
		state.UnicodeEmoji = types.StringNull()
	}
}
