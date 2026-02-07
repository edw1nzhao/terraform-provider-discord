package application_command

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
	_ resource.Resource                = &guildApplicationCommandResource{}
	_ resource.ResourceWithConfigure   = &guildApplicationCommandResource{}
	_ resource.ResourceWithImportState = &guildApplicationCommandResource{}
)

// NewGuildApplicationCommandResource is a constructor that returns a new guild application command resource.
func NewGuildApplicationCommandResource() resource.Resource {
	return &guildApplicationCommandResource{}
}

// guildApplicationCommandResource is the resource implementation.
type guildApplicationCommandResource struct {
	client        *discord.Client
	applicationID string
}

// guildApplicationCommandResourceModel maps the resource schema data.
type guildApplicationCommandResourceModel struct {
	ID                       types.String `tfsdk:"id"`
	ApplicationID            types.String `tfsdk:"application_id"`
	GuildID                  types.String `tfsdk:"guild_id"`
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	Type                     types.Int64  `tfsdk:"type"`
	DefaultMemberPermissions types.String `tfsdk:"default_member_permissions"`
	NSFW                     types.Bool   `tfsdk:"nsfw"`
	Options                  types.String `tfsdk:"options"`
}

// Metadata returns the resource type name.
func (r *guildApplicationCommandResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guild_application_command"
}

// Schema defines the schema for the resource.
func (r *guildApplicationCommandResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord guild-scoped application command.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the command.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_id": schema.StringAttribute{
				Description: "The ID of the application. Automatically set from provider configuration.",
				Computed:    true,
			},
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild this command is scoped to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the command (1-32 characters).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the command (1-100 characters). Required for CHAT_INPUT commands.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"type": schema.Int64Attribute{
				Description: "The type of command (1=CHAT_INPUT, 2=USER, 3=MESSAGE). Default: 1.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"default_member_permissions": schema.StringAttribute{
				Description: "Default member permissions required to use the command (permission bitfield string).",
				Optional:    true,
			},
			"nsfw": schema.BoolAttribute{
				Description: "Whether the command is age-restricted.",
				Optional:    true,
			},
			"options": schema.StringAttribute{
				Description: "JSON-encoded array of command options. Use JSON for complex nested option structures.",
				Optional:    true,
			},
		},
	}
}

