# SPDX-License-Identifier: MPL-2.0

# Manage the set of roles for a guild member
resource "discord_member_roles" "example" {
  guild_id = "123456789012345678" # Replace with your guild ID
  user_id  = "987654321098765432" # Replace with your user ID

  # Set of role IDs to assign to the member
  roles = [
    discord_role.moderator.id,
    discord_role.member.id,
  ]
}
