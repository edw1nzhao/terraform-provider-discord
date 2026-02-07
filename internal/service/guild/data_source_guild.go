package guild

import (
	"context"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &guildDataSource{}
	_ datasource.DataSourceWithConfigure = &guildDataSource{}
)

// guildDataSource is the data source implementation.
type guildDataSource struct {
	client *discord.Client
}

// guildDataSourceModel maps the data source schema to a Go struct.
type guildDataSourceModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Icon                        types.String `tfsdk:"icon"`
	Splash                      types.String `tfsdk:"splash"`
	Banner                      types.String `tfsdk:"banner"`
	Description                 types.String `tfsdk:"description"`
	OwnerID                     types.String `tfsdk:"owner_id"`
	AFKChannelID                types.String `tfsdk:"afk_channel_id"`
	AFKTimeout                  types.Int64  `tfsdk:"afk_timeout"`
	VerificationLevel           types.Int64  `tfsdk:"verification_level"`
	DefaultMessageNotifications types.Int64  `tfsdk:"default_message_notifications"`
	ExplicitContentFilter       types.Int64  `tfsdk:"explicit_content_filter"`
	SystemChannelID             types.String `tfsdk:"system_channel_id"`
	SystemChannelFlags          types.Int64  `tfsdk:"system_channel_flags"`
	RulesChannelID              types.String `tfsdk:"rules_channel_id"`
	PublicUpdatesChannelID      types.String `tfsdk:"public_updates_channel_id"`
	PreferredLocale             types.String `tfsdk:"preferred_locale"`
	PremiumProgressBarEnabled   types.Bool   `tfsdk:"premium_progress_bar_enabled"`
	SafetyAlertsChannelID       types.String `tfsdk:"safety_alerts_channel_id"`
	PremiumTier                 types.Int64  `tfsdk:"premium_tier"`
	PremiumSubscriptionCount    types.Int64  `tfsdk:"premium_subscription_count"`
	Features                    types.List   `tfsdk:"features"`
	MFALevel                    types.Int64  `tfsdk:"mfa_level"`
}

// NewGuildDataSource returns a new guild data source.
func NewGuildDataSource() datasource.DataSource {
	return &guildDataSource{}
}

// Metadata returns the data source type name.
func (d *guildDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guild"
}

