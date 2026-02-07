# SPDX-License-Identifier: MPL-2.0

# Manage guild widget settings
# Note: Deleting this resource disables the widget.
resource "discord_guild_widget" "example" {
  guild_id = "123456789012345678" # Replace with your guild ID
  enabled  = true

  # Channel that the widget will generate an invite to
  channel_id = "123456789012345678" # Replace with your channel ID
}
