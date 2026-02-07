package application_test

import (
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplication_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckApplication(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create (PATCH) and Read
			{
				Config: testAccApplicationConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("discord_application.test", "id"),
					resource.TestCheckResourceAttrSet("discord_application.test", "name"),
					resource.TestCheckResourceAttr("discord_application.test", "description", "Terraform acceptance test application"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_application.test",
				ImportState:       true,
				ImportStateVerify: true,
				// icon and cover_image are base64 on input but hash on read
				ImportStateVerifyIgnore: []string{"icon", "cover_image"},
			},
			// Update
			{
				Config: testAccApplicationConfig_updated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_application.test", "description", "Updated acceptance test description"),
				),
			},
		},
	})
}

func testAccApplicationConfig_basic() string {
	return `
resource "discord_application" "test" {
  description = "Terraform acceptance test application"
}
`
}

func testAccApplicationConfig_updated() string {
	return `
resource "discord_application" "test" {
  description = "Updated acceptance test description"
}
`
}
