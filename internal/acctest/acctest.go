package acctest

import (
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ProtoV6ProviderFactories returns provider factories for acceptance tests.
var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"discord": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// PreCheck validates that the DISCORD_TOKEN environment variable is set.
func PreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("DISCORD_TOKEN") == "" {
		t.Fatal("DISCORD_TOKEN must be set for acceptance tests")
	}
}

// PreCheckGuild validates guild-specific test prerequisites.
func PreCheckGuild(t *testing.T) {
	t.Helper()
	PreCheck(t)
	if os.Getenv("DISCORD_GUILD_ID") == "" {
		t.Fatal("DISCORD_GUILD_ID must be set for acceptance tests")
	}
}

// PreCheckApplication validates application-specific test prerequisites.
func PreCheckApplication(t *testing.T) {
	t.Helper()
	PreCheck(t)
	if os.Getenv("DISCORD_APPLICATION_ID") == "" {
		t.Fatal("DISCORD_APPLICATION_ID must be set for acceptance tests")
	}
}

// PreCheckUser validates user-specific test prerequisites.
func PreCheckUser(t *testing.T) {
	t.Helper()
	PreCheckGuild(t)
	if os.Getenv("DISCORD_USER_ID") == "" {
		t.Fatal("DISCORD_USER_ID must be set for acceptance tests")
	}
}

// PreCheckBanUser validates ban test prerequisites.
func PreCheckBanUser(t *testing.T) {
	t.Helper()
	PreCheckGuild(t)
	if os.Getenv("DISCORD_BAN_USER_ID") == "" {
		t.Fatal("DISCORD_BAN_USER_ID must be set for acceptance tests")
	}
}
