package application

import (
	"context"

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
	_ resource.Resource                = &applicationResource{}
	_ resource.ResourceWithConfigure   = &applicationResource{}
	_ resource.ResourceWithImportState = &applicationResource{}
)

// NewApplicationResource is a constructor that returns a new application resource.
func NewApplicationResource() resource.Resource {
	return &applicationResource{}
}

// applicationResource is the resource implementation.
type applicationResource struct {
	client        *discord.Client
	applicationID string
}

// applicationResourceModel maps the resource schema data.
type applicationResourceModel struct {
	ID                             types.String `tfsdk:"id"`
	Name                           types.String `tfsdk:"name"`
	Description                    types.String `tfsdk:"description"`
	InteractionsEndpointURL        types.String `tfsdk:"interactions_endpoint_url"`
	RoleConnectionsVerificationURL types.String `tfsdk:"role_connections_verification_url"`
	CustomInstallURL               types.String `tfsdk:"custom_install_url"`
	Tags                           types.List   `tfsdk:"tags"`
	BotPublic                      types.Bool   `tfsdk:"bot_public"`
	BotRequireCodeGrant            types.Bool   `tfsdk:"bot_require_code_grant"`
	Icon                           types.String `tfsdk:"icon"`
	CoverImage                     types.String `tfsdk:"cover_image"`
	TermsOfServiceURL              types.String `tfsdk:"terms_of_service_url"`
	PrivacyPolicyURL               types.String `tfsdk:"privacy_policy_url"`
	Flags                          types.Int64  `tfsdk:"flags"`
}

// Metadata returns the resource type name.
func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

// Schema defines the schema for the resource.
func (r *applicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord application's settings. The application cannot be created or deleted via the API. " +
			"Create behaves like Update (PATCH). Delete is a no-op.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the application. Set from provider configuration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the application.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the application.",
				Optional:    true,
				Computed:    true,
			},
			"interactions_endpoint_url": schema.StringAttribute{
				Description: "The interactions endpoint URL for the application.",
				Optional:    true,
			},
			"role_connections_verification_url": schema.StringAttribute{
				Description: "The role connections verification URL for the application.",
				Optional:    true,
			},
			"custom_install_url": schema.StringAttribute{
				Description: "The custom install URL for the application.",
				Optional:    true,
			},
			"tags": schema.ListAttribute{
				Description: "Up to 5 tags describing the content and functionality of the application.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"bot_public": schema.BoolAttribute{
				Description: "Whether the bot is public.",
				Optional:    true,
				Computed:    true,
			},
			"bot_require_code_grant": schema.BoolAttribute{
				Description: "Whether the bot requires the OAuth2 code grant.",
				Optional:    true,
				Computed:    true,
			},
			"icon": schema.StringAttribute{
				Description: "The base64 encoded icon image for the application.",
				Optional:    true,
			},
			"cover_image": schema.StringAttribute{
				Description: "The base64 encoded cover image for the application.",
				Optional:    true,
			},
			"terms_of_service_url": schema.StringAttribute{
				Description: "The URL of the application's terms of service.",
				Optional:    true,
			},
			"privacy_policy_url": schema.StringAttribute{
				Description: "The URL of the application's privacy policy.",
				Optional:    true,
			},
			"flags": schema.Int64Attribute{
				Description: "The application's public flags.",
				Computed:    true,
			},
		},
	}
}

// Configure sets the provider data on the resource.
func (r *applicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	data := common.ProviderDataFromConfig(req.ProviderData, &resp.Diagnostics)
	if data != nil {
		r.client = data.Client
		r.applicationID = data.ApplicationID
	}
}

