package discord

import "time"

// Snowflake is a Discord snowflake ID represented as a string.
type Snowflake string

// String returns the string representation of the snowflake.
func (s Snowflake) String() string {
	return string(s)
}

// IsEmpty returns true if the snowflake is empty or zero.
func (s Snowflake) IsEmpty() bool {
	return s == "" || s == "0"
}

// Channel types
const (
	ChannelTypeGuildText         = 0
	ChannelTypeGuildVoice        = 2
	ChannelTypeGuildCategory     = 4
	ChannelTypeGuildAnnouncement = 5
	ChannelTypeGuildStageVoice   = 13
	ChannelTypeGuildForum        = 15
)

// Auto-moderation trigger types
const (
	AutoModTriggerKeyword       = 1
	AutoModTriggerKeywordPreset = 4
	AutoModTriggerMentionSpam   = 5
)

// Auto-moderation action types
const (
	AutoModActionBlockMessage     = 1
	AutoModActionSendAlertMessage = 2
	AutoModActionTimeout          = 3
)

// Auto-moderation event types
const (
	AutoModEventMessageSend = 1
)

// Scheduled event entity types
const (
	ScheduledEventEntityStageInstance = 1
	ScheduledEventEntityVoice        = 2
	ScheduledEventEntityExternal     = 3
)

// Scheduled event privacy levels
const (
	ScheduledEventPrivacyGuildOnly = 2
)

// Scheduled event status
const (
	ScheduledEventStatusScheduled = 1
)

// Stage instance privacy levels
const (
	StagePrivacyGuildOnly = 2
)

// Application command types
const (
	ApplicationCommandTypeChatInput = 1
)

// Application command option types
const (
	CommandOptionString  = 3
	CommandOptionInteger = 4
)

