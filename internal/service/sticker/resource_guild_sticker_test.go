package sticker_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGuildSticker_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccGuildStickerConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_sticker.test", "name", "tf-acc-sticker"),
					resource.TestCheckResourceAttr("discord_guild_sticker.test", "description", "Acceptance test sticker"),
					resource.TestCheckResourceAttr("discord_guild_sticker.test", "tags", "test"),
					resource.TestCheckResourceAttr("discord_guild_sticker.test", "guild_id", guildID),
					resource.TestCheckResourceAttrSet("discord_guild_sticker.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_guild_sticker.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateGuildSticker("discord_guild_sticker.test"),
			},
			// Update
			{
				Config: testAccGuildStickerConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_sticker.test", "name", "tf-acc-sticker-updated"),
					resource.TestCheckResourceAttr("discord_guild_sticker.test", "description", "Updated sticker description"),
					resource.TestCheckResourceAttr("discord_guild_sticker.test", "tags", "updated"),
				),
			},
		},
	})
}

func importStateGuildSticker(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["guild_id"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccGuildStickerConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_guild_sticker" "test" {
  guild_id    = %[1]q
  name        = "tf-acc-sticker"
  description = "Acceptance test sticker"
  tags        = "test"
}
`, guildID)
}

func testAccGuildStickerConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_guild_sticker" "test" {
  guild_id    = %[1]q
  name        = "tf-acc-sticker-updated"
  description = "Updated sticker description"
  tags        = "updated"
}
`, guildID)
}
