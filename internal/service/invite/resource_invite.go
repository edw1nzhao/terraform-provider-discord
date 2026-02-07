package invite

import (
	"context"
	"fmt"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &inviteResource{}
	_ resource.ResourceWithConfigure   = &inviteResource{}
	_ resource.ResourceWithImportState = &inviteResource{}
)

// NewInviteResource is a constructor that returns a new invite resource.
func NewInviteResource() resource.Resource {
	return &inviteResource{}
}

// inviteResource is the resource implementation.
type inviteResource struct {
	client *discord.Client
}

// inviteResourceModel maps the resource schema data.
type inviteResourceModel struct {
	ID        types.String `tfsdk:"id"`
	ChannelID types.String `tfsdk:"channel_id"`
	MaxAge    types.Int64  `tfsdk:"max_age"`
	MaxUses   types.Int64  `tfsdk:"max_uses"`
	Temporary types.Bool   `tfsdk:"temporary"`
	Unique    types.Bool   `tfsdk:"unique"`
	Uses      types.Int64  `tfsdk:"uses"`
	URL       types.String `tfsdk:"url"`
}

// Metadata returns the resource type name.
func (r *inviteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_invite"
}

// Schema defines the schema for the resource.
func (r *inviteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord channel invite. This resource supports Create, Read, and Delete only. " +
			"Updates are not supported; any change to configuration will force recreation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The invite code.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"channel_id": schema.StringAttribute{
				Description: "The ID of the channel to create the invite for.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"max_age": schema.Int64Attribute{
				Description: "Duration of invite in seconds before expiry, or 0 for never. Default: 86400 (24 hours).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(86400),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"max_uses": schema.Int64Attribute{
				Description: "Max number of uses, or 0 for unlimited. Default: 0.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"temporary": schema.BoolAttribute{
				Description: "Whether this invite only grants temporary membership. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"unique": schema.BoolAttribute{
				Description: "If true, don't try to reuse a similar invite. Default: false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"uses": schema.Int64Attribute{
				Description: "Number of times this invite has been used.",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "The URL of the invite.",
				Computed:    true,
			},
		},
	}
}

// Configure sets the provider data on the resource.
func (r *inviteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource and sets the initial Terraform state.
func (r *inviteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan inviteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	maxAge := int(plan.MaxAge.ValueInt64())
	maxUses := int(plan.MaxUses.ValueInt64())
	temporary := plan.Temporary.ValueBool()
	unique := plan.Unique.ValueBool()

	params := &discord.CreateInviteParams{
		MaxAge:    &maxAge,
		MaxUses:   &maxUses,
		Temporary: &temporary,
		Unique:    &unique,
	}

	invite, err := r.client.CreateChannelInvite(ctx, discord.Snowflake(plan.ChannelID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating invite", err.Error())
		return
	}

	plan.ID = types.StringValue(invite.Code)
	plan.Uses = types.Int64Value(int64(invite.Uses))
	plan.URL = types.StringValue(fmt.Sprintf("https://discord.gg/%s", invite.Code))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *inviteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state inviteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	invite, err := r.client.GetInvite(ctx, state.ID.ValueString())
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading invite", err.Error())
		return
	}

	state.ID = types.StringValue(invite.Code)
	if invite.Channel != nil {
		state.ChannelID = types.StringValue(invite.Channel.ID.String())
	}
	state.MaxAge = types.Int64Value(int64(invite.MaxAge))
	state.MaxUses = types.Int64Value(int64(invite.MaxUses))
	state.Temporary = types.BoolValue(invite.Temporary)
	state.Uses = types.Int64Value(int64(invite.Uses))
	state.URL = types.StringValue(fmt.Sprintf("https://discord.gg/%s", invite.Code))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is not supported for invites. Changes force recreation.
func (r *inviteResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Discord invites cannot be updated. Changes require recreation.",
	)
}

// Delete deletes the resource and removes the Terraform state.
func (r *inviteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state inviteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteInvite(ctx, state.ID.ValueString())
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting invite", err.Error())
	}
}

// ImportState imports the resource state using the invite code.
func (r *inviteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
