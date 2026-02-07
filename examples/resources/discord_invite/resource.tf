# SPDX-License-Identifier: MPL-2.0

# Create a channel invite
# Note: Invites cannot be updated. Any change forces recreation.
resource "discord_invite" "example" {
  channel_id = "123456789012345678" # Replace with your channel ID

  # Duration in seconds before expiry (0 = never). Default: 86400 (24 hours)
  max_age = 0

  # Max number of uses (0 = unlimited). Default: 0
  max_uses = 10

  # Whether this invite only grants temporary membership. Default: false
  temporary = false

  # If true, don't try to reuse a similar invite. Default: false
  unique = true
}

output "invite_url" {
  value = discord_invite.example.url
}
