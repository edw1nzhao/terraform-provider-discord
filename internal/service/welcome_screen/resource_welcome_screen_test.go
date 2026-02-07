package welcome_screen_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWelcomeScreen_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccWelcomeScreenConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_welcome_screen.test", "guild_id", guildID),
					resource.TestCheckResourceAttr("discord_welcome_screen.test", "enabled", "true"),
					resource.TestCheckResourceAttr("discord_welcome_screen.test", "description", "Welcome to the server!"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_welcome_screen.test",
				ImportState:       true,
				ImportStateVerify: true,
				// enabled is not always returned by API on read
				ImportStateVerifyIgnore: []string{"enabled"},
			},
			// Update
			{
				Config: testAccWelcomeScreenConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_welcome_screen.test", "description", "Updated welcome description"),
				),
			},
		},
	})
}

func testAccWelcomeScreenConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "welcome_ch" {
  guild_id = %[1]q
  name     = "tf-acc-welcome"
  type     = 0
}

resource "discord_welcome_screen" "test" {
  guild_id    = %[1]q
  enabled     = true
  description = "Welcome to the server!"

  welcome_channels = [
    {
      channel_id  = discord_channel.welcome_ch.id
      description = "General chat channel"
    }
  ]
}
`, guildID)
}

func testAccWelcomeScreenConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "welcome_ch" {
  guild_id = %[1]q
  name     = "tf-acc-welcome"
  type     = 0
}

resource "discord_welcome_screen" "test" {
  guild_id    = %[1]q
  enabled     = true
  description = "Updated welcome description"

  welcome_channels = [
    {
      channel_id  = discord_channel.welcome_ch.id
      description = "Updated channel description"
    }
  ]
}
`, guildID)
}
