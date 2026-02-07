package channel_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccChannelPermission_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccChannelPermissionConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("discord_channel_permission.test", "id"),
					resource.TestCheckResourceAttr("discord_channel_permission.test", "type", "0"),
					resource.TestCheckResourceAttr("discord_channel_permission.test", "allow", "1024"),
					resource.TestCheckResourceAttr("discord_channel_permission.test", "deny", "2048"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_channel_permission.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccChannelPermissionConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "perm_test" {
  guild_id = %[1]q
  name     = "tf-acc-perm-test"
  type     = 0
}

resource "discord_role" "perm_test" {
  guild_id = %[1]q
  name     = "tf-acc-perm-role"
}

resource "discord_channel_permission" "test" {
  channel_id   = discord_channel.perm_test.id
  overwrite_id = discord_role.perm_test.id
  type         = 0
  allow        = "1024"
  deny         = "2048"
}
`, guildID)
}
