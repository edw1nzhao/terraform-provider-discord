package template_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGuildTemplate_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccGuildTemplateConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_template.test", "name", "tf-acc-template"),
					resource.TestCheckResourceAttr("discord_guild_template.test", "guild_id", guildID),
					resource.TestCheckResourceAttrSet("discord_guild_template.test", "id"),
					resource.TestCheckResourceAttrSet("discord_guild_template.test", "source_guild_id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_guild_template.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateGuildTemplate("discord_guild_template.test"),
			},
			// Update
			{
				Config: testAccGuildTemplateConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_template.test", "name", "tf-acc-template-updated"),
					resource.TestCheckResourceAttr("discord_guild_template.test", "description", "Updated description"),
				),
			},
		},
	})
}

func importStateGuildTemplate(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["guild_id"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccGuildTemplateConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_guild_template" "test" {
  guild_id = %[1]q
  name     = "tf-acc-template"
}
`, guildID)
}

func testAccGuildTemplateConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_guild_template" "test" {
  guild_id    = %[1]q
  name        = "tf-acc-template-updated"
  description = "Updated description"
}
`, guildID)
}
