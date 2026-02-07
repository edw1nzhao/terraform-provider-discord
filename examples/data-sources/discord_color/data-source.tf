# SPDX-License-Identifier: MPL-2.0

# Convert a hex color to a Discord color integer
data "discord_color" "brand_blue" {
  hex = "#3498DB"
}

# Convert RGB values to a Discord color integer
data "discord_color" "custom_red" {
  rgb = {
    r = 231
    g = 76
    b = 60
  }
}

# Use in a role definition
resource "discord_role" "colored" {
  guild_id = "123456789012345678" # Replace with your guild ID
  name     = "Blue Team"
  color    = data.discord_color.brand_blue.int
}
