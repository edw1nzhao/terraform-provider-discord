package application_command_test

import (
	"fmt"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGlobalApplicationCommand_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckApplication(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccGlobalApplicationCommandConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_global_application_command.test", "name", "tf-acc-test-cmd"),
					resource.TestCheckResourceAttr("discord_global_application_command.test", "description", "Acceptance test command"),
					resource.TestCheckResourceAttr("discord_global_application_command.test", "type", "1"),
					resource.TestCheckResourceAttrSet("discord_global_application_command.test", "id"),
					resource.TestCheckResourceAttrSet("discord_global_application_command.test", "application_id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_global_application_command.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateGlobalAppCommand("discord_global_application_command.test"),
			},
			// Update
			{
				Config: testAccGlobalApplicationCommandConfig_updated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_global_application_command.test", "name", "tf-acc-test-cmd"),
					resource.TestCheckResourceAttr("discord_global_application_command.test", "description", "Updated description"),
				),
			},
		},
	})
}

func importStateGlobalAppCommand(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["application_id"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccGlobalApplicationCommandConfig_basic() string {
	return `
resource "discord_global_application_command" "test" {
  name        = "tf-acc-test-cmd"
  description = "Acceptance test command"
  type        = 1
}
`
}

func testAccGlobalApplicationCommandConfig_updated() string {
	return `
resource "discord_global_application_command" "test" {
  name        = "tf-acc-test-cmd"
  description = "Updated description"
  type        = 1
}
`
}
