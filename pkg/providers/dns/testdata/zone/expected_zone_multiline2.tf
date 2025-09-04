terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 9.0.0"
    }
  }
}

locals {
  zone = "0007770b-08a8-4b5f-a46b-081b772ba605-test.com"
}

resource "akamai_dns_zone" "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com" {
  contract                 = var.contractid
  group                    = var.groupid
  comment                  = <<EOT
first
second
EOT
  end_customer_id          = ""
  masters                  = ["1.1.1.1"]
  sign_and_serve           = false
  sign_and_serve_algorithm = ""
  outbound_zone_transfer {
    acl            = ["192.0.2.156/24"]
    enabled        = true
    notify_targets = ["192.0.2.192"]
    tsig_key {
      algorithm = "hmac-sha1"
      name      = "other.com.akamai.com"
      secret    = "fakeSecretajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw=="
    }
  }
  target = ""
  type   = "SECONDARY"
  zone   = local.zone
  tsig_key {
    name      = "some-name"
    algorithm = "some-algorithm"
    secret    = "some-secret"
  }
}

