package ban_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBan_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")
	userID := os.Getenv("DISCORD_BAN_USER_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckBanUser(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccBanConfig_basic(guildID, userID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_ban.test", "guild_id", guildID),
					resource.TestCheckResourceAttr("discord_ban.test", "user_id", userID),
					resource.TestCheckResourceAttrSet("discord_ban.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_ban.test",
				ImportState:       true,
				ImportStateVerify: true,
				// delete_message_seconds is write-only
				ImportStateVerifyIgnore: []string{"delete_message_seconds"},
				ImportStateIdFunc:       importStateBan("discord_ban.test"),
			},
		},
	})
}

func importStateBan(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["guild_id"], rs.Primary.Attributes["user_id"]), nil
	}
}

func testAccBanConfig_basic(guildID, userID string) string {
	return fmt.Sprintf(`
resource "discord_ban" "test" {
  guild_id = %[1]q
  user_id  = %[2]q
}
`, guildID, userID)
}
