package invite_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/edw1nzhao/terraform-provider-discord/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInvite_basic(t *testing.T) {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckGuild(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccInviteConfig_basic(guildID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("discord_invite.test", "id"),
					resource.TestCheckResourceAttrSet("discord_invite.test", "url"),
					resource.TestCheckResourceAttr("discord_invite.test", "max_age", "3600"),
					resource.TestCheckResourceAttr("discord_invite.test", "max_uses", "10"),
					resource.TestCheckResourceAttr("discord_invite.test", "temporary", "false"),
					resource.TestCheckResourceAttr("discord_invite.test", "unique", "true"),
				),
			},
			// ImportState
			{
				ResourceName:      "discord_invite.test",
				ImportState:       true,
				ImportStateVerify: true,
				// unique is write-only, not returned by API
				ImportStateVerifyIgnore: []string{"unique"},
			},
		},
	})
}

func testAccInviteConfig_basic(guildID string) string {
	return fmt.Sprintf(`
resource "discord_channel" "invite_test" {
  guild_id = %[1]q
  name     = "tf-acc-invite-ch"
  type     = 0
}

resource "discord_invite" "test" {
  channel_id = discord_channel.invite_test.id
  max_age    = 3600
  max_uses   = 10
  temporary  = false
  unique     = true
}
`, guildID)
}
