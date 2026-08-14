terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 9.0.0"
    }
  }
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

locals {
  zone = "example.com"
}

resource "akamai_dns_zone" "example_com" {
  contract                 = var.contractid
  group                    = var.groupid
  zone                     = local.zone
  type                     = "PRIMARY"
  masters                  = []
  comment                  = ""
  sign_and_serve           = false
  sign_and_serve_algorithm = ""
  target                   = ""
  end_customer_id          = ""
}


resource "akamai_dns_record" "example_com_example_com_NS" {
  zone       = local.zone
  name       = "example.com"
  recordtype = "NS"
  target     = ["ns1.example.com."]
  ttl        = 3600
}

resource "akamai_dns_record" "example_com_example_com_SOA" {
  zone         = local.zone
  contact      = "admin.example.com."
  expiry       = 604800
  minimum      = 300
  name         = "example.com"
  originserver = "ns1.example.com."
  recordtype   = "SOA"
  refresh      = 3600
  retry        = 600
  ttl          = 3600
}
