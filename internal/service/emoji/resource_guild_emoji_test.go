package emoji_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGuildEmoji_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	// A minimal 1x1 white PNG as base64 data URI for testing
	testImage := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccGuildEmojiConfig_basic(guildID, testImage),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_emoji.test", "name", "tfacctest"),
					resource.TestCheckResourceAttr("discord_guild_emoji.test", "guild_id", guildID),
					resource.TestCheckResourceAttrSet("discord_guild_emoji.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_guild_emoji.test",
				ImportState:       true,
				ImportStateVerify: true,
				// image is write-only (only used on create)
				ImportStateVerifyIgnore: []string{"image"},
				ImportStateIdFunc:       importStateGuildEmoji("discord_guild_emoji.test"),
			},
			// Update name
			{
				Config: testAccGuildEmojiConfig_updated(guildID, testImage),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_emoji.test", "name", "tfacctestupdated"),
				),
			},
		},
	})
}

func importStateGuildEmoji(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["guild_id"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccGuildEmojiConfig_basic(guildID, image string) string {
	return fmt.Sprintf(`
resource "discord_guild_emoji" "test" {
  guild_id = %[1]q
  name     = "tfacctest"
  image    = %[2]q
}
`, guildID, image)
}

func testAccGuildEmojiConfig_updated(guildID, image string) string {
	return fmt.Sprintf(`
resource "discord_guild_emoji" "test" {
  guild_id = %[1]q
  name     = "tfacctestupdated"
  image    = %[2]q
}
`, guildID, image)
}
