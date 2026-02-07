package emoji

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
	_ resource.Resource                = &guildEmojiResource{}
	_ resource.ResourceWithConfigure   = &guildEmojiResource{}
	_ resource.ResourceWithImportState = &guildEmojiResource{}
)

// NewGuildEmojiResource is a constructor that returns a new guild emoji resource.
func NewGuildEmojiResource() resource.Resource {
	return &guildEmojiResource{}
}

// guildEmojiResource is the resource implementation.
type guildEmojiResource struct {
	client *discord.Client
}

// guildEmojiResourceModel maps the resource schema data.
type guildEmojiResourceModel struct {
	ID        types.String `tfsdk:"id"`
	GuildID   types.String `tfsdk:"guild_id"`
	Name      types.String `tfsdk:"name"`
	Image     types.String `tfsdk:"image"`
	Roles     types.Set    `tfsdk:"roles"`
	Animated  types.Bool   `tfsdk:"animated"`
	Available types.Bool   `tfsdk:"available"`
}

// Metadata returns the resource type name.
func (r *guildEmojiResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guild_emoji"
}

// Schema defines the schema for the resource.
func (r *guildEmojiResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord guild emoji.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the emoji.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild this emoji belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the emoji.",
				Required:    true,
			},
			"image": schema.StringAttribute{
				Description: "The base64 encoded image for the emoji (data URI scheme). Only used on create.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"roles": schema.SetAttribute{
				Description: "Set of role IDs allowed to use this emoji.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"animated": schema.BoolAttribute{
				Description: "Whether the emoji is animated.",
				Computed:    true,
			},
			"available": schema.BoolAttribute{
				Description: "Whether the emoji is available for use.",
				Computed:    true,
			},
		},
	}
}

// Configure sets the provider data on the resource.
func (r *guildEmojiResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource and sets the initial Terraform state.
func (r *guildEmojiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan guildEmojiResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.CreateEmojiParams{
		Name:  plan.Name.ValueString(),
		Image: plan.Image.ValueString(),
	}

	if !plan.Roles.IsNull() && !plan.Roles.IsUnknown() {
		var roleIDs []string
		resp.Diagnostics.Append(plan.Roles.ElementsAs(ctx, &roleIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, id := range roleIDs {
			params.Roles = append(params.Roles, discord.Snowflake(id))
		}
	}

	emoji, err := r.client.CreateGuildEmoji(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating guild emoji", err.Error())
		return
	}

	resp.Diagnostics.Append(r.flattenEmoji(ctx, emoji, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *guildEmojiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state guildEmojiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	emoji, err := r.client.GetGuildEmoji(ctx, discord.Snowflake(state.GuildID.ValueString()), discord.Snowflake(state.ID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading guild emoji", err.Error())
		return
	}

	resp.Diagnostics.Append(r.flattenEmoji(ctx, emoji, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *guildEmojiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan guildEmojiResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state guildEmojiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	params := &discord.ModifyEmojiParams{
		Name: &name,
	}

	if !plan.Roles.IsNull() && !plan.Roles.IsUnknown() {
		var roleIDs []string
		resp.Diagnostics.Append(plan.Roles.ElementsAs(ctx, &roleIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		roles := make([]discord.Snowflake, len(roleIDs))
		for i, id := range roleIDs {
			roles[i] = discord.Snowflake(id)
		}
		params.Roles = roles
	} else {
		// Send empty array to clear roles.
		params.Roles = []discord.Snowflake{}
	}

	emoji, err := r.client.ModifyGuildEmoji(ctx, discord.Snowflake(state.GuildID.ValueString()), discord.Snowflake(state.ID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating guild emoji", err.Error())
		return
	}

	resp.Diagnostics.Append(r.flattenEmoji(ctx, emoji, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state.
func (r *guildEmojiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state guildEmojiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGuildEmoji(ctx, discord.Snowflake(state.GuildID.ValueString()), discord.Snowflake(state.ID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting guild emoji", err.Error())
	}
}

// ImportState imports the resource state from guild_id/emoji_id.
func (r *guildEmojiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format guild_id/emoji_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// flattenEmoji maps the API response to the Terraform state model.
func (r *guildEmojiResource) flattenEmoji(ctx context.Context, emoji *discord.Emoji, model *guildEmojiResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if emoji.ID != nil {
		model.ID = types.StringValue(emoji.ID.String())
	}
	if emoji.Name != nil {
		model.Name = types.StringValue(*emoji.Name)
	}
	model.Animated = types.BoolValue(emoji.Animated)
	model.Available = types.BoolValue(emoji.Available)

	if len(emoji.Roles) > 0 {
		roleIDs := make([]string, len(emoji.Roles))
		for i, role := range emoji.Roles {
			roleIDs[i] = role.String()
		}
		rolesSet, d := types.SetValueFrom(ctx, types.StringType, roleIDs)
		diags.Append(d...)
		model.Roles = rolesSet
	} else {
		model.Roles = types.SetNull(types.StringType)
	}

	return diags
}
