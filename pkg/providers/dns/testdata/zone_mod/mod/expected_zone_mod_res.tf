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
  value = akamai_dns_zone._0007770b-08a8-4b5f-a46b-081b772ba605-test_com.name
}

locals {
  zone = var.name
}

resource "akamai_dns_zone" "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com" {
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