// User is a Discord user object.
type User struct {
	ID            Snowflake `json:"id"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	GlobalName    *string   `json:"global_name,omitempty"`
	Avatar        *string   `json:"avatar,omitempty"`
	Bot           bool      `json:"bot,omitempty"`
	System        bool      `json:"system,omitempty"`
	MFAEnabled    bool      `json:"mfa_enabled,omitempty"`
	Banner        *string   `json:"banner,omitempty"`
	AccentColor   *int      `json:"accent_color,omitempty"`
	Locale        string    `json:"locale,omitempty"`
	Verified      bool      `json:"verified,omitempty"`
	Email         *string   `json:"email,omitempty"`
	Flags         int       `json:"flags,omitempty"`
	PremiumType   int       `json:"premium_type,omitempty"`
	PublicFlags   int       `json:"public_flags,omitempty"`
	AvatarDecoration *string `json:"avatar_decoration,omitempty"`
}

// Guild represents a Discord guild (server).
type Guild struct {
	ID                          Snowflake  `json:"id"`
	Name                        string     `json:"name"`
	Icon                        *string    `json:"icon,omitempty"`
	IconHash                    *string    `json:"icon_hash,omitempty"`
	Splash                      *string    `json:"splash,omitempty"`
	DiscoverySplash             *string    `json:"discovery_splash,omitempty"`
	Owner                       bool       `json:"owner,omitempty"`
	OwnerID                     Snowflake  `json:"owner_id"`
	Permissions                 *string    `json:"permissions,omitempty"`
	Region                      *string    `json:"region,omitempty"`
	AFKChannelID                *Snowflake `json:"afk_channel_id,omitempty"`
	AFKTimeout                  int        `json:"afk_timeout"`
	WidgetEnabled               bool       `json:"widget_enabled,omitempty"`
	WidgetChannelID             *Snowflake `json:"widget_channel_id,omitempty"`
	VerificationLevel           int        `json:"verification_level"`
	DefaultMessageNotifications int        `json:"default_message_notifications"`
	ExplicitContentFilter       int        `json:"explicit_content_filter"`
	Roles                       []*Role    `json:"roles,omitempty"`
	Emojis                      []*Emoji   `json:"emojis,omitempty"`
	Features                    []string   `json:"features"`
	MFALevel                    int        `json:"mfa_level"`
	ApplicationID               *Snowflake `json:"application_id,omitempty"`
	SystemChannelID             *Snowflake `json:"system_channel_id,omitempty"`
	SystemChannelFlags          int        `json:"system_channel_flags"`
	RulesChannelID              *Snowflake `json:"rules_channel_id,omitempty"`
	MaxPresences                *int       `json:"max_presences,omitempty"`
	MaxMembers                  int        `json:"max_members,omitempty"`
	VanityURLCode               *string    `json:"vanity_url_code,omitempty"`
	Description                 *string    `json:"description,omitempty"`
	Banner                      *string    `json:"banner,omitempty"`
	PremiumTier                 int        `json:"premium_tier"`
	PremiumSubscriptionCount    int        `json:"premium_subscription_count,omitempty"`
	PreferredLocale             string     `json:"preferred_locale"`
	PublicUpdatesChannelID      *Snowflake `json:"public_updates_channel_id,omitempty"`
	MaxVideoChannelUsers        int        `json:"max_video_channel_users,omitempty"`
	MaxStageVideoChannelUsers   int        `json:"max_stage_video_channel_users,omitempty"`
	ApproximateMemberCount      int        `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    int        `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               *WelcomeScreen `json:"welcome_screen,omitempty"`
	NSFWLevel                   int        `json:"nsfw_level"`
	Stickers                    []*Sticker `json:"stickers,omitempty"`
	PremiumProgressBarEnabled   bool       `json:"premium_progress_bar_enabled"`
	SafetyAlertsChannelID       *Snowflake `json:"safety_alerts_channel_id,omitempty"`
}

// Channel represents a Discord channel.
type Channel struct {
	ID                            Snowflake              `json:"id"`
	Type                          int                    `json:"type"`
	GuildID                       *Snowflake             `json:"guild_id,omitempty"`
	Position                      *int                   `json:"position,omitempty"`
	PermissionOverwrites          []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                          *string                `json:"name,omitempty"`
	Topic                         *string                `json:"topic,omitempty"`
	NSFW                          bool                   `json:"nsfw,omitempty"`
	LastMessageID                 *Snowflake             `json:"last_message_id,omitempty"`
	Bitrate                       *int                   `json:"bitrate,omitempty"`
	UserLimit                     *int                   `json:"user_limit,omitempty"`
	RateLimitPerUser              *int                   `json:"rate_limit_per_user,omitempty"`
	Recipients                    []*User                `json:"recipients,omitempty"`
	Icon                          *string                `json:"icon,omitempty"`
	OwnerID                       *Snowflake             `json:"owner_id,omitempty"`
	ApplicationID                 *Snowflake             `json:"application_id,omitempty"`
	Managed                       bool                   `json:"managed,omitempty"`
	ParentID                      *Snowflake             `json:"parent_id,omitempty"`
	LastPinTimestamp               *time.Time             `json:"last_pin_timestamp,omitempty"`
	RTCRegion                     *string                `json:"rtc_region,omitempty"`
	VideoQualityMode              *int                   `json:"video_quality_mode,omitempty"`
	MessageCount                  *int                   `json:"message_count,omitempty"`
	MemberCount                   *int                   `json:"member_count,omitempty"`
	ThreadMetadata                *ThreadMetadata        `json:"thread_metadata,omitempty"`
	DefaultAutoArchiveDuration    *int                   `json:"default_auto_archive_duration,omitempty"`
	Permissions                   *string                `json:"permissions,omitempty"`
	Flags                         *int                   `json:"flags,omitempty"`
	TotalMessageSent              *int                   `json:"total_message_sent,omitempty"`
	AvailableTags                 []*ForumTag            `json:"available_tags,omitempty"`
	AppliedTags                   []Snowflake            `json:"applied_tags,omitempty"`
	DefaultReactionEmoji          *DefaultReaction       `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser *int                   `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              *int                   `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *int                   `json:"default_forum_layout,omitempty"`
}

