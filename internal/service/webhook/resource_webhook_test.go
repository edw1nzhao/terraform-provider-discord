package webhook_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebhook_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccWebhookConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_webhook.test", "name", "tf-acc-webhook"),
					resource.TestCheckResourceAttrSet("discord_webhook.test", "id"),
					resource.TestCheckResourceAttrSet("discord_webhook.test", "token"),
					resource.TestCheckResourceAttrSet("discord_webhook.test", "url"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_webhook.test",
				ImportState:       true,
				ImportStateVerify: true,
				// avatar is base64 on input but hash on read
				ImportStateVerifyIgnore: []string{"avatar"},
			},
			// Update
			{
				Config: testAccWebhookConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_webhook.test", "name", "tf-acc-webhook-updated"),
				),
			},
		},
	})
}

func testAccWebhookConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "webhook_test" {
  guild_id = %[1]q
  name     = "tf-acc-webhook-ch"
  type     = 0
}

resource "discord_webhook" "test" {
  channel_id = discord_channel.webhook_test.id
  name       = "tf-acc-webhook"
}
`, guildID)
}

func testAccWebhookConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "webhook_test" {
  guild_id = %[1]q
  name     = "tf-acc-webhook-ch"
  type     = 0
}

resource "discord_webhook" "test" {
  channel_id = discord_channel.webhook_test.id
  name       = "tf-acc-webhook-updated"
}
`, guildID)
}
