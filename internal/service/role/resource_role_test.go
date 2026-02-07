package role_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRole_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccRoleConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_role.test", "name", "tf-acc-test-role"),
					resource.TestCheckResourceAttr("discord_role.test", "guild_id", guildID),
					resource.TestCheckResourceAttrSet("discord_role.test", "id"),
					resource.TestCheckResourceAttrSet("discord_role.test", "permissions"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccRoleConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_role.test", "name", "tf-acc-test-role-updated"),
					resource.TestCheckResourceAttr("discord_role.test", "color", "16711680"),
					resource.TestCheckResourceAttr("discord_role.test", "hoist", "true"),
					resource.TestCheckResourceAttr("discord_role.test", "mentionable", "true"),
				),
			},
		},
	})
}

func testAccRoleConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_role" "test" {
  guild_id = %[1]q
  name     = "tf-acc-test-role"
}
`, guildID)
}

func testAccRoleConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_role" "test" {
  guild_id    = %[1]q
  name        = "tf-acc-test-role-updated"
  color       = 16711680
  hoist       = true
  mentionable = true
}
`, guildID)
}
