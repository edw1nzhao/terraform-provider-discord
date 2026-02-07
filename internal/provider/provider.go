package provider

import (
	"context"
	"os"

	"github.com/edw1nzhao/terraform-provider-discord/internal/conns"
	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/application"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/application_command"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/automod"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/ban"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/channel"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/emoji"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/guild"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/invite"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/member"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/message"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/onboarding"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/role"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/scheduled_event"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/soundboard"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/stage"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/sticker"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/template"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/user"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/voice"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/webhook"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/welcome_screen"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/widget"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure discordProvider satisfies the provider.Provider interface.
var _ provider.Provider = &discordProvider{}

// discordProvider implements the Discord Terraform provider.
type discordProvider struct {
	// "dev" for local builds, "test" for acceptance tests, or semver on release.
	version string
}

// discordProviderModel describes the provider configuration data model.
type discordProviderModel struct {
	// Token is the Discord bot token used to authenticate API requests.
	// Required. Can also be set via the DISCORD_TOKEN environment variable.
	Token types.String `tfsdk:"token"`

	// ApplicationID is the Discord application (bot) ID. Optional.
	// Can also be set via the DISCORD_APPLICATION_ID environment variable.
	// Needed for application command resources.
	ApplicationID types.String `tfsdk:"application_id"`
}

// New returns a factory function that creates a new instance of the provider.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &discordProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *discordProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "discord"
	resp.Version = p.version
}

// Schema defines the provider-level configuration schema.
func (p *discordProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage Discord servers (guilds), channels, roles, and related resources.",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "The Discord bot token used to authenticate API requests. " +
					"Can also be set via the DISCORD_TOKEN environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"application_id": schema.StringAttribute{
				Description: "The Discord application (bot) ID. " +
					"Can also be set via the DISCORD_APPLICATION_ID environment variable. " +
					"Required for managing application command resources.",
				Optional: true,
			},
		},
	}
}

// Configure prepares the Discord API client for resources and data sources.
func (p *discordProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config discordProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve the token: config value takes precedence, then env var.
	token := os.Getenv("DISCORD_TOKEN")
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"Missing Discord Bot Token",
			"The provider requires a Discord bot token to authenticate API requests. "+
				"Set the token in the provider configuration block or via the DISCORD_TOKEN environment variable.",
		)
		return
	}

	// Resolve the application ID: config value takes precedence, then env var.
	applicationID := os.Getenv("DISCORD_APPLICATION_ID")
	if !config.ApplicationID.IsNull() {
		applicationID = config.ApplicationID.ValueString()
	}

	// Create the Discord REST client.
	client := discord.NewClient(token, p.version)

	// Store the client and application ID so resources and data sources can
	// retrieve them.
	data := &conns.ProviderData{
		Client:        client,
		ApplicationID: applicationID,
	}

	resp.DataSourceData = data
	resp.ResourceData = data
}

// Resources defines the resources implemented in the provider.
func (p *discordProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		guild.NewGuildResource,
		channel.NewChannelResource,
		channel.NewChannelPermissionResource,
		role.NewRoleResource,
		member.NewMemberRolesResource,
		soundboard.NewSoundboardSoundResource,
		message.NewMessageResource,
		widget.NewGuildWidgetResource,
		webhook.NewWebhookResource,
		invite.NewInviteResource,
		emoji.NewGuildEmojiResource,
		sticker.NewGuildStickerResource,
		template.NewGuildTemplateResource,
		stage.NewStageInstanceResource,
		application_command.NewGlobalApplicationCommandResource,
		application_command.NewGuildApplicationCommandResource,
		application.NewApplicationResource,
		automod.NewAutoModerationRuleResource,
		ban.NewBanResource,
		scheduled_event.NewGuildScheduledEventResource,
		welcome_screen.NewWelcomeScreenResource,
		onboarding.NewGuildOnboardingResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *discordProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		guild.NewGuildDataSource,
		guild.NewPermissionDataSource,
		guild.NewColorDataSource,
		channel.NewChannelDataSource,
		role.NewRoleDataSource,
		user.NewUserDataSource,
		voice.NewVoiceRegionsDataSource,
	}
}
