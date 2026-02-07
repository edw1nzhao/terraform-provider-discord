# SPDX-License-Identifier: MPL-2.0

# Manage a Discord webhook
resource "discord_webhook" "notifications" {
  channel_id = "123456789012345678" # Replace with your channel ID
  name       = "Build Notifications"
}

# The webhook token and URL are available as computed attributes
output "webhook_url" {
  value     = discord_webhook.notifications.url
  sensitive = true
}
