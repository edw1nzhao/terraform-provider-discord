package member_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMemberRoles_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")
	userID := os.Getenv("DISCORD_USER_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckUser(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccMemberRolesConfig_basic(guildID, userID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_member_roles.test", "guild_id", guildID),
					resource.TestCheckResourceAttr("discord_member_roles.test", "user_id", userID),
					resource.TestCheckResourceAttrSet("discord_member_roles.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_member_roles.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMemberRolesConfig_basic(guildID, userID string) string {
	return fmt.Sprintf(`
resource "discord_role" "member_test" {
  guild_id = %[1]q
  name     = "tf-acc-member-role"
}

resource "discord_member_roles" "test" {
  guild_id = %[1]q
  user_id  = %[2]q
  roles    = [discord_role.member_test.id]
}
`, guildID, userID)
}
