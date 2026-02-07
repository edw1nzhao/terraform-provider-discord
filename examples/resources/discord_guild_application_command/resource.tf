# SPDX-License-Identifier: MPL-2.0

# Manage a guild-scoped slash command (only available in one guild)
resource "discord_guild_application_command" "config" {
  guild_id    = "123456789012345678" # Replace with your guild ID
  name        = "config"
  description = "Server configuration commands"

  # Command type: 1=CHAT_INPUT (default), 2=USER, 3=MESSAGE
  type = 1

  # Require specific permissions to use this command
  default_member_permissions = "32" # MANAGE_GUILD permission bit
}
