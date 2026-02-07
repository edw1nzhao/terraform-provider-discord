# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    discord = {
      source  = "edw1nzhao/discord"
      version = "~> 0.1"
    }
  }
}

# Configure the Discord provider
provider "discord" {
  # The bot token can also be set via the DISCORD_TOKEN environment variable
  token = var.discord_token

  # Application ID is needed for application command resources
  # Can also be set via DISCORD_APPLICATION_ID environment variable
  application_id = var.discord_application_id
}

variable "discord_token" {
  type      = string
  sensitive = true
}

variable "discord_application_id" {
  type    = string
  default = ""
}
