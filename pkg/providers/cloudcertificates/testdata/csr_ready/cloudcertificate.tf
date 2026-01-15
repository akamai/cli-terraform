terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 9.3.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cloudcertificates_certificate" "test-name_example_com1234567890" {
  base_name      = "test-name.example.com1234567890"
  contract_id    = "test_contract"
  key_size       = "2048"
  key_type       = "RSA"
  secure_network = "ENHANCED_TLS"
  sans           = ["test.example.com", "test.example2.com"]
  subject = {
    common_name  = "test.example.com"
    organization = "Test Org"
    country      = "US"
    state        = "CA"
    locality     = "Test City"
  }
}

/*
# Paste the PEM‑encoded signed certificate issued by your CA (exactly one certificate)
# in the signed_certificate_pem attribute.
# Keep the BEGIN CERTIFICATE/END CERTIFICATE lines; replace only the body.

# For the trust_chain_pem attribute, you can specify one or more PEM‑encoded certificates.
# Multiple trust chain certificates should each be enclosed in their own BEGIN CERTIFICATE/END CERTIFICATE block.
# If a trust chain is not required, you can remove the trust_chain_pem attribute entirely.

resource "akamai_cloudcertificates_upload_signed_certificate" "test-name_example_com1234567890" {
  certificate_id         = "12345"
  signed_certificate_pem = <<-EOT
-----BEGIN CERTIFICATE-----
Please place your signed certificate here
-----END CERTIFICATE-----
EOT
  trust_chain_pem = <<-EOT
-----BEGIN CERTIFICATE-----
Please place your trust chain here (optional)
-----END CERTIFICATE-----
EOT
}
*/