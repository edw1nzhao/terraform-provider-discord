# SPDX-License-Identifier: MPL-2.0

# Manage a guild sticker
# Note: Sticker file upload is not supported via the API. Create the sticker
# in Discord first, then import it to bring it under Terraform management.
resource "discord_guild_sticker" "example" {
  guild_id    = "123456789012345678" # Replace with your guild ID
  name        = "my_sticker"
  description = "A cool custom sticker"
  tags        = "cool,sticker,custom"
}
