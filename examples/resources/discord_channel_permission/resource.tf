# SPDX-License-Identifier: MPL-2.0

# Manage a channel permission overwrite for a role
resource "discord_channel_permission" "example" {
  channel_id   = "123456789012345678" # Replace with your channel ID
  overwrite_id = discord_role.example.id # Role or user ID
  type         = 0                       # 0=role, 1=member

  # Permission bitfield values as strings
  # Use the discord_permission data source to compute these
  allow = data.discord_permission.allow.allow_bits
  deny  = data.discord_permission.deny.allow_bits
}

# Compute the allow permission integer
data "discord_permission" "allow" {
  view_channel = true
  send_messages = true
  read_message_history = true
}

# Compute the deny permission integer
data "discord_permission" "deny" {
  manage_messages = true
}
