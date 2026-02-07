package ban

import (
	"context"
	"fmt"
	"strings"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure interface compliance.
var (
	_ resource.Resource                = &banResource{}
	_ resource.ResourceWithConfigure   = &banResource{}
	_ resource.ResourceWithImportState = &banResource{}
)

// NewBanResource returns a new resource for discord_ban.
func NewBanResource() resource.Resource {
	return &banResource{}
}

type banResource struct {
	client *discord.Client
}

// banModel maps the Terraform schema to Go types.
type banModel struct {
	ID                   types.String `tfsdk:"id"`
	GuildID              types.String `tfsdk:"guild_id"`
	UserID               types.String `tfsdk:"user_id"`
	Reason               types.String `tfsdk:"reason"`
	DeleteMessageSeconds types.Int64  `tfsdk:"delete_message_seconds"`
}

// Metadata sets the type name.
func (r *banResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ban"
}

// Schema defines the schema.
func (r *banResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord guild ban.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The composite ID of the ban (guild_id/user_id).",
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
				Description: "The ID of the user to ban.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"reason": schema.StringAttribute{
				Description: "The reason for the ban.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"delete_message_seconds": schema.Int64Attribute{
				Description: "Number of seconds to delete messages for, between 0 and 604800 (7 days).",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure stores the provider data.
func (r *banResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create creates a guild ban.
func (r *banResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan banModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guildID := discord.Snowflake(plan.GuildID.ValueString())
	userID := discord.Snowflake(plan.UserID.ValueString())

	params := &discord.CreateBanParams{}

	if !plan.DeleteMessageSeconds.IsNull() && !plan.DeleteMessageSeconds.IsUnknown() {
		v := int(plan.DeleteMessageSeconds.ValueInt64())
		params.DeleteMessageSeconds = &v
	}

	err := r.client.CreateGuildBan(ctx, guildID, userID, params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ban", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.GuildID.ValueString(), plan.UserID.ValueString()))

	// The reason is write-only at creation time via the audit log header.
	// We preserve it from the plan. The API returns it in the ban object.
	// Read it back to capture the server-side state.
	ban, err := r.client.GetGuildBan(ctx, guildID, userID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ban after creation", err.Error())
		return
	}

	flattenBan(ban, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state from the API.
func (r *banResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state banModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guildID := discord.Snowflake(state.GuildID.ValueString())
	userID := discord.Snowflake(state.UserID.ValueString())

	ban, err := r.client.GetGuildBan(ctx, guildID, userID)
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading ban", err.Error())
		return
	}

	flattenBan(ban, &state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is not supported for bans; all attributes are ForceNew.
func (r *banResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"discord_ban does not support in-place updates. All changes require replacement.",
	)
}

// Delete removes the guild ban (unbans the user).
func (r *banResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state banModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guildID := discord.Snowflake(state.GuildID.ValueString())
	userID := discord.Snowflake(state.UserID.ValueString())

	err := r.client.RemoveGuildBan(ctx, guildID, userID)
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting ban", err.Error())
	}
}

// ImportState supports importing by guild_id/user_id.
func (r *banResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: guild_id/user_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// --- Helper functions ---

func flattenBan(ban *discord.Ban, model *banModel) {
	// Preserve guild_id and user_id from existing model (they are the key).
	if ban.User != nil {
		model.UserID = types.StringValue(ban.User.ID.String())
	}

	model.ID = types.StringValue(fmt.Sprintf("%s/%s", model.GuildID.ValueString(), model.UserID.ValueString()))

	if ban.Reason != nil {
		model.Reason = types.StringValue(*ban.Reason)
	} else {
		model.Reason = types.StringNull()
	}

	// delete_message_seconds is write-only; preserve from plan/state.
}
