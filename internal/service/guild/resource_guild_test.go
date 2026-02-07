package guild_test

import (
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGuild_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccGuildConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild.test", "name", "tf-acc-test-guild"),
					resource.TestCheckResourceAttrSet("discord_guild.test", "id"),
					resource.TestCheckResourceAttrSet("discord_guild.test", "owner_id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_guild.test",
				ImportState:       true,
				ImportStateVerify: true,
				// icon is base64 on input but hash on read
				ImportStateVerifyIgnore: []string{"icon"},
			},
			// Update
			{
				Config: testAccGuildConfig_updated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild.test", "name", "tf-acc-test-guild-updated"),
				),
			},
		},
	})
}

func testAccGuildConfig_basic() string {
	return `
resource "discord_guild" "test" {
  name = "tf-acc-test-guild"
}
`
}

func testAccGuildConfig_updated() string {
	return `
resource "discord_guild" "test" {
  name = "tf-acc-test-guild-updated"
}
`
}
