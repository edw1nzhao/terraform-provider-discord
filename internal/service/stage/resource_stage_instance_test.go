package stage_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccStageInstance_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccStageInstanceConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_stage_instance.test", "topic", "Terraform Acceptance Test"),
					resource.TestCheckResourceAttr("discord_stage_instance.test", "privacy_level", "2"),
					resource.TestCheckResourceAttrSet("discord_stage_instance.test", "id"),
					resource.TestCheckResourceAttrSet("discord_stage_instance.test", "guild_id"),
				),
			},
			// ImportState - stage instances import by channel_id
			{
				ResourceName:      "discord_stage_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateStageByChannelID("discord_stage_instance.test"),
			},
			// Update
			{
				Config: testAccStageInstanceConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_stage_instance.test", "topic", "Updated Stage Topic"),
				),
			},
		},
	})
}

func importStateStageByChannelID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return rs.Primary.Attributes["channel_id"], nil
	}
}

func testAccStageInstanceConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "stage_ch" {
  guild_id = %[1]q
  name     = "tf-acc-stage"
  type     = 13
}

resource "discord_stage_instance" "test" {
  channel_id    = discord_channel.stage_ch.id
  topic         = "Terraform Acceptance Test"
  privacy_level = 2
}
`, guildID)
}

func testAccStageInstanceConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "stage_ch" {
  guild_id = %[1]q
  name     = "tf-acc-stage"
  type     = 13
}

resource "discord_stage_instance" "test" {
  channel_id    = discord_channel.stage_ch.id
  topic         = "Updated Stage Topic"
  privacy_level = 2
}
`, guildID)
}
