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
  multi_provider_dnssec    = false
  target                   = ""
  end_customer_id          = ""
}


resource "akamai_dns_record" "example_com_abc_example_com_TXT" {
  zone       = local.zone
  name       = "abc.example.com"
  recordtype = "TXT"
  target     = ["\"dummy text abc\""]
  ttl        = 300
}

resource "akamai_dns_record" "example_com_def_example_com_TXT" {
  zone       = local.zone
  name       = "def.example.com"
  recordtype = "TXT"
  target     = ["\"dummy text def\""]
  ttl        = 300
}
