package template

import (
	"context"
	"fmt"
	"strings"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &guildTemplateResource{}
	_ resource.ResourceWithConfigure   = &guildTemplateResource{}
	_ resource.ResourceWithImportState = &guildTemplateResource{}
)

// NewGuildTemplateResource is a constructor that returns a new guild template resource.
func NewGuildTemplateResource() resource.Resource {
	return &guildTemplateResource{}
}

// guildTemplateResource is the resource implementation.
type guildTemplateResource struct {
	client *discord.Client
}

// guildTemplateResourceModel maps the resource schema data.
type guildTemplateResourceModel struct {
	ID            types.String `tfsdk:"id"`
	GuildID       types.String `tfsdk:"guild_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	UsageCount    types.Int64  `tfsdk:"usage_count"`
	SourceGuildID types.String `tfsdk:"source_guild_id"`
	IsDirty       types.Bool   `tfsdk:"is_dirty"`
}

// Metadata returns the resource type name.
func (r *guildTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guild_template"
}

// Schema defines the schema for the resource.
func (r *guildTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord guild template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The template code (unique ID).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild this template is for.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the template (1-100 characters).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the template (0-120 characters).",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(120),
				},
			},
			"usage_count": schema.Int64Attribute{
				Description: "Number of times this template has been used.",
				Computed:    true,
			},
			"source_guild_id": schema.StringAttribute{
				Description: "The ID of the guild this template is based on.",
				Computed:    true,
			},
			"is_dirty": schema.BoolAttribute{
				Description: "Whether the template has unsynced changes.",
				Computed:    true,
			},
		},
	}
}

// Configure sets the provider data on the resource.
func (r *guildTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource and sets the initial Terraform state.
func (r *guildTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan guildTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.CreateTemplateParams{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		v := plan.Description.ValueString()
		params.Description = &v
	}

	tmpl, err := r.client.CreateGuildTemplate(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating guild template", err.Error())
		return
	}

	r.flattenTemplate(tmpl, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *guildTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state guildTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The Discord API does not have a "get single template by code" endpoint that
	// takes guild_id and code. We list all templates and find the matching one.
	templates, err := r.client.GetGuildTemplates(ctx, discord.Snowflake(state.GuildID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading guild templates", err.Error())
		return
	}

	var found *discord.GuildTemplate
	for _, t := range templates {
		if t.Code == state.ID.ValueString() {
			found = t
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	r.flattenTemplate(found, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *guildTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan guildTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state guildTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	params := &discord.ModifyTemplateParams{
		Name: &name,
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		v := plan.Description.ValueString()
		params.Description = &v
	}

	tmpl, err := r.client.ModifyGuildTemplate(ctx, discord.Snowflake(state.GuildID.ValueString()), state.ID.ValueString(), params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating guild template", err.Error())
		return
	}

	r.flattenTemplate(tmpl, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state.
func (r *guildTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state guildTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGuildTemplate(ctx, discord.Snowflake(state.GuildID.ValueString()), state.ID.ValueString())
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting guild template", err.Error())
	}
}

// ImportState imports the resource state from guild_id/template_code.
func (r *guildTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format guild_id/template_code, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// flattenTemplate maps the API response to the Terraform state model.
func (r *guildTemplateResource) flattenTemplate(tmpl *discord.GuildTemplate, model *guildTemplateResourceModel) {
	model.ID = types.StringValue(tmpl.Code)
	model.Name = types.StringValue(tmpl.Name)
	model.UsageCount = types.Int64Value(int64(tmpl.UsageCount))
	model.SourceGuildID = types.StringValue(tmpl.SourceGuildID.String())

	if tmpl.Description != nil {
		model.Description = types.StringValue(*tmpl.Description)
	} else {
		model.Description = types.StringNull()
	}

	if tmpl.IsDirty != nil {
		model.IsDirty = types.BoolValue(*tmpl.IsDirty)
	} else {
		model.IsDirty = types.BoolNull()
	}
}