// PermissionOverwrite represents a channel permission overwrite.
type PermissionOverwrite struct {
	ID    Snowflake `json:"id"`
	Type  int       `json:"type"` // 0 = role, 1 = member
	Allow string    `json:"allow"`
	Deny  string    `json:"deny"`
}

// ThreadMetadata contains thread-specific channel fields.
type ThreadMetadata struct {
	Archived            bool       `json:"archived"`
	AutoArchiveDuration int        `json:"auto_archive_duration"`
	ArchiveTimestamp     time.Time  `json:"archive_timestamp"`
	Locked              bool       `json:"locked"`
	Invitable           *bool      `json:"invitable,omitempty"`
	CreateTimestamp      *time.Time `json:"create_timestamp,omitempty"`
}

// ForumTag represents a tag available in a forum channel.
type ForumTag struct {
	ID        Snowflake  `json:"id"`
	Name      string     `json:"name"`
	Moderated bool       `json:"moderated"`
	EmojiID   *Snowflake `json:"emoji_id,omitempty"`
	EmojiName *string    `json:"emoji_name,omitempty"`
}

// DefaultReaction is the default reaction emoji for a forum channel.
type DefaultReaction struct {
	EmojiID   *Snowflake `json:"emoji_id,omitempty"`
	EmojiName *string    `json:"emoji_name,omitempty"`
}

// Role represents a Discord role.
type Role struct {
	ID           Snowflake  `json:"id"`
	Name         string     `json:"name"`
	Color        int        `json:"color"`
	Hoist        bool       `json:"hoist"`
	Icon         *string    `json:"icon,omitempty"`
	UnicodeEmoji *string    `json:"unicode_emoji,omitempty"`
	Position     int        `json:"position"`
	Permissions  string     `json:"permissions"`
	Managed      bool       `json:"managed"`
	Mentionable  bool       `json:"mentionable"`
	Tags         *RoleTags  `json:"tags,omitempty"`
	Flags        int        `json:"flags"`
}

// RoleTags contains special role tag information.
type RoleTags struct {
	BotID                 *Snowflake `json:"bot_id,omitempty"`
	IntegrationID         *Snowflake `json:"integration_id,omitempty"`
	PremiumSubscriber     *bool      `json:"premium_subscriber,omitempty"`
	SubscriptionListingID *Snowflake `json:"subscription_listing_id,omitempty"`
	AvailableForPurchase  *bool      `json:"available_for_purchase,omitempty"`
	GuildConnections      *bool      `json:"guild_connections,omitempty"`
}

// Member represents a guild member.
type Member struct {
	User                       *User       `json:"user,omitempty"`
	Nick                       *string     `json:"nick,omitempty"`
	Avatar                     *string     `json:"avatar,omitempty"`
	Roles                      []Snowflake `json:"roles"`
	JoinedAt                   time.Time   `json:"joined_at"`
	PremiumSince               *time.Time  `json:"premium_since,omitempty"`
	Deaf                       bool        `json:"deaf"`
	Mute                       bool        `json:"mute"`
	Flags                      int         `json:"flags"`
	Pending                    bool        `json:"pending,omitempty"`
	Permissions                *string     `json:"permissions,omitempty"`
	CommunicationDisabledUntil *time.Time  `json:"communication_disabled_until,omitempty"`
}

// Ban represents a guild ban.
type Ban struct {
	Reason *string `json:"reason,omitempty"`
	User   *User   `json:"user"`
}

// Emoji represents a Discord emoji.
type Emoji struct {
	ID            *Snowflake  `json:"id,omitempty"`
	Name          *string     `json:"name,omitempty"`
	Roles         []Snowflake `json:"roles,omitempty"`
	User          *User       `json:"user,omitempty"`
	RequireColons bool        `json:"require_colons,omitempty"`
	Managed       bool        `json:"managed,omitempty"`
	Animated      bool        `json:"animated,omitempty"`
	Available     bool        `json:"available,omitempty"`
}

