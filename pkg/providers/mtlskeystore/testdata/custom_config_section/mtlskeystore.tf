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

resource "akamai_mtlskeystore_client_certificate_akamai" "test_akamai_cert" {
  certificate_name    = "test-akamai-cert"
  contract_id         = "C-0NTR4CT"
  group_id            = 98765
  geography           = "CORE"
  key_algorithm       = "RSA"
  notification_emails = ["test@mail.com"]
  secure_network      = "ENHANCED_TLS"
  subject             = "/C=US/O=Akamai Technologies, Inc./OU=12345 C-0NTR4CT 98765/CN=testCommonName/"
}