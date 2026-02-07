package guild

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &permissionDataSource{}
)

// permissionDataSource is the data source implementation.
type permissionDataSource struct{}

// permissionFlag maps a permission name to its bit position.
type permissionFlag struct {
	Name string
	Bit  uint64
}

// allPermissions lists all Discord permission flags in order.
var allPermissions = []permissionFlag{
	{"create_instant_invite", 0x0000000000000001},
	{"kick_members", 0x0000000000000002},
	{"ban_members", 0x0000000000000004},
	{"administrator", 0x0000000000000008},
	{"manage_channels", 0x0000000000000010},
	{"manage_guild", 0x0000000000000020},
	{"add_reactions", 0x0000000000000040},
	{"view_audit_log", 0x0000000000000080},
	{"priority_speaker", 0x0000000000000100},
	{"stream", 0x0000000000000200},
	{"view_channel", 0x0000000000000400},
	{"send_messages", 0x0000000000000800},
	{"send_tts_messages", 0x0000000000001000},
	{"manage_messages", 0x0000000000002000},
	{"embed_links", 0x0000000000004000},
	{"attach_files", 0x0000000000008000},
	{"read_message_history", 0x0000000000010000},
	{"mention_everyone", 0x0000000000020000},
	{"use_external_emojis", 0x0000000000040000},
	{"view_guild_insights", 0x0000000000080000},
	{"connect", 0x0000000000100000},
	{"speak", 0x0000000000200000},
	{"mute_members", 0x0000000000400000},
	{"deafen_members", 0x0000000000800000},
	{"move_members", 0x0000000001000000},
	{"use_vad", 0x0000000002000000},
	{"change_nickname", 0x0000000004000000},
	{"manage_nicknames", 0x0000000008000000},
	{"manage_roles", 0x0000000010000000},
	{"manage_webhooks", 0x0000000020000000},
	{"manage_guild_expressions", 0x0000000040000000},
	{"use_application_commands", 0x0000000080000000},
	{"request_to_speak", 0x0000000100000000},
	{"manage_events", 0x0000000200000000},
	{"manage_threads", 0x0000000400000000},
	{"create_public_threads", 0x0000000800000000},
	{"create_private_threads", 0x0000001000000000},
	{"use_external_stickers", 0x0000002000000000},
	{"send_messages_in_threads", 0x0000004000000000},
	{"use_embedded_activities", 0x0000008000000000},
	{"moderate_members", 0x0000010000000000},
	{"view_creator_monetization_analytics", 0x0000020000000000},
	{"use_soundboard", 0x0000040000000000},
	{"use_external_sounds", 0x0000200000000000},
	{"send_voice_messages", 0x0000400000000000},
}

