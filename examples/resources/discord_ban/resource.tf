# SPDX-License-Identifier: MPL-2.0

# Ban a user from a guild
# Note: All attributes force replacement. Bans cannot be updated in place.
resource "discord_ban" "example" {
  guild_id = "123456789012345678" # Replace with your guild ID
  user_id  = "987654321098765432" # Replace with your user ID

  # Optional reason for the ban (recorded in the audit log)
  reason = "Violation of server rules"

  # Number of seconds of messages to delete (0-604800, i.e., up to 7 days)
  delete_message_seconds = 86400 # Delete 1 day of messages
}
