# SPDX-License-Identifier: MPL-2.0

locals {
  guild_id = "123456789012345678" # Replace with your guild ID
}

# Look up a role by ID
data "discord_role" "by_id" {
  guild_id = local.guild_id
  id       = "987654321098765432" # Replace with your role ID
}

# Look up a role by name
data "discord_role" "by_name" {
  guild_id = local.guild_id
  name     = "Moderator"
}

output "role_permissions" {
  value = data.discord_role.by_name.permissions
}
