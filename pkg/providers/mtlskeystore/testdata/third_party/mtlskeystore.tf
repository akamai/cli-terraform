terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 8.1.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_mtlskeystore_client_certificate_third_party" "test_third_party_cert" {
  certificate_name    = "test-third-party-cert"
  contract_id         = "C-0NTR4CT"
  group_id            = 98765
  geography           = "RUSSIA_AND_CORE"
  key_algorithm       = "ECDSA"
  notification_emails = ["test@mail.com"]
  secure_network      = "STANDARD_TLS"
  subject             = "/C=US/O=Akamai Technologies, Inc./OU=12345 C-0NTR4CT 98765/CN=testCommonName/"
  versions = {
    "2025-10-01T12:00:00_v3" = {},
    "2024-10-01T12:00:00_v2" = {},
    "2023-10-01T12:00:00_v1" = {}
  }
}