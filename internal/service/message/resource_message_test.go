package message_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMessage_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMessageConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_message.test", "content", "Hello from Terraform acceptance test"),
					resource.TestCheckResourceAttrSet("discord_message.test", "id"),
					resource.TestCheckResourceAttrSet("discord_message.test", "channel_id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_message.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateMessage("discord_message.test"),
			},
			// Update
			{
				Config: testAccMessageConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_message.test", "content", "Updated message content"),
				),
			},
		},
	})
}

func TestAccMessage_embed(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMessageConfig_embed(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("discord_message.embed", "id"),
					resource.TestCheckResourceAttr("discord_message.embed", "embed.0.title", "Test Embed"),
					resource.TestCheckResourceAttr("discord_message.embed", "embed.0.description", "Test description"),
				),
			},
		},
	})
}

func importStateMessage(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["channel_id"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccMessageConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "msg_test" {
  guild_id = %[1]q
  name     = "tf-acc-msg-test"
  type     = 0
}

resource "discord_message" "test" {
  channel_id = discord_channel.msg_test.id
  content    = "Hello from Terraform acceptance test"
}
`, guildID)
}

func testAccMessageConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "msg_test" {
  guild_id = %[1]q
  name     = "tf-acc-msg-test"
  type     = 0
}

resource "discord_message" "test" {
  channel_id = discord_channel.msg_test.id
  content    = "Updated message content"
}
`, guildID)
}

func testAccMessageConfig_embed(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "embed_test" {
  guild_id = %[1]q
  name     = "tf-acc-embed-test"
  type     = 0
}

resource "discord_message" "embed" {
  channel_id = discord_channel.embed_test.id
  content    = "Message with embed"

  embed {
    title       = "Test Embed"
    description = "Test description"
    color       = 16711680
  }
}
`, guildID)
}
