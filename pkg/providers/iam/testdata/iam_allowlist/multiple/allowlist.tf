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

resource "akamai_iam_ip_allowlist" "allowlist" {
  enable = false
}

resource "akamai_iam_cidr_blocks" "cidr_blocks" {
  cidr_blocks = [
    {
      cidr_block = "1.1.1.1/1"
      enabled    = true
      comments   = "comment"
    },
    {
      cidr_block = "2.2.2.2/2"
      enabled    = false
    },
  ]
}
