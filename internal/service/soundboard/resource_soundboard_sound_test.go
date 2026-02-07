package soundboard_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSoundboardSound_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccSoundboardSoundConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_soundboard_sound.test", "name", "tf-acc-sound"),
					resource.TestCheckResourceAttr("discord_soundboard_sound.test", "guild_id", guildID),
					resource.TestCheckResourceAttr("discord_soundboard_sound.test", "volume", "0.5"),
					resource.TestCheckResourceAttrSet("discord_soundboard_sound.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_soundboard_sound.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateSoundboardSound("discord_soundboard_sound.test"),
			},
			// Update
			{
				Config: testAccSoundboardSoundConfig_updated(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_soundboard_sound.test", "name", "tf-acc-sound-updated"),
					resource.TestCheckResourceAttr("discord_soundboard_sound.test", "volume", "0.8"),
				),
			},
		},
	})
}

func importStateSoundboardSound(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["guild_id"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccSoundboardSoundConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_soundboard_sound" "test" {
  guild_id = %[1]q
  name     = "tf-acc-sound"
  volume   = 0.5
}
`, guildID)
}

func testAccSoundboardSoundConfig_updated(guildID string) string {
	return fmt.Sprintf(`
resource "discord_soundboard_sound" "test" {
  guild_id = %[1]q
  name     = "tf-acc-sound-updated"
  volume   = 0.8
}
`, guildID)
}
