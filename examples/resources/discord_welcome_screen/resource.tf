# SPDX-License-Identifier: MPL-2.0

# Manage the guild welcome screen
# Note: Deleting this resource disables the welcome screen and clears channels.
resource "discord_welcome_screen" "example" {
  guild_id    = "123456789012345678" # Replace with your guild ID
  enabled     = true
  description = "Welcome to our community! Pick some channels to get started."

  # Up to 5 welcome screen channels
  welcome_channels = [
    {
      channel_id  = "111111111111111111" # Replace with your channel ID
      description = "Get the latest server announcements"
      emoji_name  = "\ud83d\udce2"
    },
    {
      channel_id  = "222222222222222222" # Replace with your channel ID
      description = "Introduce yourself to the community"
      emoji_name  = "\ud83d\udc4b"
    },
    {
      channel_id  = "333333333333333333" # Replace with your channel ID
      description = "Chat with other members"
      emoji_name  = "\ud83d\udcac"
    },
  ]
}
