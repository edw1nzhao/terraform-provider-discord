# SPDX-License-Identifier: MPL-2.0

# Manage the Discord application settings
# Note: The application cannot be created or deleted via the API.
# Create behaves like Update (PATCH). Delete is a no-op.
resource "discord_application" "example" {
  description = "A bot that manages our community server"

  # Tags describing the application (up to 5)
  tags = ["moderation", "community", "utilities"]

  # Whether the bot can be added by anyone
  bot_public = true

  # Whether the bot requires the OAuth2 code grant
  bot_require_code_grant = false
}
