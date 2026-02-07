# SPDX-License-Identifier: MPL-2.0

# Manage guild onboarding configuration
# Note: Deleting this resource disables onboarding.
resource "discord_guild_onboarding" "example" {
  guild_id = "123456789012345678" # Replace with your guild ID
  enabled  = true

  # Mode: 0=ONBOARDING_DEFAULT, 1=ONBOARDING_ADVANCED
  mode = 0

  # Channels that members get opted into automatically
  default_channel_ids = [
    "111111111111111111", # Replace with your default channel IDs
    "222222222222222222", # Replace with your default channel IDs
  ]

  # Onboarding prompts
  prompts = [
    {
      # Type: 0=MULTIPLE_CHOICE, 1=DROPDOWN
      type          = 0
      title         = "What are you interested in?"
      single_select = false
      required      = true
      in_onboarding = true

      options = [
        {
          title       = "Gaming"
          description = "Join our gaming channels"
          channel_ids = ["333333333333333333"] # Replace with your channel IDs
          role_ids    = ["444444444444444444"]  # Replace with your role IDs
          emoji_name  = "\ud83c\udfae"
        },
        {
          title       = "Art"
          description = "Join our art channels"
          channel_ids = ["555555555555555555"] # Replace with your channel IDs
          role_ids    = ["666666666666666666"]  # Replace with your role IDs
          emoji_name  = "\ud83c\udfa8"
        },
      ]
    },
  ]
}
