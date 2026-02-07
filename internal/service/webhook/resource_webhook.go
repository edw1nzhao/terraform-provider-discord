package webhook

import (
	"context"
	"fmt"

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
	_ resource.Resource                = &webhookResource{}
	_ resource.ResourceWithConfigure   = &webhookResource{}
	_ resource.ResourceWithImportState = &webhookResource{}
)

// NewWebhookResource is a constructor that returns a new webhook resource.
func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

// webhookResource is the resource implementation.
type webhookResource struct {
	client *discord.Client
}

// webhookResourceModel maps the resource schema data.
type webhookResourceModel struct {
	ID        types.String `tfsdk:"id"`
	ChannelID types.String `tfsdk:"channel_id"`
	Name      types.String `tfsdk:"name"`
	Avatar    types.String `tfsdk:"avatar"`
	Type      types.Int64  `tfsdk:"type"`
	GuildID   types.String `tfsdk:"guild_id"`
	Token     types.String `tfsdk:"token"`
	URL       types.String `tfsdk:"url"`
}

// Metadata returns the resource type name.
func (r *webhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

// Schema defines the schema for the resource.
func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord webhook.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the webhook.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"channel_id": schema.StringAttribute{
				Description: "The ID of the channel the webhook belongs to.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the webhook (1-80 characters).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 80),
				},
			},
			"avatar": schema.StringAttribute{
				Description: "The base64 encoded image for the webhook avatar.",
				Optional:    true,
			},
			"type": schema.Int64Attribute{
				Description: "The type of the webhook.",
				Computed:    true,
			},
			"guild_id": schema.StringAttribute{
				Description: "The guild ID the webhook belongs to.",
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "The secure token of the webhook.",
				Computed:    true,
				Sensitive:   true,
			},
			"url": schema.StringAttribute{
				Description: "The URL of the webhook.",
				Computed:    true,
			},
		},
	}
}

// Configure sets the provider data on the resource.
func (r *webhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource and sets the initial Terraform state.
func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan webhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.CreateWebhookParams{
		Name: plan.Name.ValueString(),
	}
	if !plan.Avatar.IsNull() && !plan.Avatar.IsUnknown() {
		v := plan.Avatar.ValueString()
		params.Avatar = &v
	}

	webhook, err := r.client.CreateWebhook(ctx, discord.Snowflake(plan.ChannelID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating webhook", err.Error())
		return
	}

	r.flattenWebhook(webhook, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhook, err := r.client.GetWebhook(ctx, discord.Snowflake(state.ID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading webhook", err.Error())
		return
	}

	r.flattenWebhook(webhook, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan webhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state webhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.ModifyWebhookParams{}

	name := plan.Name.ValueString()
	params.Name = &name

	if !plan.Avatar.IsNull() && !plan.Avatar.IsUnknown() {
		v := plan.Avatar.ValueString()
		params.Avatar = &v
	}

	channelID := discord.Snowflake(plan.ChannelID.ValueString())
	params.ChannelID = &channelID

	webhook, err := r.client.ModifyWebhook(ctx, discord.Snowflake(state.ID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating webhook", err.Error())
		return
	}

	r.flattenWebhook(webhook, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state.
func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWebhook(ctx, discord.Snowflake(state.ID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting webhook", err.Error())
	}
}

// ImportState imports the resource state.
func (r *webhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// flattenWebhook maps the API response to the Terraform state model.
func (r *webhookResource) flattenWebhook(webhook *discord.Webhook, model *webhookResourceModel) {
	model.ID = types.StringValue(webhook.ID.String())
	model.Type = types.Int64Value(int64(webhook.Type))

	if webhook.ChannelID != nil {
		model.ChannelID = types.StringValue(webhook.ChannelID.String())
	}
	if webhook.GuildID != nil {
		model.GuildID = types.StringValue(webhook.GuildID.String())
	} else {
		model.GuildID = types.StringNull()
	}
	if webhook.Name != nil {
		model.Name = types.StringValue(*webhook.Name)
	}
	if webhook.Avatar != nil {
		model.Avatar = types.StringValue(*webhook.Avatar)
	}
	if webhook.Token != nil {
		model.Token = types.StringValue(*webhook.Token)
		model.URL = types.StringValue(fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", webhook.ID, *webhook.Token))
	} else {
		model.Token = types.StringNull()
		model.URL = types.StringNull()
	}
}