// Sticker represents a Discord sticker.
type Sticker struct {
	ID          Snowflake  `json:"id"`
	PackID      *Snowflake `json:"pack_id,omitempty"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Tags        string     `json:"tags"`
	Type        int        `json:"type"`
	FormatType  int        `json:"format_type"`
	Available   bool       `json:"available,omitempty"`
	GuildID     *Snowflake `json:"guild_id,omitempty"`
	User        *User      `json:"user,omitempty"`
	SortValue   *int       `json:"sort_value,omitempty"`
}

// Webhook represents a Discord webhook.
type Webhook struct {
	ID            Snowflake  `json:"id"`
	Type          int        `json:"type"`
	GuildID       *Snowflake `json:"guild_id,omitempty"`
	ChannelID     *Snowflake `json:"channel_id,omitempty"`
	User          *User      `json:"user,omitempty"`
	Name          *string    `json:"name,omitempty"`
	Avatar        *string    `json:"avatar,omitempty"`
	Token         *string    `json:"token,omitempty"`
	ApplicationID *Snowflake `json:"application_id,omitempty"`
	URL           *string    `json:"url,omitempty"`
}

// Invite represents a Discord invite.
type Invite struct {
	Code                     string             `json:"code"`
	Guild                    *Guild             `json:"guild,omitempty"`
	Channel                  *Channel           `json:"channel,omitempty"`
	Inviter                  *User              `json:"inviter,omitempty"`
	TargetType               *int               `json:"target_type,omitempty"`
	TargetUser               *User              `json:"target_user,omitempty"`
	TargetApplication        *Application       `json:"target_application,omitempty"`
	ApproximatePresenceCount *int               `json:"approximate_presence_count,omitempty"`
	ApproximateMemberCount   *int               `json:"approximate_member_count,omitempty"`
	ExpiresAt                *time.Time         `json:"expires_at,omitempty"`
	GuildScheduledEvent      *GuildScheduledEvent `json:"guild_scheduled_event,omitempty"`
	Uses                     int                `json:"uses,omitempty"`
	MaxUses                  int                `json:"max_uses,omitempty"`
	MaxAge                   int                `json:"max_age,omitempty"`
	Temporary                bool               `json:"temporary,omitempty"`
	CreatedAt                *time.Time         `json:"created_at,omitempty"`
}

// AutoModerationRule represents an auto-moderation rule.
type AutoModerationRule struct {
	ID              Snowflake          `json:"id"`
	GuildID         Snowflake          `json:"guild_id"`
	Name            string             `json:"name"`
	CreatorID       Snowflake          `json:"creator_id"`
	EventType       int                `json:"event_type"`
	TriggerType     int                `json:"trigger_type"`
	TriggerMetadata *TriggerMetadata   `json:"trigger_metadata"`
	Actions         []*AutoModAction   `json:"actions"`
	Enabled         bool               `json:"enabled"`
	ExemptRoles     []Snowflake        `json:"exempt_roles"`
	ExemptChannels  []Snowflake        `json:"exempt_channels"`
}

// TriggerMetadata contains additional data for auto-mod triggers.
type TriggerMetadata struct {
	KeywordFilter                []string `json:"keyword_filter,omitempty"`
	RegexPatterns                []string `json:"regex_patterns,omitempty"`
	Presets                      []int    `json:"presets,omitempty"`
	AllowList                    []string `json:"allow_list,omitempty"`
	MentionTotalLimit            *int     `json:"mention_total_limit,omitempty"`
	MentionRaidProtectionEnabled bool     `json:"mention_raid_protection_enabled,omitempty"`
}

// AutoModAction represents an action taken by an auto-mod rule.
type AutoModAction struct {
	Type     int                `json:"type"`
	Metadata *AutoModActionMeta `json:"metadata,omitempty"`
}

// AutoModActionMeta contains metadata for auto-mod actions.
type AutoModActionMeta struct {
	ChannelID       *Snowflake `json:"channel_id,omitempty"`
	DurationSeconds *int       `json:"duration_seconds,omitempty"`
	CustomMessage   *string    `json:"custom_message,omitempty"`
}

// GuildScheduledEvent represents a scheduled event in a guild.
type GuildScheduledEvent struct {
	ID                 Snowflake                    `json:"id"`
	GuildID            Snowflake                    `json:"guild_id"`
	ChannelID          *Snowflake                   `json:"channel_id,omitempty"`
	CreatorID          *Snowflake                   `json:"creator_id,omitempty"`
	Name               string                       `json:"name"`
	Description        *string                      `json:"description,omitempty"`
	ScheduledStartTime time.Time                    `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time                   `json:"scheduled_end_time,omitempty"`
	PrivacyLevel       int                          `json:"privacy_level"`
	Status             int                          `json:"status"`
	EntityType         int                          `json:"entity_type"`
	EntityID           *Snowflake                   `json:"entity_id,omitempty"`
	EntityMetadata     *ScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Creator            *User                        `json:"creator,omitempty"`
	UserCount          *int                         `json:"user_count,omitempty"`
	Image              *string                      `json:"image,omitempty"`
}

