package conns

import "github.com/edw1nzhao/terraform-provider-discord/internal/discord"

// ProviderData holds the configured Discord client and associated metadata.
// It is passed to resources and data sources via Configure.
type ProviderData struct {
	Client        *discord.Client
	ApplicationID string
}
