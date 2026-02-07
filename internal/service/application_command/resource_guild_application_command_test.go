package application_command_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGuildApplicationCommand_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckApplication(t)
			acctest.PreCheckGuild(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccGuildApplicationCommandConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_application_command.test", "name", "tf-acc-guild-cmd"),
					resource.TestCheckResourceAttr("discord_guild_application_command.test", "description", "Guild-scoped test command"),
					resource.TestCheckResourceAttr("discord_guild_application_command.test", "guild_id", guildID),
					resource.TestCheckResourceAttrSet("discord_guild_application_command.test", "id"),
					resource.TestCheckResourceAttrSet("discord_guild_application_command.test", "application_id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_guild_application_command.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateGuildAppCommand("discord_guild_application_command.test"),
			},
			// Update
			{
				Config: testAccGuildApplicationCommandConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_application_command.test", "description", "Updated guild command"),
				),
			},
		},
	})
}

func importStateGuildAppCommand(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s",
			rs.Primary.Attributes["application_id"],
			rs.Primary.Attributes["guild_id"],
			rs.Primary.Attributes["id"],
		), nil
	}
}

func testAccGuildApplicationCommandConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_guild_application_command" "test" {
  guild_id    = %[1]q
  name        = "tf-acc-guild-cmd"
  description = "Guild-scoped test command"
  type        = 1
}
`, guildID)
}

func testAccGuildApplicationCommandConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_guild_application_command" "test" {
  guild_id    = %[1]q
  name        = "tf-acc-guild-cmd"
  description = "Updated guild command"
  type        = 1
}
`, guildID)
}
