package channel_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccChannelDataSource_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChannelDataSourceConfig(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.discord_channel.test", "id"),
					resource.TestCheckResourceAttrSet("data.discord_channel.test", "name"),
					resource.TestCheckResourceAttrSet("data.discord_channel.test", "type"),
				),
			},
		},
	})
}

func testAccChannelDataSourceConfig(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "ds_test" {
  guild_id = %[1]q
  name     = "tf-acc-ds-test"
  type     = 0
}

data "discord_channel" "test" {
  id = discord_channel.ds_test.id
}
`, guildID)
}
