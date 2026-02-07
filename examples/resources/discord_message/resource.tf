# SPDX-License-Identifier: MPL-2.0

locals {
  channel_id = "123456789012345678" # Replace with your channel ID
}

# Send a simple text message to a channel
resource "discord_message" "welcome" {
  channel_id = local.channel_id
  content    = "Welcome to the server! Please read the rules."
}

# Send a message with an embed
resource "discord_message" "announcement" {
  channel_id = local.channel_id
  content    = "Check out our latest update!"

  embed {
    title       = "Server Update v2.0"
    description = "We have added new features to the server."
    url         = "https://example.com/updates"
    color       = 3447003 # Blue color
    footer_text = "Posted by the bot"
  }
}
