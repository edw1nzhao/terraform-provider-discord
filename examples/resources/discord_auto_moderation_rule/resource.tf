# SPDX-License-Identifier: MPL-2.0

locals {
  guild_id = "123456789012345678" # Replace with your guild ID
}

# Manage a keyword-based auto-moderation rule
resource "discord_auto_moderation_rule" "no_bad_words" {
  guild_id = local.guild_id
  name     = "Block Bad Words"

  # Event type: 1=MESSAGE_SEND, 2=MEMBER_UPDATE
  event_type = 1

  # Trigger type: 1=KEYWORD, 3=SPAM, 4=KEYWORD_PRESET, 5=MENTION_SPAM, 6=MEMBER_PROFILE
  trigger_type = 1

  trigger_metadata = {
    keyword_filter = ["badword1", "badword2"]
    regex_patterns = ["b[a@]dw[o0]rd"]
  }

  # Actions to take when the rule is triggered
  actions = [
    {
      # Block the message (type 1)
      type = 1
      metadata = {
        custom_message = "Your message was blocked by AutoMod."
      }
    },
    {
      # Send alert to a channel (type 2)
      type = 2
      metadata = {
        channel_id = "123456789012345678" # Replace with your alert channel ID
      }
    },
  ]

  enabled = true

  # Exempt specific roles and channels from this rule
  exempt_roles    = ["123456789012345678"] # Replace with your exempt role IDs
  exempt_channels = ["123456789012345678"] # Replace with your exempt channel IDs
}

# Auto-moderation rule using keyword presets
resource "discord_auto_moderation_rule" "preset_filter" {
  guild_id     = local.guild_id
  name         = "Preset Content Filter"
  event_type   = 1
  trigger_type = 4 # KEYWORD_PRESET

  trigger_metadata = {
    # Presets: 1=Profanity, 2=Sexual Content, 3=Slurs
    presets    = [1, 2, 3]
    allow_list = ["acceptable_word"]
  }

  actions = [
    {
      type = 1 # BLOCK_MESSAGE
    },
  ]

  enabled = true
}
