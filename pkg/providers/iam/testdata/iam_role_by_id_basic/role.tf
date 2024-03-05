terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_iam_role" "role_id_12345" {
  name          = "Custom role"
  description   = "Custom role\ndescription"
  granted_roles = [992, 707, 452, 677, 726, 296, 457, 987]
}

