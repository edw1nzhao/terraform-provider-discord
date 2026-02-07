# SPDX-License-Identifier: MPL-2.0

# Manage a guild soundboard sound
resource "discord_soundboard_sound" "example" {
  guild_id = "123456789012345678" # Replace with your guild ID
  name     = "airhorn"

  # Volume from 0.0 to 1.0 (default: 1.0)
  volume = 0.75

  # Optional emoji association
  emoji_name = "\ud83d\udce3"
}
