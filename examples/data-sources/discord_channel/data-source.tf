# SPDX-License-Identifier: MPL-2.0

# Look up information about a Discord channel
data "discord_channel" "general" {
  id = "123456789012345678" # Replace with your channel ID
}

output "channel_name" {
  value = data.discord_channel.general.name
}

output "channel_type" {
  value = data.discord_channel.general.type
}
