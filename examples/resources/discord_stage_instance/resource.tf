# SPDX-License-Identifier: MPL-2.0

# Manage a stage instance for a stage channel
resource "discord_stage_instance" "example" {
  channel_id = "123456789012345678" # Replace with your stage channel ID (type 13)
  topic      = "Community Q&A Session"

  # Privacy level: 2 = GUILD_ONLY (default)
  privacy_level = 2
}
