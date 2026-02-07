# SPDX-License-Identifier: MPL-2.0

# Manage a Discord guild role
resource "discord_role" "moderator" {
  guild_id = "123456789012345678" # Replace with your guild ID
  name     = "Moderator"

  # Use the discord_permission data source to compute permissions
  permissions = data.discord_permission.moderator.allow_bits

  # Color as an integer (use discord_color data source to convert hex)
  color = data.discord_color.blue.int

  # Display role members separately in the sidebar
  hoist = true

  # Allow anyone to @mention this role
  mentionable = true

  # Position in the role hierarchy
  position = 5
}

# Compute moderator permissions
data "discord_permission" "moderator" {
  manage_messages = true
  kick_members    = true
  ban_members     = true
  moderate_members = true
}

# Convert hex color to integer
data "discord_color" "blue" {
  hex = "#3498DB"
}
