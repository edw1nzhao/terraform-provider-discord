# SPDX-License-Identifier: MPL-2.0

# Look up information about a Discord guild
data "discord_guild" "example" {
  id = "123456789012345678" # Replace with your guild ID
}

output "guild_name" {
  value = data.discord_guild.example.name
}

output "guild_owner_id" {
  value = data.discord_guild.example.owner_id
}
