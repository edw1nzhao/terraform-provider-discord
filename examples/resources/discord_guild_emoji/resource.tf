# SPDX-License-Identifier: MPL-2.0

# Manage a custom guild emoji
resource "discord_guild_emoji" "example" {
  guild_id = "123456789012345678" # Replace with your guild ID
  name     = "my_emoji"

  # Base64-encoded image data URI (required on create, forces replacement if changed)
  image = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUg..."

  # Optional: restrict emoji usage to specific roles
  roles = [
    discord_role.member.id,
  ]
}