// ScheduledEventEntityMetadata contains additional metadata for scheduled events.
type ScheduledEventEntityMetadata struct {
	Location *string `json:"location,omitempty"`
}

// StageInstance represents a live stage instance.
type StageInstance struct {
	ID                    Snowflake `json:"id"`
	GuildID               Snowflake `json:"guild_id"`
	ChannelID             Snowflake `json:"channel_id"`
	Topic                 string    `json:"topic"`
	PrivacyLevel          int       `json:"privacy_level"`
	DiscoverableDisabled  bool      `json:"discoverable_disabled"`
	GuildScheduledEventID *Snowflake `json:"guild_scheduled_event_id,omitempty"`
}

// SoundboardSound represents a soundboard sound.
type SoundboardSound struct {
	Name     string     `json:"name"`
	SoundID  Snowflake  `json:"sound_id"`
	Volume   float64    `json:"volume"`
	EmojiID  *Snowflake `json:"emoji_id,omitempty"`
	EmojiName *string   `json:"emoji_name,omitempty"`
	GuildID  *Snowflake `json:"guild_id,omitempty"`
	Available bool      `json:"available"`
	User     *User      `json:"user,omitempty"`
}

// ApplicationCommand represents an application (slash) command.
type ApplicationCommand struct {
	ID                       Snowflake                   `json:"id"`
	Type                     *int                        `json:"type,omitempty"`
	ApplicationID            Snowflake                   `json:"application_id"`
	GuildID                  *Snowflake                  `json:"guild_id,omitempty"`
	Name                     string                      `json:"name"`
	NameLocalizations        map[string]string           `json:"name_localizations,omitempty"`
	Description              string                      `json:"description"`
	DescriptionLocalizations map[string]string           `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                     `json:"default_member_permissions,omitempty"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"`
	NSFW                     *bool                       `json:"nsfw,omitempty"`
	Version                  Snowflake                   `json:"version"`
}

// ApplicationCommandOption represents a command option.
type ApplicationCommandOption struct {
	Type                     int                            `json:"type"`
	Name                     string                         `json:"name"`
	NameLocalizations        map[string]string              `json:"name_localizations,omitempty"`
	Description              string                         `json:"description"`
	DescriptionLocalizations map[string]string              `json:"description_localizations,omitempty"`
	Required                 bool                           `json:"required,omitempty"`
	Choices                  []*ApplicationCommandChoice    `json:"choices,omitempty"`
	Options                  []*ApplicationCommandOption    `json:"options,omitempty"`
	ChannelTypes             []int                          `json:"channel_types,omitempty"`
	MinValue                 *float64                       `json:"min_value,omitempty"`
	MaxValue                 *float64                       `json:"max_value,omitempty"`
	MinLength                *int                           `json:"min_length,omitempty"`
	MaxLength                *int                           `json:"max_length,omitempty"`
	Autocomplete             *bool                          `json:"autocomplete,omitempty"`
}

// ApplicationCommandChoice represents a choice for a command option.
type ApplicationCommandChoice struct {
	Name              string            `json:"name"`
	NameLocalizations map[string]string `json:"name_localizations,omitempty"`
	Value             interface{}       `json:"value"`
}