// permissionDataSourceModel maps the data source schema data.
type permissionDataSourceModel struct {
	CreateInstantInvite              types.Bool   `tfsdk:"create_instant_invite"`
	KickMembers                      types.Bool   `tfsdk:"kick_members"`
	BanMembers                       types.Bool   `tfsdk:"ban_members"`
	Administrator                    types.Bool   `tfsdk:"administrator"`
	ManageChannels                   types.Bool   `tfsdk:"manage_channels"`
	ManageGuild                      types.Bool   `tfsdk:"manage_guild"`
	AddReactions                     types.Bool   `tfsdk:"add_reactions"`
	ViewAuditLog                     types.Bool   `tfsdk:"view_audit_log"`
	PrioritySpeaker                  types.Bool   `tfsdk:"priority_speaker"`
	Stream                           types.Bool   `tfsdk:"stream"`
	ViewChannel                      types.Bool   `tfsdk:"view_channel"`
	SendMessages                     types.Bool   `tfsdk:"send_messages"`
	SendTTSMessages                  types.Bool   `tfsdk:"send_tts_messages"`
	ManageMessages                   types.Bool   `tfsdk:"manage_messages"`
	EmbedLinks                       types.Bool   `tfsdk:"embed_links"`
	AttachFiles                      types.Bool   `tfsdk:"attach_files"`
	ReadMessageHistory               types.Bool   `tfsdk:"read_message_history"`
	MentionEveryone                  types.Bool   `tfsdk:"mention_everyone"`
	UseExternalEmojis                types.Bool   `tfsdk:"use_external_emojis"`
	ViewGuildInsights                types.Bool   `tfsdk:"view_guild_insights"`
	Connect                          types.Bool   `tfsdk:"connect"`
	Speak                            types.Bool   `tfsdk:"speak"`
	MuteMembers                      types.Bool   `tfsdk:"mute_members"`
	DeafenMembers                    types.Bool   `tfsdk:"deafen_members"`
	MoveMembers                      types.Bool   `tfsdk:"move_members"`
	UseVAD                           types.Bool   `tfsdk:"use_vad"`
	ChangeNickname                   types.Bool   `tfsdk:"change_nickname"`
	ManageNicknames                  types.Bool   `tfsdk:"manage_nicknames"`
	ManageRoles                      types.Bool   `tfsdk:"manage_roles"`
	ManageWebhooks                   types.Bool   `tfsdk:"manage_webhooks"`
	ManageGuildExpressions           types.Bool   `tfsdk:"manage_guild_expressions"`
	UseApplicationCommands           types.Bool   `tfsdk:"use_application_commands"`
	RequestToSpeak                   types.Bool   `tfsdk:"request_to_speak"`
	ManageEvents                     types.Bool   `tfsdk:"manage_events"`
	ManageThreads                    types.Bool   `tfsdk:"manage_threads"`
	CreatePublicThreads              types.Bool   `tfsdk:"create_public_threads"`
	CreatePrivateThreads             types.Bool   `tfsdk:"create_private_threads"`
	UseExternalStickers              types.Bool   `tfsdk:"use_external_stickers"`
	SendMessagesInThreads            types.Bool   `tfsdk:"send_messages_in_threads"`
	UseEmbeddedActivities            types.Bool   `tfsdk:"use_embedded_activities"`
	ModerateMembers                  types.Bool   `tfsdk:"moderate_members"`
	ViewCreatorMonetizationAnalytics types.Bool   `tfsdk:"view_creator_monetization_analytics"`
	UseSoundboard                    types.Bool   `tfsdk:"use_soundboard"`
	UseExternalSounds                types.Bool   `tfsdk:"use_external_sounds"`
	SendVoiceMessages                types.Bool   `tfsdk:"send_voice_messages"`
	AllowBits                        types.String `tfsdk:"allow_bits"`
	DenyBits                         types.String `tfsdk:"deny_bits"`
}

// NewPermissionDataSource returns a new permission data source.
func NewPermissionDataSource() datasource.DataSource {
	return &permissionDataSource{}
}

