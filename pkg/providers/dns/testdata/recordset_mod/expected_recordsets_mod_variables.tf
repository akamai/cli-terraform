variable "zonename" {
  type        = string
  description = "zone name for this name record set config"
}

locals {
  zone = var.zonename
}

resource "akamai_dns_record" "zoneName_someName_someType" {
  zone       = local.zone
  hardware   = "INTEL-386"
  software   = "Unix"
  name       = "someName"
  recordtype = "someType"
  ttl        = 1000
}
