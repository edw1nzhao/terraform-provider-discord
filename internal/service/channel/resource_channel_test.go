package channel_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccChannel_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccChannelConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_channel.test", "name", "tf-acc-test"),
					resource.TestCheckResourceAttr("discord_channel.test", "guild_id", guildID),
					resource.TestCheckResourceAttr("discord_channel.test", "type", "0"),
					resource.TestCheckResourceAttrSet("discord_channel.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_channel.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccChannelConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_channel.test", "name", "tf-acc-test-updated"),
					resource.TestCheckResourceAttr("discord_channel.test", "topic", "Updated topic"),
				),
			},
		},
	})
}

func TestAccChannel_voice(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChannelConfig_voice(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_channel.voice", "name", "tf-acc-voice"),
					resource.TestCheckResourceAttr("discord_channel.voice", "type", "2"),
					resource.TestCheckResourceAttrSet("discord_channel.voice", "bitrate"),
				),
			},
		},
	})
}

func TestAccChannel_category(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChannelConfig_category(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_channel.category", "name", "tf-acc-category"),
					resource.TestCheckResourceAttr("discord_channel.category", "type", "4"),
				),
			},
		},
	})
}

func testAccChannelConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "test" {
  guild_id = %[1]q
  name     = "tf-acc-test"
  type     = 0
}
`, guildID)
}

func testAccChannelConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "test" {
  guild_id = %[1]q
  name     = "tf-acc-test-updated"
  type     = 0
  topic    = "Updated topic"
}
`, guildID)
}

func testAccChannelConfig_voice(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "voice" {
  guild_id = %[1]q
  name     = "tf-acc-voice"
  type     = 2
}
`, guildID)
}

func testAccChannelConfig_category(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "category" {
  guild_id = %[1]q
  name     = "tf-acc-category"
  type     = 4
}
`, guildID)
}
