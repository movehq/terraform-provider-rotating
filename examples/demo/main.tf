terraform {
  required_providers {
    rotating = {
      source  = "movehq/rotating"
      version = "1.0.0"
    }
  }
}

variable "rotation_minutes" {
  type        = number
  description = "The number of minutes before the iam creds are rotated."
  default     = 2

  validation {
    condition     = var.rotation_minutes >= 2
    error_message = "rotation_minutes must be greater than or equal to 2"
  }
}

resource "rotating_blue_green" "this" {
  rotate_after_minutes = var.rotation_minutes
}

resource "random_pet" "blue" {
  keepers = {
    uuid = rotating_blue_green.this.blue_uuid
  }
}

resource "random_pet" "green" {
  keepers = {
    uuid = rotating_blue_green.this.green_uuid
  }
}

resource "null_resource" "blue" {
  lifecycle {
    replace_triggered_by = [
      random_pet.blue
    ]
  }
}

resource "null_resource" "green" {
  lifecycle {
    replace_triggered_by = [
      random_pet.green
    ]
  }
}

output "active" {
  value = rotating_blue_green.this.active
}