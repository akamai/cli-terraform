variable "zonename" {
  type        = string
  description = "zone name for this name record set config"
}

locals {
  zone = var.zonename
}

resource "akamai_dns_record" "example_com_example_com_NS" {
  zone       = local.zone
  name       = "example.com"
  recordtype = "NS"
  target     = ["ns1.example.com."]
  ttl        = 3600
}