// Configure adds the provider configured client to the data source.
func (d *guildDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the data source.
func (d *guildDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a Discord guild (server).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the guild.",
				Computed:    true,
			},
			"icon": schema.StringAttribute{
				Description: "The guild icon hash.",
				Computed:    true,
			},
			"splash": schema.StringAttribute{
				Description: "The guild splash hash.",
				Computed:    true,
			},
			"banner": schema.StringAttribute{
				Description: "The guild banner hash.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the guild.",
				Computed:    true,
			},
			"owner_id": schema.StringAttribute{
				Description: "The ID of the guild owner.",
				Computed:    true,
			},
			"afk_channel_id": schema.StringAttribute{
				Description: "The ID of the AFK voice channel.",
				Computed:    true,
			},
			"afk_timeout": schema.Int64Attribute{
				Description: "AFK timeout in seconds.",
				Computed:    true,
			},
			"verification_level": schema.Int64Attribute{
				Description: "The verification level required for the guild.",
				Computed:    true,
			},
			"default_message_notifications": schema.Int64Attribute{
				Description: "The default message notification level.",
				Computed:    true,
			},
			"explicit_content_filter": schema.Int64Attribute{
				Description: "The explicit content filter level.",
				Computed:    true,
			},
			"system_channel_id": schema.StringAttribute{
				Description: "The ID of the channel where system messages are sent.",
				Computed:    true,
			},
			"system_channel_flags": schema.Int64Attribute{
				Description: "System channel flags.",
				Computed:    true,
			},
			"rules_channel_id": schema.StringAttribute{
				Description: "The ID of the channel where community rules are displayed.",
				Computed:    true,
			},
			"public_updates_channel_id": schema.StringAttribute{
				Description: "The ID of the channel where admins and moderators receive notices from Discord.",
				Computed:    true,
			},
			"preferred_locale": schema.StringAttribute{
				Description: "The preferred locale of a community guild.",
				Computed:    true,
			},
			"premium_progress_bar_enabled": schema.BoolAttribute{
				Description: "Whether the guild has the boost progress bar enabled.",
				Computed:    true,
			},
			"safety_alerts_channel_id": schema.StringAttribute{
				Description: "The ID of the channel where safety alerts are sent.",
				Computed:    true,
			},
			"premium_tier": schema.Int64Attribute{
				Description: "The premium tier (Server Boost level).",
				Computed:    true,
			},
			"premium_subscription_count": schema.Int64Attribute{
				Description: "The number of boosts this guild currently has.",
				Computed:    true,
			},
			"features": schema.ListAttribute{
				Description: "The list of enabled guild features.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"mfa_level": schema.Int64Attribute{
				Description: "The required MFA level for the guild.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *guildDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config guildDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guild, err := d.client.GetGuild(ctx, discord.Snowflake(config.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Discord Guild",
			"Could not read guild ID "+config.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map the API response to the data source model.
	config.Name = types.StringValue(guild.Name)
	config.OwnerID = types.StringValue(guild.OwnerID.String())
	config.AFKTimeout = types.Int64Value(int64(guild.AFKTimeout))
	config.VerificationLevel = types.Int64Value(int64(guild.VerificationLevel))
	config.DefaultMessageNotifications = types.Int64Value(int64(guild.DefaultMessageNotifications))
	config.ExplicitContentFilter = types.Int64Value(int64(guild.ExplicitContentFilter))
	config.SystemChannelFlags = types.Int64Value(int64(guild.SystemChannelFlags))
	config.PremiumTier = types.Int64Value(int64(guild.PremiumTier))
	config.PremiumSubscriptionCount = types.Int64Value(int64(guild.PremiumSubscriptionCount))
	config.MFALevel = types.Int64Value(int64(guild.MFALevel))
	config.PremiumProgressBarEnabled = types.BoolValue(guild.PremiumProgressBarEnabled)
	config.PreferredLocale = types.StringValue(guild.PreferredLocale)

	// Nullable string fields.
	if guild.Icon != nil {
		config.Icon = types.StringValue(*guild.Icon)
	} else {
		config.Icon = types.StringNull()
	}
	if guild.Splash != nil {
		config.Splash = types.StringValue(*guild.Splash)
	} else {
		config.Splash = types.StringNull()
	}
	if guild.Banner != nil {
		config.Banner = types.StringValue(*guild.Banner)
	} else {
		config.Banner = types.StringNull()
	}
	if guild.Description != nil {
		config.Description = types.StringValue(*guild.Description)
	} else {
		config.Description = types.StringNull()
	}

	// Nullable snowflake fields.
	if guild.AFKChannelID != nil {
		config.AFKChannelID = types.StringValue(guild.AFKChannelID.String())
	} else {
		config.AFKChannelID = types.StringNull()
	}
	if guild.SystemChannelID != nil {
		config.SystemChannelID = types.StringValue(guild.SystemChannelID.String())
	} else {
		config.SystemChannelID = types.StringNull()
	}
	if guild.RulesChannelID != nil {
		config.RulesChannelID = types.StringValue(guild.RulesChannelID.String())
	} else {
		config.RulesChannelID = types.StringNull()
	}
	if guild.PublicUpdatesChannelID != nil {
		config.PublicUpdatesChannelID = types.StringValue(guild.PublicUpdatesChannelID.String())
	} else {
		config.PublicUpdatesChannelID = types.StringNull()
	}
	if guild.SafetyAlertsChannelID != nil {
		config.SafetyAlertsChannelID = types.StringValue(guild.SafetyAlertsChannelID.String())
	} else {
		config.SafetyAlertsChannelID = types.StringNull()
	}

	// Features list.
	features := make([]types.String, len(guild.Features))
	for i, f := range guild.Features {
		features[i] = types.StringValue(f)
	}
	featuresList, diags := types.ListValueFrom(ctx, types.StringType, features)
	resp.Diagnostics.Append(diags...)
	config.Features = featuresList

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
