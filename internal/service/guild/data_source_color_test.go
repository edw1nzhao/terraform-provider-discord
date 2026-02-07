package guild_test

import (
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccColorDataSource_hex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccColorDataSourceConfig_hex(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// #FF0000 = 16711680
					resource.TestCheckResourceAttr("data.discord_color.test", "int", "16711680"),
				),
			},
		},
	})
}

func TestAccColorDataSource_rgb(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccColorDataSourceConfig_rgb(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// (0 << 16) | (255 << 8) | 0 = 65280
					resource.TestCheckResourceAttr("data.discord_color.test_rgb", "int", "65280"),
				),
			},
		},
	})
}

func testAccColorDataSourceConfig_hex() string {
	return `
data "discord_color" "test" {
  hex = "#FF0000"
}
`
}

func testAccColorDataSourceConfig_rgb() string {
	return `
data "discord_color" "test_rgb" {
  rgb = {
    r = 0
    g = 255
    b = 0
  }
}
`
}
