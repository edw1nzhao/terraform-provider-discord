# SPDX-License-Identifier: MPL-2.0

# Manage a global slash command (available in all guilds)
resource "discord_global_application_command" "ping" {
  name        = "ping"
  description = "Check if the bot is alive"

  # Command type: 1=CHAT_INPUT (default), 2=USER, 3=MESSAGE
  type = 1
}

# Global command with options (using JSON-encoded options)
resource "discord_global_application_command" "greet" {
  name        = "greet"
  description = "Greet a user"

  options = jsonencode([
    {
      name        = "user"
      description = "The user to greet"
      type        = 6 # USER type
      required    = true
    },
    {
      name        = "message"
      description = "Custom greeting message"
      type        = 3 # STRING type
      required    = false
    }
  ])
}
