# SPDX-License-Identifier: MPL-2.0

locals {
  guild_id = "123456789012345678" # Replace with your guild ID
}

# Manage an external scheduled event (e.g., a meetup)
resource "discord_guild_scheduled_event" "meetup" {
  guild_id = local.guild_id
  name     = "Community Meetup"

  description        = "Join us for our monthly community meetup!"
  scheduled_start_time = "2025-12-01T18:00:00Z"
  scheduled_end_time   = "2025-12-01T20:00:00Z"

  # Entity type: 1=STAGE_INSTANCE, 2=VOICE, 3=EXTERNAL
  entity_type = 3

  # Required for EXTERNAL events
  entity_metadata_location = "Discord Stage Channel"

  # Privacy level: 2=GUILD_ONLY (default)
  privacy_level = 2
}

# Manage a voice channel event
resource "discord_guild_scheduled_event" "game_night" {
  guild_id = local.guild_id
  name     = "Game Night"

  description          = "Weekly game night in voice chat"
  scheduled_start_time = "2025-12-05T20:00:00Z"
  entity_type          = 2                        # VOICE
  channel_id           = "123456789012345678"      # Replace with your voice channel ID
}
