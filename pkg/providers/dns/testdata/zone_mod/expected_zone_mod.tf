terraform {
  required_version = ">= 0.13"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = "~> 1.6.1"
    }
  }
}

locals {
  zone = "0007770b-08a8-4b5f-a46b-081b772ba605-test.com"
}

module "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com" {
  source = "modules/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com"

  contract = var.contractid
  group    = var.groupid
  name     = local.zone
}
