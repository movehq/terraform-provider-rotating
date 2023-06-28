terraform {
  required_providers {
    rotating = {
      source  = "movehq/rotating"
      version = "1.0.0"
    }
  }
}

variable "name" {
  type        = string
  description = "The name of the user to create."
}

variable "rotation_days" {
  type        = number
  description = "The number of days before the iam creds are rotated."
  default     = 30

  validation {
    condition     = var.rotation_days >= 2
    error_message = "rotation_days must be greater than or equal to 2"
  }
}

locals {
  BG = {
    blue  = aws_iam_access_key.blue
    green = aws_iam_access_key.green
  }
}

resource "rotating_blue_green" "this" {
  rotate_after_days = var.rotation_days
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

resource "aws_iam_user" "this" {
  name = var.name
}

resource "aws_iam_access_key" "blue" {
  user = aws_iam_user.this.name

  lifecycle {
    replace_triggered_by = [
      random_pet.blue
    ]
  }
}

resource "aws_iam_access_key" "green" {
  user = aws_iam_user.this.name

  lifecycle {
    replace_triggered_by = [
      random_pet.green
    ]
  }
}

output "active" {
  value = local.BG[rotating_blue_green.this.active]
}

output "blue" {
  value = aws_iam_access_key.blue
}

output "green" {
  value = aws_iam_access_key.green
}

output "name" {
  value = aws_iam_user.this.name
}

output "arn" {
  value = aws_iam_user.this.arn
}
