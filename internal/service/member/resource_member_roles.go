package member

import (
	"context"
	"fmt"
	"strings"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &memberRolesResource{}
	_ resource.ResourceWithConfigure   = &memberRolesResource{}
	_ resource.ResourceWithImportState = &memberRolesResource{}
)

// memberRolesResource is the resource implementation.
type memberRolesResource struct {
	client *discord.Client
}

// memberRolesResourceModel maps the resource schema to a Go struct.
type memberRolesResourceModel struct {
	ID      types.String `tfsdk:"id"`
	GuildID types.String `tfsdk:"guild_id"`
	UserID  types.String `tfsdk:"user_id"`
	Roles   types.Set    `tfsdk:"roles"`
}

// NewMemberRolesResource returns a new member roles resource.
func NewMemberRolesResource() resource.Resource {
	return &memberRolesResource{}
}

// Metadata returns the resource type name.
func (r *memberRolesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_member_roles"
}

// Configure adds the provider configured client to the resource.
func (r *memberRolesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the resource.
func (r *memberRolesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the set of roles for a Discord guild member.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The composite ID (guild_id/user_id).",
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
			"user_id": schema.StringAttribute{
				Description: "The ID of the user.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"roles": schema.SetAttribute{
				Description: "The set of role IDs assigned to the member.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Create sets the member roles.
func (r *memberRolesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan memberRolesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guildID := discord.Snowflake(plan.GuildID.ValueString())
	userID := discord.Snowflake(plan.UserID.ValueString())

	roleIDs, diags := extractRoleIDs(ctx, plan.Roles)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set all roles at once using ModifyGuildMember.
	snowflakes := make([]discord.Snowflake, len(roleIDs))
	for i, id := range roleIDs {
		snowflakes[i] = discord.Snowflake(id)
	}

	_, err := r.client.ModifyGuildMember(ctx, guildID, userID, &discord.ModifyMemberParams{
		Roles: snowflakes,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Setting Discord Member Roles",
			"Could not set member roles: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(plan.GuildID.ValueString() + "/" + plan.UserID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *memberRolesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state memberRolesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guildID := discord.Snowflake(state.GuildID.ValueString())
	userID := discord.Snowflake(state.UserID.ValueString())

	member, err := r.client.GetGuildMember(ctx, guildID, userID)
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Discord Member Roles",
			"Could not read member "+state.UserID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Convert the member roles to a set.
	roleStrings := make([]string, len(member.Roles))
	for i, r := range member.Roles {
		roleStrings[i] = r.String()
	}

	rolesSet, diags := types.SetValueFrom(ctx, types.StringType, roleStrings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Roles = rolesSet
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update modifies the member roles.
func (r *memberRolesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan memberRolesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state memberRolesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guildID := discord.Snowflake(state.GuildID.ValueString())
	userID := discord.Snowflake(state.UserID.ValueString())

	roleIDs, diags := extractRoleIDs(ctx, plan.Roles)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	snowflakes := make([]discord.Snowflake, len(roleIDs))
	for i, id := range roleIDs {
		snowflakes[i] = discord.Snowflake(id)
	}

	_, err := r.client.ModifyGuildMember(ctx, guildID, userID, &discord.ModifyMemberParams{
		Roles: snowflakes,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Discord Member Roles",
			"Could not update member roles: "+err.Error(),
		)
		return
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete removes all managed roles from the member. This sets the member's
// roles to an empty list, effectively removing all roles.
func (r *memberRolesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state memberRolesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guildID := discord.Snowflake(state.GuildID.ValueString())
	userID := discord.Snowflake(state.UserID.ValueString())

	_, err := r.client.ModifyGuildMember(ctx, guildID, userID, &discord.ModifyMemberParams{
		Roles: []discord.Snowflake{},
	})
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Removing Discord Member Roles",
			"Could not remove member roles: "+err.Error(),
		)
	}
}

// ImportState allows importing by guild_id/user_id.
func (r *memberRolesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format 'guild_id/user_id', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), parts[1])...)
}

// extractRoleIDs extracts role ID strings from a types.Set.
func extractRoleIDs(ctx context.Context, rolesSet types.Set) ([]string, diag.Diagnostics) {
	var roleIDs []string
	diags := rolesSet.ElementsAs(ctx, &roleIDs, false)
	return roleIDs, diags
}
