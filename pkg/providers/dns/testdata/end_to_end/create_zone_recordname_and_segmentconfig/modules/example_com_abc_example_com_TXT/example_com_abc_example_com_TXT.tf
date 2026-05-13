variable "zonename" {
  type        = string
  description = "zone name for this name record set config"
}

locals {
  zone = var.zonename
}

resource "akamai_dns_record" "example_com_abc_example_com_TXT" {
  zone       = local.zone
  name       = "abc.example.com"
  recordtype = "TXT"
  target     = ["\"dummy text abc\""]
  ttl        = 300
}