// GuildTemplate represents a guild template.
type GuildTemplate struct {
	Code                  string     `json:"code"`
	Name                  string     `json:"name"`
	Description           *string    `json:"description,omitempty"`
	UsageCount            int        `json:"usage_count"`
	CreatorID             Snowflake  `json:"creator_id"`
	Creator               *User      `json:"creator,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	SourceGuildID         Snowflake  `json:"source_guild_id"`
	SerializedSourceGuild *Guild     `json:"serialized_source_guild"`
	IsDirty               *bool      `json:"is_dirty,omitempty"`
}

// WelcomeScreen represents a guild's welcome screen.
type WelcomeScreen struct {
	Description     *string                `json:"description,omitempty"`
	WelcomeChannels []*WelcomeScreenChannel `json:"welcome_channels"`
}

// WelcomeScreenChannel represents a channel shown in the welcome screen.
type WelcomeScreenChannel struct {
	ChannelID   Snowflake  `json:"channel_id"`
	Description string     `json:"description"`
	EmojiID     *Snowflake `json:"emoji_id,omitempty"`
	EmojiName   *string    `json:"emoji_name,omitempty"`
}

// GuildWidget represents a guild's widget settings.
type GuildWidget struct {
	Enabled   bool       `json:"enabled"`
	ChannelID *Snowflake `json:"channel_id,omitempty"`
}

// GuildOnboarding represents guild onboarding configuration.
type GuildOnboarding struct {
	GuildID           Snowflake           `json:"guild_id"`
	Prompts           []*OnboardingPrompt `json:"prompts"`
	DefaultChannelIDs []Snowflake         `json:"default_channel_ids"`
	Enabled           bool                `json:"enabled"`
	Mode              int                 `json:"mode"`
}

// OnboardingPrompt represents an onboarding prompt.
type OnboardingPrompt struct {
	ID           Snowflake                `json:"id"`
	Type         int                      `json:"type"`
	Options      []*OnboardingPromptOption `json:"options"`
	Title        string                   `json:"title"`
	SingleSelect bool                     `json:"single_select"`
	Required     bool                     `json:"required"`
	InOnboarding bool                     `json:"in_onboarding"`
}

// OnboardingPromptOption represents an option in an onboarding prompt.
type OnboardingPromptOption struct {
	ID          Snowflake   `json:"id"`
	ChannelIDs  []Snowflake `json:"channel_ids"`
	RoleIDs     []Snowflake `json:"role_ids"`
	Emoji       *Emoji      `json:"emoji,omitempty"`
	EmojiID     *Snowflake  `json:"emoji_id,omitempty"`
	EmojiName   *string     `json:"emoji_name,omitempty"`
	EmojiAnimated *bool     `json:"emoji_animated,omitempty"`
	Title       string      `json:"title"`
	Description *string     `json:"description,omitempty"`
}

// Application represents a Discord application.
type Application struct {
	ID                             Snowflake  `json:"id"`
	Name                           string     `json:"name"`
	Icon                           *string    `json:"icon,omitempty"`
	Description                    string     `json:"description"`
	RPCOrigins                     []string   `json:"rpc_origins,omitempty"`
	BotPublic                      bool       `json:"bot_public"`
	BotRequireCodeGrant            bool       `json:"bot_require_code_grant"`
	Bot                            *User      `json:"bot,omitempty"`
	TermsOfServiceURL              *string    `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyURL               *string    `json:"privacy_policy_url,omitempty"`
	Owner                          *User      `json:"owner,omitempty"`
	VerifyKey                      string     `json:"verify_key"`
	GuildID                        *Snowflake `json:"guild_id,omitempty"`
	Guild                          *Guild     `json:"guild,omitempty"`
	PrimarySkuID                   *Snowflake `json:"primary_sku_id,omitempty"`
	Slug                           *string    `json:"slug,omitempty"`
	CoverImage                     *string    `json:"cover_image,omitempty"`
	Flags                          *int       `json:"flags,omitempty"`
	ApproximateGuildCount          *int       `json:"approximate_guild_count,omitempty"`
	ApproximateUserInstallCount    *int       `json:"approximate_user_install_count,omitempty"`
	RedirectURIs                   []string   `json:"redirect_uris,omitempty"`
	InteractionsEndpointURL        *string    `json:"interactions_endpoint_url,omitempty"`
	RoleConnectionsVerificationURL *string    `json:"role_connections_verification_url,omitempty"`
	Tags                           []string   `json:"tags,omitempty"`
	CustomInstallURL               *string    `json:"custom_install_url,omitempty"`
}

