package guild_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGuildDataSource_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGuildDataSourceConfig(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.discord_guild.test", "id", guildID),
					resource.TestCheckResourceAttrSet("data.discord_guild.test", "name"),
					resource.TestCheckResourceAttrSet("data.discord_guild.test", "owner_id"),
				),
			},
		},
	})
}

func testAccGuildDataSourceConfig(guildID string) string {
	return fmt.Sprintf(`
data "discord_guild" "test" {
  id = %[1]q
}
`, guildID)
}
