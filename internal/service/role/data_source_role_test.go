package role_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleDataSource_byName(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleDataSourceConfig_byName(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.discord_role.test", "id"),
					resource.TestCheckResourceAttr("data.discord_role.test", "name", "tf-acc-ds-role"),
					resource.TestCheckResourceAttrSet("data.discord_role.test", "permissions"),
				),
			},
		},
	})
}

func testAccRoleDataSourceConfig_byName(guildID string) string {
	return fmt.Sprintf(`
resource "discord_role" "ds_test" {
  guild_id = %[1]q
  name     = "tf-acc-ds-role"
}

data "discord_role" "test" {
  guild_id = %[1]q
  name     = discord_role.ds_test.name
}
`, guildID)
}
