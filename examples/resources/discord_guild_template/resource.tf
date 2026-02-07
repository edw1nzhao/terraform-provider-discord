# SPDX-License-Identifier: MPL-2.0

# Manage a guild template
resource "discord_guild_template" "example" {
  guild_id    = "123456789012345678" # Replace with your guild ID
  name        = "My Server Template"
  description = "A template for setting up new community servers"
}
