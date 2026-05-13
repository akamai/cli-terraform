variable "zonename" {
  type        = string
  description = "zone name for this name record set config"
}

locals {
  zone = var.zonename
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
