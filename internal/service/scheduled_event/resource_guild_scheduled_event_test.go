package scheduled_event_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGuildScheduledEvent_external(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	// Schedule the event 1 hour from now to ensure it is in the future.
	startTime := time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)
	endTime := time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccScheduledEventConfig_external(guildID, startTime, endTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_scheduled_event.test", "name", "tf-acc-event"),
					resource.TestCheckResourceAttr("discord_guild_scheduled_event.test", "guild_id", guildID),
					resource.TestCheckResourceAttr("discord_guild_scheduled_event.test", "entity_type", "3"),
					resource.TestCheckResourceAttr("discord_guild_scheduled_event.test", "entity_metadata_location", "Test Location"),
					resource.TestCheckResourceAttrSet("discord_guild_scheduled_event.test", "id"),
					resource.TestCheckResourceAttrSet("discord_guild_scheduled_event.test", "status"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_guild_scheduled_event.test",
				ImportState:       true,
				ImportStateVerify: true,
				// image is base64 on input but hash on read
				ImportStateVerifyIgnore: []string{"image"},
				ImportStateIdFunc:       importStateScheduledEvent("discord_guild_scheduled_event.test"),
			},
		},
	})
}

func TestAccGuildScheduledEvent_voice(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	startTime := time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScheduledEventConfig_voice(guildID, startTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discord_guild_scheduled_event.voice", "name", "tf-acc-voice-event"),
					resource.TestCheckResourceAttr("discord_guild_scheduled_event.voice", "entity_type", "2"),
					resource.TestCheckResourceAttrSet("discord_guild_scheduled_event.voice", "channel_id"),
				),
			},
		},
	})
}

func importStateScheduledEvent(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["guild_id"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccScheduledEventConfig_external(guildID, startTime, endTime string) string {
	return fmt.Sprintf(`
resource "discord_guild_scheduled_event" "test" {
  guild_id                 = %[1]q
  name                     = "tf-acc-event"
  scheduled_start_time     = %[2]q
  scheduled_end_time       = %[3]q
  entity_type              = 3
  entity_metadata_location = "Test Location"
  privacy_level            = 2
}
`, guildID, startTime, endTime)
}

func testAccScheduledEventConfig_voice(guildID, startTime string) string {
	return fmt.Sprintf(`
resource "discord_channel" "event_voice" {
  guild_id = %[1]q
  name     = "tf-acc-event-voice"
  type     = 2
}

resource "discord_guild_scheduled_event" "voice" {
  guild_id             = %[1]q
  name                 = "tf-acc-voice-event"
  scheduled_start_time = %[2]q
  entity_type          = 2
  channel_id           = discord_channel.event_voice.id
  privacy_level        = 2
}
`, guildID, startTime)
}
