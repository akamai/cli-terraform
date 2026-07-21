variable "contractid" {
  type        = string
  description = "contract id for zone creation"
}

variable "groupid" {
  type        = string
  description = "group id for zone creation"
}

variable "name" {
  type        = string
  description = "zone name"
}

output "zonename" {
  value = akamai_dns_zone.example_com.name
}

locals {
  zone = var.name
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
  multi_provider_dnssec = false
  target                   = ""
  end_customer_id          = ""
}
