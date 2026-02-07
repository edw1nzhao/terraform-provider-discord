package widget_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGuildWidget_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccGuildWidgetConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_widget.test", "guild_id", guildID),
					resource.TestCheckResourceAttr("discord_guild_widget.test", "enabled", "true"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_guild_widget.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - disable
			{
				Config: testAccGuildWidgetConfig_disabled(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_widget.test", "enabled", "false"),
				),
			},
		},
	})
}

func testAccGuildWidgetConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_guild_widget" "test" {
  guild_id = %[1]q
  enabled  = true
}
`, guildID)
}

func testAccGuildWidgetConfig_disabled(guildID string) string {
	return fmt.Sprintf(`
resource "discord_guild_widget" "test" {
  guild_id = %[1]q
  enabled  = false
}
`, guildID)
}
