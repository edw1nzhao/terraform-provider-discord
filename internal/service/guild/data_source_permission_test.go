package guild_test

import (
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPermissionDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.discord_permission.test", "allow_bits", "1024"),
					resource.TestCheckResourceAttr("data.discord_permission.test", "deny_bits", "2048"),
				),
			},
		},
	})
}

func testAccPermissionDataSourceConfig() string {
	return `
data "discord_permission" "test" {
  view_channel  = true
  send_messages = false
}
`
}
