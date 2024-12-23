terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 6.5.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_iam_ip_allowlist" "allowlist" {
  enable = true
}

resource "akamai_iam_cidr_block" "cidr_1_1_1_1-1" {
  cidr_block = "1.1.1.1/1"
  enabled    = true
  comments   = "comment"
}

