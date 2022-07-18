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
  zone = "0007770b-08a8-4b5f-a46b-081b772ba605-sbodden-calvin.com"
}

resource "akamai_dns_zone" "_0007770b-08a8-4b5f-a46b-081b772ba605-sbodden-calvin_com" {
  contract                 = var.contractid
  group                    = var.groupid
  comment                  = ""
  end_customer_id          = ""
  masters                  = []
  sign_and_serve           = false
  sign_and_serve_algorithm = ""
  target                   = ""
  type                     = "PRIMARY"
  zone                     = local.zone
  tsig_key {
    name      = "some-name"
    algorithm = "some-algorithm"
    secret    = "some-secret"
  }
}