// Metadata returns the data source type name.
func (d *permissionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

// Schema defines the schema for the data source.
func (d *permissionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A utility data source that computes a Discord permission integer from named permissions. " +
			"This is not backed by an API call. Set permission flags to true to include them in " +
			"allow_bits, or false to include them in deny_bits. Unset (null) flags are not included in either.",
		Attributes: map[string]schema.Attribute{
			"create_instant_invite": schema.BoolAttribute{
				Description: "Allows creation of instant invites.",
				Optional:    true,
			},
			"kick_members": schema.BoolAttribute{
				Description: "Allows kicking members.",
				Optional:    true,
			},
			"ban_members": schema.BoolAttribute{
				Description: "Allows banning members.",
				Optional:    true,
			},
			"administrator": schema.BoolAttribute{
				Description: "Allows all permissions and bypasses channel permission overwrites.",
				Optional:    true,
			},
			"manage_channels": schema.BoolAttribute{
				Description: "Allows management and editing of channels.",
				Optional:    true,
			},
			"manage_guild": schema.BoolAttribute{
				Description: "Allows management and editing of the guild.",
				Optional:    true,
			},
			"add_reactions": schema.BoolAttribute{
				Description: "Allows for the addition of reactions to messages.",
				Optional:    true,
			},
			"view_audit_log": schema.BoolAttribute{
				Description: "Allows for viewing of audit logs.",
				Optional:    true,
			},
			"priority_speaker": schema.BoolAttribute{
				Description: "Allows for using priority speaker in a voice channel.",
				Optional:    true,
			},
			"stream": schema.BoolAttribute{
				Description: "Allows the user to go live.",
				Optional:    true,
			},
			"view_channel": schema.BoolAttribute{
				Description: "Allows guild members to view a channel.",
				Optional:    true,
			},
			"send_messages": schema.BoolAttribute{
				Description: "Allows for sending messages in a channel.",
				Optional:    true,
			},
			"send_tts_messages": schema.BoolAttribute{
				Description: "Allows for sending of TTS messages.",
				Optional:    true,
			},
			"manage_messages": schema.BoolAttribute{
				Description: "Allows for deletion of other users messages.",
				Optional:    true,
			},
			"embed_links": schema.BoolAttribute{
				Description: "Links sent by users with this permission will be auto-embedded.",
				Optional:    true,
			},
			"attach_files": schema.BoolAttribute{
				Description: "Allows for uploading images and files.",
				Optional:    true,
			},
			"read_message_history": schema.BoolAttribute{
				Description: "Allows for reading of message history.",
				Optional:    true,
			},
			"mention_everyone": schema.BoolAttribute{
				Description: "Allows for using the @everyone tag to notify all users in a channel.",
				Optional:    true,
			},
			"use_external_emojis": schema.BoolAttribute{
				Description: "Allows the usage of custom emojis from other servers.",
				Optional:    true,
			},
			"view_guild_insights": schema.BoolAttribute{
				Description: "Allows for viewing guild insights.",
				Optional:    true,
			},
			"connect": schema.BoolAttribute{
				Description: "Allows for joining of a voice channel.",
				Optional:    true,
			},
			"speak": schema.BoolAttribute{
				Description: "Allows for speaking in a voice channel.",
				Optional:    true,
			},
			"mute_members": schema.BoolAttribute{
				Description: "Allows for muting members in a voice channel.",
				Optional:    true,
			},
			"deafen_members": schema.BoolAttribute{
				Description: "Allows for deafening of members in a voice channel.",
				Optional:    true,
			},
			"move_members": schema.BoolAttribute{
				Description: "Allows for moving of members between voice channels.",
				Optional:    true,
			},
			"use_vad": schema.BoolAttribute{
				Description: "Allows for using voice-activity-detection in a voice channel.",
				Optional:    true,
			},
			"change_nickname": schema.BoolAttribute{
				Description: "Allows for modification of own nickname.",
				Optional:    true,
			},
			"manage_nicknames": schema.BoolAttribute{
				Description: "Allows for modification of other users nicknames.",
				Optional:    true,
			},
			"manage_roles": schema.BoolAttribute{
				Description: "Allows management and editing of roles.",
				Optional:    true,
			},
			"manage_webhooks": schema.BoolAttribute{
				Description: "Allows management and editing of webhooks.",
				Optional:    true,
			},
			"manage_guild_expressions": schema.BoolAttribute{
				Description: "Allows management and editing of emojis, stickers, and soundboard sounds.",
				Optional:    true,
			},
			"use_application_commands": schema.BoolAttribute{
				Description: "Allows members to use application commands.",
				Optional:    true,
			},
			"request_to_speak": schema.BoolAttribute{
				Description: "Allows for requesting to speak in stage channels.",
				Optional:    true,
			},
			"manage_events": schema.BoolAttribute{
				Description: "Allows for management of guild scheduled events.",
				Optional:    true,
			},
			"manage_threads": schema.BoolAttribute{
				Description: "Allows for deleting and archiving threads, and viewing all private threads.",
				Optional:    true,
			},
			"create_public_threads": schema.BoolAttribute{
				Description: "Allows for creating public threads.",
				Optional:    true,
			},
			"create_private_threads": schema.BoolAttribute{
				Description: "Allows for creating private threads.",
				Optional:    true,
			},
			"use_external_stickers": schema.BoolAttribute{
				Description: "Allows the usage of custom stickers from other servers.",
				Optional:    true,
			},
			"send_messages_in_threads": schema.BoolAttribute{
				Description: "Allows for sending messages in threads.",
				Optional:    true,
			},
			"use_embedded_activities": schema.BoolAttribute{
				Description: "Allows for using Activities in a voice channel.",
				Optional:    true,
			},
			"moderate_members": schema.BoolAttribute{
				Description: "Allows for timing out users.",
				Optional:    true,
			},
			"view_creator_monetization_analytics": schema.BoolAttribute{
				Description: "Allows for viewing role subscription insights.",
				Optional:    true,
			},
			"use_soundboard": schema.BoolAttribute{
				Description: "Allows for using soundboard in a voice channel.",
				Optional:    true,
			},
			"use_external_sounds": schema.BoolAttribute{
				Description: "Allows the usage of custom soundboard sounds from other servers.",
				Optional:    true,
			},
			"send_voice_messages": schema.BoolAttribute{
				Description: "Allows sending voice messages.",
				Optional:    true,
			},
			"allow_bits": schema.StringAttribute{
				Description: "The computed allow permission integer as a string.",
				Computed:    true,
			},
			"deny_bits": schema.StringAttribute{
				Description: "The computed deny permission integer as a string.",
				Computed:    true,
			},
		},
	}
}

