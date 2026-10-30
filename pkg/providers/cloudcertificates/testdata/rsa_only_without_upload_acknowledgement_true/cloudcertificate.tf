terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 12.0.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cloudcertificates_lineage" "rsa-only-without-upload-acknowledgement-true" {
  contract_id    = "ctr_1-ABC"
  group_id       = 123
  geo_class      = "STANDARD_WORLDWIDE"
  lineage_name   = "rsa-only-without-upload-acknowledgement-true"
  lineage_type   = "MULTIPLE_GENERATION"
  secure_network = "ENHANCED_TLS"
  key_specs = {
    "RSA" = "2048"
  }
  sans = [
    "*.example.com",
    "example.com",
  ]
  subject = {
    common_name         = "example.com"
    organization        = "Example, Inc."
    organizational_unit = "Security"
    country             = "US"
    state               = "Massachusetts"
    locality            = "Cambridge"
  }
}
/*

resource "akamai_cloudcertificates_upload" "rsa-only-without-upload-acknowledgement-true" {
  lineage_id = akamai_cloudcertificates_lineage.rsa-only-without-upload-acknowledgement-true.lineage_id
  acknowledge_warnings = true
  algorithms = {
    # Paste the PEM-encoded signed certificate issued by your CA (exactly one certificate)
    # in the signed_certificate_pem attribute.
    # Keep the BEGIN CERTIFICATE/END CERTIFICATE lines; replace only the body.
    "RSA" = {
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
  }
}
*/
/*
resource "akamai_cloudcertificates_activation" "rsa-only-without-upload-acknowledgement-true" {
  lineage_id = akamai_cloudcertificates_upload.rsa-only-without-upload-acknowledgement-true.lineage_id
  staging_generation_id = (
    var.ccm_activate_latest_on_staging
    ? akamai_cloudcertificates_upload.rsa-only-without-upload-acknowledgement-true.generation_id
    : akamai_cloudcertificates_lineage.rsa-only-without-upload-acknowledgement-true.current_staging.generation_id
  )
  production_generation_id = (
    var.ccm_activate_latest_on_production
    ? akamai_cloudcertificates_upload.rsa-only-without-upload-acknowledgement-true.generation_id
    : akamai_cloudcertificates_lineage.rsa-only-without-upload-acknowledgement-true.current_production.generation_id
  )
}
*/