// Message represents a Discord message.
type Message struct {
	ID              Snowflake    `json:"id"`
	ChannelID       Snowflake    `json:"channel_id"`
	Author          *User        `json:"author,omitempty"`
	Content         string       `json:"content"`
	Timestamp       time.Time    `json:"timestamp"`
	EditedTimestamp  *time.Time   `json:"edited_timestamp,omitempty"`
	TTS             bool         `json:"tts"`
	MentionEveryone bool         `json:"mention_everyone"`
	Mentions        []*User      `json:"mentions,omitempty"`
	MentionRoles    []Snowflake  `json:"mention_roles,omitempty"`
	Attachments     []*Attachment `json:"attachments,omitempty"`
	Embeds          []*Embed     `json:"embeds,omitempty"`
	Pinned          bool         `json:"pinned"`
	Type            int          `json:"type"`
	Flags           *int         `json:"flags,omitempty"`
	GuildID         *Snowflake   `json:"guild_id,omitempty"`
}

// Embed represents a Discord embed.
type Embed struct {
	Title       *string        `json:"title,omitempty"`
	Type        *string        `json:"type,omitempty"`
	Description *string        `json:"description,omitempty"`
	URL         *string        `json:"url,omitempty"`
	Timestamp   *time.Time     `json:"timestamp,omitempty"`
	Color       *int           `json:"color,omitempty"`
	Footer      *EmbedFooter   `json:"footer,omitempty"`
	Image       *EmbedImage    `json:"image,omitempty"`
	Thumbnail   *EmbedImage    `json:"thumbnail,omitempty"`
	Video       *EmbedVideo    `json:"video,omitempty"`
	Provider    *EmbedProvider `json:"provider,omitempty"`
	Author      *EmbedAuthor   `json:"author,omitempty"`
	Fields      []*EmbedField  `json:"fields,omitempty"`
}

// EmbedFooter is the footer of an embed.
type EmbedFooter struct {
	Text         string  `json:"text"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

// EmbedImage is an image in an embed.
type EmbedImage struct {
	URL      *string `json:"url,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
}

// EmbedVideo is a video in an embed.
type EmbedVideo struct {
	URL      *string `json:"url,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
}

// EmbedProvider is a provider in an embed.
type EmbedProvider struct {
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

// EmbedAuthor is the author of an embed.
type EmbedAuthor struct {
	Name         *string `json:"name,omitempty"`
	URL          *string `json:"url,omitempty"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

// EmbedField is a field in an embed.
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline *bool  `json:"inline,omitempty"`
}

// Attachment represents a message attachment.
type Attachment struct {
	ID          Snowflake `json:"id"`
	Filename    string    `json:"filename"`
	Description *string   `json:"description,omitempty"`
	ContentType *string   `json:"content_type,omitempty"`
	Size        int       `json:"size"`
	URL         string    `json:"url"`
	ProxyURL    string    `json:"proxy_url"`
	Height      *int      `json:"height,omitempty"`
	Width       *int      `json:"width,omitempty"`
	Ephemeral   bool      `json:"ephemeral,omitempty"`
	DurationSecs *float64 `json:"duration_secs,omitempty"`
	Waveform    *string   `json:"waveform,omitempty"`
	Flags       *int      `json:"flags,omitempty"`
}

// VoiceRegion represents a voice region.
type VoiceRegion struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Optimal  bool   `json:"optimal"`
	Deprecated bool `json:"deprecated"`
	Custom   bool   `json:"custom"`
}