// getBoolField returns the types.Bool value for a given permission name from the model.
func getBoolField(m *permissionDataSourceModel, name string) types.Bool {
	switch name {
	case "create_instant_invite":
		return m.CreateInstantInvite
	case "kick_members":
		return m.KickMembers
	case "ban_members":
		return m.BanMembers
	case "administrator":
		return m.Administrator
	case "manage_channels":
		return m.ManageChannels
	case "manage_guild":
		return m.ManageGuild
	case "add_reactions":
		return m.AddReactions
	case "view_audit_log":
		return m.ViewAuditLog
	case "priority_speaker":
		return m.PrioritySpeaker
	case "stream":
		return m.Stream
	case "view_channel":
		return m.ViewChannel
	case "send_messages":
		return m.SendMessages
	case "send_tts_messages":
		return m.SendTTSMessages
	case "manage_messages":
		return m.ManageMessages
	case "embed_links":
		return m.EmbedLinks
	case "attach_files":
		return m.AttachFiles
	case "read_message_history":
		return m.ReadMessageHistory
	case "mention_everyone":
		return m.MentionEveryone
	case "use_external_emojis":
		return m.UseExternalEmojis
	case "view_guild_insights":
		return m.ViewGuildInsights
	case "connect":
		return m.Connect
	case "speak":
		return m.Speak
	case "mute_members":
		return m.MuteMembers
	case "deafen_members":
		return m.DeafenMembers
	case "move_members":
		return m.MoveMembers
	case "use_vad":
		return m.UseVAD
	case "change_nickname":
		return m.ChangeNickname
	case "manage_nicknames":
		return m.ManageNicknames
	case "manage_roles":
		return m.ManageRoles
	case "manage_webhooks":
		return m.ManageWebhooks
	case "manage_guild_expressions":
		return m.ManageGuildExpressions
	case "use_application_commands":
		return m.UseApplicationCommands
	case "request_to_speak":
		return m.RequestToSpeak
	case "manage_events":
		return m.ManageEvents
	case "manage_threads":
		return m.ManageThreads
	case "create_public_threads":
		return m.CreatePublicThreads
	case "create_private_threads":
		return m.CreatePrivateThreads
	case "use_external_stickers":
		return m.UseExternalStickers
	case "send_messages_in_threads":
		return m.SendMessagesInThreads
	case "use_embedded_activities":
		return m.UseEmbeddedActivities
	case "moderate_members":
		return m.ModerateMembers
	case "view_creator_monetization_analytics":
		return m.ViewCreatorMonetizationAnalytics
	case "use_soundboard":
		return m.UseSoundboard
	case "use_external_sounds":
		return m.UseExternalSounds
	case "send_voice_messages":
		return m.SendVoiceMessages
	default:
		return types.BoolNull()
	}
}

// Read computes the permission integers from the flags.
func (d *permissionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config permissionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var allowBits uint64
	var denyBits uint64

	for _, perm := range allPermissions {
		field := getBoolField(&config, perm.Name)
		if field.IsNull() || field.IsUnknown() {
			continue
		}
		if field.ValueBool() {
			allowBits |= perm.Bit
		} else {
			denyBits |= perm.Bit
		}
	}

	config.AllowBits = types.StringValue(fmt.Sprintf("%d", allowBits))
	config.DenyBits = types.StringValue(fmt.Sprintf("%d", denyBits))

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
