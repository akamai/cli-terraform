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
}

resource "akamai_cloudcertificates_upload_signed_certificate" "test-name_example_com1234567890" {
  certificate_id = "12345"
  signed_certificate_pem = trimsuffix(<<EOT
-----BEGIN CERTIFICATE-----
testsignedcertificate
-----END CERTIFICATE-----
EOT
  , "\n")
  trust_chain_pem = trimsuffix(<<EOT
-----BEGIN CERTIFICATE-----
testtrustchaincertificate1
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
testtrustchaincertificate2
-----END CERTIFICATE-----
EOT
  , "\n")
}