package onboarding_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGuildOnboarding_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccGuildOnboardingConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_onboarding.test", "guild_id", guildID),
					resource.TestCheckResourceAttr("discord_guild_onboarding.test", "enabled", "true"),
					resource.TestCheckResourceAttr("discord_guild_onboarding.test", "mode", "0"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_guild_onboarding.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccGuildOnboardingConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_onboarding.test", "enabled", "true"),
					resource.TestCheckResourceAttr("discord_guild_onboarding.test", "mode", "1"),
				),
			},
		},
	})
}

func testAccGuildOnboardingConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "onboard_ch" {
  guild_id = %[1]q
  name     = "tf-acc-onboard"
  type     = 0
}

resource "discord_guild_onboarding" "test" {
  guild_id            = %[1]q
  enabled             = true
  mode                = 0
  default_channel_ids = [discord_channel.onboard_ch.id]

  prompts = [
    {
      type          = 0
      title         = "What brings you here?"
      single_select = false
      required      = true
      in_onboarding = true

      options = [
        {
          title       = "Chatting"
          channel_ids = [discord_channel.onboard_ch.id]
        }
      ]
    }
  ]
}
`, guildID)
}

func testAccGuildOnboardingConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "onboard_ch" {
  guild_id = %[1]q
  name     = "tf-acc-onboard"
  type     = 0
}

resource "discord_guild_onboarding" "test" {
  guild_id            = %[1]q
  enabled             = true
  mode                = 1
  default_channel_ids = [discord_channel.onboard_ch.id]

  prompts = [
    {
      type          = 0
      title         = "What brings you here?"
      single_select = true
      required      = true
      in_onboarding = true

      options = [
        {
          title       = "Chatting"
          channel_ids = [discord_channel.onboard_ch.id]
        },
        {
          title       = "Gaming"
          channel_ids = [discord_channel.onboard_ch.id]
        }
      ]
    }
  ]
}
`, guildID)
}
