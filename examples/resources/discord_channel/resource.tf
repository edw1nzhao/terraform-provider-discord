# SPDX-License-Identifier: MPL-2.0

locals {
  guild_id = "123456789012345678" # Replace with your guild ID
}

# Manage a text channel
resource "discord_channel" "general" {
  guild_id = local.guild_id
  name     = "general"
  type     = 0 # 0=text, 2=voice, 4=category, 5=announcement, 13=stage, 15=forum, 16=media
  topic    = "General discussion channel"
}

# Manage a voice channel inside a category
resource "discord_channel" "voice" {
  guild_id  = local.guild_id
  name      = "voice-chat"
  type      = 2
  parent_id = discord_channel.category.id
  bitrate   = 64000
}

# Manage a category channel
resource "discord_channel" "category" {
  guild_id = local.guild_id
  name     = "My Category"
  type     = 4
}

# Manage a forum channel
resource "discord_channel" "forum" {
  guild_id             = local.guild_id
  name                 = "help-forum"
  type                 = 15
  topic                = "Ask questions and get help from the community"
  default_sort_order   = 0 # 0=latest_activity, 1=creation_date
  default_forum_layout = 1 # 0=not_set, 1=list_view, 2=gallery_view
}