// Configure sets the provider data on the resource.
func (r *guildApplicationCommandResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	data := common.ProviderDataFromConfig(req.ProviderData, &resp.Diagnostics)
	if data != nil {
		r.client = data.Client
		r.applicationID = data.ApplicationID
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *guildApplicationCommandResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan guildApplicationCommandResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.applicationID == "" {
		resp.Diagnostics.AddError("Missing Application ID", "The application_id must be set in the provider configuration to manage application commands.")
		return
	}

	params := &discord.CreateCommandParams{
		Name: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		v := plan.Description.ValueString()
		params.Description = &v
	}

	cmdType := int(plan.Type.ValueInt64())
	params.Type = &cmdType

	if !plan.DefaultMemberPermissions.IsNull() && !plan.DefaultMemberPermissions.IsUnknown() {
		v := plan.DefaultMemberPermissions.ValueString()
		params.DefaultMemberPermissions = &v
	}

	if !plan.NSFW.IsNull() && !plan.NSFW.IsUnknown() {
		v := plan.NSFW.ValueBool()
		params.NSFW = &v
	}

	if !plan.Options.IsNull() && !plan.Options.IsUnknown() {
		var options []*discord.ApplicationCommandOption
		if err := json.Unmarshal([]byte(plan.Options.ValueString()), &options); err != nil {
			resp.Diagnostics.AddError("Invalid options JSON", fmt.Sprintf("Failed to parse options: %s", err.Error()))
			return
		}
		params.Options = options
	}

	cmd, err := r.client.CreateGuildApplicationCommand(
		ctx,
		discord.Snowflake(r.applicationID),
		discord.Snowflake(plan.GuildID.ValueString()),
		params,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating guild application command", err.Error())
		return
	}

	r.flattenCommand(cmd, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *guildApplicationCommandResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state guildApplicationCommandResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appID := r.applicationID
	if appID == "" {
		appID = state.ApplicationID.ValueString()
	}

	cmd, err := r.client.GetGuildApplicationCommand(
		ctx,
		discord.Snowflake(appID),
		discord.Snowflake(state.GuildID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
	)
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading guild application command", err.Error())
		return
	}

	r.flattenCommand(cmd, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *guildApplicationCommandResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan guildApplicationCommandResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state guildApplicationCommandResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	params := &discord.EditCommandParams{
		Name: &name,
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		v := plan.Description.ValueString()
		params.Description = &v
	}

	if !plan.DefaultMemberPermissions.IsNull() && !plan.DefaultMemberPermissions.IsUnknown() {
		v := plan.DefaultMemberPermissions.ValueString()
		params.DefaultMemberPermissions = &v
	}

	if !plan.NSFW.IsNull() && !plan.NSFW.IsUnknown() {
		v := plan.NSFW.ValueBool()
		params.NSFW = &v
	}

	if !plan.Options.IsNull() && !plan.Options.IsUnknown() {
		var options []*discord.ApplicationCommandOption
		if err := json.Unmarshal([]byte(plan.Options.ValueString()), &options); err != nil {
			resp.Diagnostics.AddError("Invalid options JSON", fmt.Sprintf("Failed to parse options: %s", err.Error()))
			return
		}
		params.Options = options
	}

	appID := r.applicationID
	if appID == "" {
		appID = state.ApplicationID.ValueString()
	}

	cmd, err := r.client.EditGuildApplicationCommand(
		ctx,
		discord.Snowflake(appID),
		discord.Snowflake(state.GuildID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
		params,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating guild application command", err.Error())
		return
	}

	r.flattenCommand(cmd, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state.
func (r *guildApplicationCommandResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state guildApplicationCommandResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appID := r.applicationID
	if appID == "" {
		appID = state.ApplicationID.ValueString()
	}

	err := r.client.DeleteGuildApplicationCommand(
		ctx,
		discord.Snowflake(appID),
		discord.Snowflake(state.GuildID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
	)
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting guild application command", err.Error())
	}
}

// ImportState imports the resource state from application_id/guild_id/command_id.
func (r *guildApplicationCommandResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format application_id/guild_id/command_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
}

// flattenCommand maps the API response to the Terraform state model.
func (r *guildApplicationCommandResource) flattenCommand(cmd *discord.ApplicationCommand, model *guildApplicationCommandResourceModel) {
	model.ID = types.StringValue(cmd.ID.String())
	model.ApplicationID = types.StringValue(cmd.ApplicationID.String())
	model.Name = types.StringValue(cmd.Name)

	if cmd.GuildID != nil {
		model.GuildID = types.StringValue(cmd.GuildID.String())
	}

	if cmd.Description != "" {
		model.Description = types.StringValue(cmd.Description)
	} else {
		model.Description = types.StringNull()
	}

	if cmd.Type != nil {
		model.Type = types.Int64Value(int64(*cmd.Type))
	} else {
		model.Type = types.Int64Value(1)
	}

	if cmd.DefaultMemberPermissions != nil {
		model.DefaultMemberPermissions = types.StringValue(*cmd.DefaultMemberPermissions)
	} else {
		model.DefaultMemberPermissions = types.StringNull()
	}

	if cmd.NSFW != nil {
		model.NSFW = types.BoolValue(*cmd.NSFW)
	} else {
		model.NSFW = types.BoolNull()
	}

	if len(cmd.Options) > 0 {
		optionsJSON, err := json.Marshal(cmd.Options)
		if err == nil {
			model.Options = types.StringValue(string(optionsJSON))
		}
	} else {
		model.Options = types.StringNull()
	}
}
