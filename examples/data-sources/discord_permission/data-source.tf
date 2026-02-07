# SPDX-License-Identifier: MPL-2.0

# Compute permission integers from named permissions
# Set a permission to true to include it in allow_bits,
# or false to include it in deny_bits.
# Unset (null) permissions are excluded from both.
data "discord_permission" "moderator" {
  view_channel         = true
  send_messages        = true
  manage_messages      = true
  read_message_history = true
  kick_members         = true
  ban_members          = true
  moderate_members     = true
}

# Use the computed values in channel permission overwrites
resource "discord_channel_permission" "mod_perms" {
  channel_id   = "123456789012345678" # Replace with your channel ID
  overwrite_id = discord_role.moderator.id
  type         = 0 # role

  allow = data.discord_permission.moderator.allow_bits
  deny  = data.discord_permission.moderator.deny_bits
}