// Create behaves like Update for an application (PATCH the current application).
func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan applicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.applicationID == "" {
		resp.Diagnostics.AddError("Missing Application ID", "The application_id must be set in the provider configuration to manage the application.")
		return
	}

	params, diags := r.buildEditParams(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.EditCurrentApplication(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", err.Error())
		return
	}

	resp.Diagnostics.Append(r.flattenApplication(ctx, app, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state applicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.GetCurrentApplication(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading application", err.Error())
		return
	}

	resp.Diagnostics.Append(r.flattenApplication(ctx, app, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan applicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params, diags := r.buildEditParams(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.EditCurrentApplication(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", err.Error())
		return
	}

	resp.Diagnostics.Append(r.flattenApplication(ctx, app, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is a no-op for applications (cannot be deleted via the API).
func (r *applicationResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Applications cannot be deleted via the Discord API. This is a no-op.
}

// ImportState imports the resource state using the application ID.
func (r *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildEditParams constructs the edit parameters from the plan.
func (r *applicationResource) buildEditParams(ctx context.Context, plan *applicationResourceModel) (*discord.EditApplicationParams, diag.Diagnostics) {
	var diags diag.Diagnostics
	params := &discord.EditApplicationParams{}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		v := plan.Description.ValueString()
		params.Description = &v
	}

	if !plan.InteractionsEndpointURL.IsNull() && !plan.InteractionsEndpointURL.IsUnknown() {
		v := plan.InteractionsEndpointURL.ValueString()
		params.InteractionsEndpointURL = &v
	}

	if !plan.RoleConnectionsVerificationURL.IsNull() && !plan.RoleConnectionsVerificationURL.IsUnknown() {
		v := plan.RoleConnectionsVerificationURL.ValueString()
		params.RoleConnectionsVerificationURL = &v
	}

	if !plan.CustomInstallURL.IsNull() && !plan.CustomInstallURL.IsUnknown() {
		v := plan.CustomInstallURL.ValueString()
		params.CustomInstallURL = &v
	}

	if !plan.Icon.IsNull() && !plan.Icon.IsUnknown() {
		v := plan.Icon.ValueString()
		params.Icon = &v
	}

	if !plan.CoverImage.IsNull() && !plan.CoverImage.IsUnknown() {
		v := plan.CoverImage.ValueString()
		params.CoverImage = &v
	}

	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags []string
		diags.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		params.Tags = tags
	}

	return params, diags
}

// flattenApplication maps the API response to the Terraform state model.
func (r *applicationResource) flattenApplication(ctx context.Context, app *discord.Application, model *applicationResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.ID = types.StringValue(app.ID.String())
	model.Name = types.StringValue(app.Name)
	model.Description = types.StringValue(app.Description)
	model.BotPublic = types.BoolValue(app.BotPublic)
	model.BotRequireCodeGrant = types.BoolValue(app.BotRequireCodeGrant)

	if app.InteractionsEndpointURL != nil {
		model.InteractionsEndpointURL = types.StringValue(*app.InteractionsEndpointURL)
	} else {
		model.InteractionsEndpointURL = types.StringNull()
	}

	if app.RoleConnectionsVerificationURL != nil {
		model.RoleConnectionsVerificationURL = types.StringValue(*app.RoleConnectionsVerificationURL)
	} else {
		model.RoleConnectionsVerificationURL = types.StringNull()
	}

	if app.CustomInstallURL != nil {
		model.CustomInstallURL = types.StringValue(*app.CustomInstallURL)
	} else {
		model.CustomInstallURL = types.StringNull()
	}

	if app.Icon != nil {
		model.Icon = types.StringValue(*app.Icon)
	} else {
		model.Icon = types.StringNull()
	}

	if app.CoverImage != nil {
		model.CoverImage = types.StringValue(*app.CoverImage)
	} else {
		model.CoverImage = types.StringNull()
	}

	if app.TermsOfServiceURL != nil {
		model.TermsOfServiceURL = types.StringValue(*app.TermsOfServiceURL)
	} else {
		model.TermsOfServiceURL = types.StringNull()
	}

	if app.PrivacyPolicyURL != nil {
		model.PrivacyPolicyURL = types.StringValue(*app.PrivacyPolicyURL)
	} else {
		model.PrivacyPolicyURL = types.StringNull()
	}

	if app.Flags != nil {
		model.Flags = types.Int64Value(int64(*app.Flags))
	} else {
		model.Flags = types.Int64Value(0)
	}

	if len(app.Tags) > 0 {
		tagsList, d := types.ListValueFrom(ctx, types.StringType, app.Tags)
		diags.Append(d...)
		model.Tags = tagsList
	} else {
		model.Tags = types.ListNull(types.StringType)
	}

	return diags
}
