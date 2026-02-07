package sticker

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
	_ resource.Resource                = &guildStickerResource{}
	_ resource.ResourceWithConfigure   = &guildStickerResource{}
	_ resource.ResourceWithImportState = &guildStickerResource{}
)

// NewGuildStickerResource is a constructor that returns a new guild sticker resource.
func NewGuildStickerResource() resource.Resource {
	return &guildStickerResource{}
}

// guildStickerResource is the resource implementation.
type guildStickerResource struct {
	client *discord.Client
}

// guildStickerResourceModel maps the resource schema data.
type guildStickerResourceModel struct {
	ID          types.String `tfsdk:"id"`
	GuildID     types.String `tfsdk:"guild_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Tags        types.String `tfsdk:"tags"`
	FormatType  types.Int64  `tfsdk:"format_type"`
	Available   types.Bool   `tfsdk:"available"`
}

// Metadata returns the resource type name.
func (r *guildStickerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guild_sticker"
}

// Schema defines the schema for the resource.
func (r *guildStickerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord guild sticker. Note: Stickers must be created in Discord first " +
			"(file upload is not supported). Use `terraform import` to bring an existing sticker under management. " +
			"Create will attempt to use the JSON API but may fail without multipart form data for the file.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the sticker.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild this sticker belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the sticker (2-30 characters).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 30),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the sticker (2-100 characters).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 100),
				},
			},
			"tags": schema.StringAttribute{
				Description: "Autocomplete/suggestion tags for the sticker (max 200 characters).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"format_type": schema.Int64Attribute{
				Description: "The format type of the sticker (1=PNG, 2=APNG, 3=LOTTIE, 4=GIF).",
				Computed:    true,
			},
			"available": schema.BoolAttribute{
				Description: "Whether the sticker is available for use.",
				Computed:    true,
			},
		},
	}
}

// Configure sets the provider data on the resource.
func (r *guildStickerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource and sets the initial Terraform state.
func (r *guildStickerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan guildStickerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	description := plan.Description.ValueString()
	params := &discord.CreateStickerParams{
		Name:        plan.Name.ValueString(),
		Description: &description,
		Tags:        plan.Tags.ValueString(),
	}

	sticker, err := r.client.CreateGuildSticker(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating guild sticker",
			"Sticker creation requires multipart form data for the file upload. "+
				"Consider creating the sticker in Discord first and importing it. Error: "+err.Error(),
		)
		return
	}

	r.flattenSticker(sticker, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *guildStickerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state guildStickerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sticker, err := r.client.GetGuildSticker(ctx, discord.Snowflake(state.GuildID.ValueString()), discord.Snowflake(state.ID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading guild sticker", err.Error())
		return
	}

	r.flattenSticker(sticker, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *guildStickerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan guildStickerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state guildStickerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	description := plan.Description.ValueString()
	tags := plan.Tags.ValueString()

	params := &discord.ModifyStickerParams{
		Name:        &name,
		Description: &description,
		Tags:        &tags,
	}

	sticker, err := r.client.ModifyGuildSticker(ctx, discord.Snowflake(state.GuildID.ValueString()), discord.Snowflake(state.ID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating guild sticker", err.Error())
		return
	}

	r.flattenSticker(sticker, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state.
func (r *guildStickerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state guildStickerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGuildSticker(ctx, discord.Snowflake(state.GuildID.ValueString()), discord.Snowflake(state.ID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting guild sticker", err.Error())
	}
}

// ImportState imports the resource state from guild_id/sticker_id.
func (r *guildStickerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format guild_id/sticker_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// flattenSticker maps the API response to the Terraform state model.
func (r *guildStickerResource) flattenSticker(sticker *discord.Sticker, model *guildStickerResourceModel) {
	model.ID = types.StringValue(sticker.ID.String())
	model.Name = types.StringValue(sticker.Name)
	model.Tags = types.StringValue(sticker.Tags)
	model.FormatType = types.Int64Value(int64(sticker.FormatType))
	model.Available = types.BoolValue(sticker.Available)

	if sticker.Description != nil {
		model.Description = types.StringValue(*sticker.Description)
	} else {
		model.Description = types.StringNull()
	}

	if sticker.GuildID != nil {
		model.GuildID = types.StringValue(sticker.GuildID.String())
	}
}
