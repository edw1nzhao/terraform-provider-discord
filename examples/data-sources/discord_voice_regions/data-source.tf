# SPDX-License-Identifier: MPL-2.0

# List all available voice regions
data "discord_voice_regions" "all" {}

output "regions" {
  value = data.discord_voice_regions.all.regions
}
