# SPDX-License-Identifier: MPL-2.0

# Manage a Discord guild (server)
resource "discord_guild" "example" {
  name = "My Terraform Server"

  # Verification level: 0=None, 1=Low, 2=Medium, 3=High, 4=Very High
  verification_level = 1

  # Default notification level: 0=All Messages, 1=Only Mentions
  default_message_notifications = 1

  # Explicit content filter: 0=Disabled, 1=Members without roles, 2=All members
  explicit_content_filter = 2

  # AFK timeout in seconds (must be 60, 300, 900, 1800, or 3600)
  afk_timeout = 300

  # Enable the server boost progress bar
  premium_progress_bar_enabled = true
}
