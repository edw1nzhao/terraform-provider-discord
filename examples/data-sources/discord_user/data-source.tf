# SPDX-License-Identifier: MPL-2.0

# Look up information about a Discord user
data "discord_user" "example" {
  id = "123456789012345678" # Replace with your user ID
}

output "username" {
  value = data.discord_user.example.username
}

output "is_bot" {
  value = data.discord_user.example.bot
}
