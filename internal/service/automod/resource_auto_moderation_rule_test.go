package automod_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAutoModerationRule_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccAutoModRuleConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_auto_moderation_rule.test", "name", "tf-acc-automod"),
					resource.TestCheckResourceAttr("discord_auto_moderation_rule.test", "guild_id", guildID),
					resource.TestCheckResourceAttr("discord_auto_moderation_rule.test", "event_type", "1"),
					resource.TestCheckResourceAttr("discord_auto_moderation_rule.test", "trigger_type", "1"),
					resource.TestCheckResourceAttr("discord_auto_moderation_rule.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("discord_auto_moderation_rule.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_auto_moderation_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateAutoModRule("discord_auto_moderation_rule.test"),
			},
			// Update
			{
				Config: testAccAutoModRuleConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_auto_moderation_rule.test", "name", "tf-acc-automod-updated"),
					resource.TestCheckResourceAttr("discord_auto_moderation_rule.test", "enabled", "false"),
				),
			},
		},
	})
}

func importStateAutoModRule(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["guild_id"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccAutoModRuleConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_auto_moderation_rule" "test" {
  guild_id     = %[1]q
  name         = "tf-acc-automod"
  event_type   = 1
  trigger_type = 1
  enabled      = true

  trigger_metadata = {
    keyword_filter = ["badword"]
  }

  actions = [
    {
      type = 1
    }
  ]
}
`, guildID)
}

func testAccAutoModRuleConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_auto_moderation_rule" "test" {
  guild_id     = %[1]q
  name         = "tf-acc-automod-updated"
  event_type   = 1
  trigger_type = 1
  enabled      = false

  trigger_metadata = {
    keyword_filter = ["badword", "anotherbadword"]
  }

  actions = [
    {
      type = 1
    }
  ]
}
`, guildID)
}
