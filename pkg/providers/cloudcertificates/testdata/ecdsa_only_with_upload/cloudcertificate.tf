terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 11.1.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cloudcertificates_lineage" "ecdsa-only-with-upload" {
  contract_id    = "ctr_1-ABC"
  group_id       = 123
  geo_class      = "STANDARD_WORLDWIDE"
  lineage_name   = "ecdsa-only-with-upload"
  lineage_type   = "MULTIPLE_GENERATION"
  secure_network = "ENHANCED_TLS"
  key_specs = {
    "ECDSA" = "P-256"
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

resource "akamai_cloudcertificates_upload" "ecdsa-only-with-upload" {
  lineage_id = akamai_cloudcertificates_lineage.ecdsa-only-with-upload.lineage_id
  algorithms = {
    "ECDSA" = {
      signed_certificate_pem = <<EOT
-----BEGIN CERTIFICATE-----
ecdsa-signed-certificate
-----END CERTIFICATE-----
EOT
    }
  }
}
/*
resource "akamai_cloudcertificates_activation" "ecdsa-only-with-upload" {
  lineage_id = akamai_cloudcertificates_upload.ecdsa-only-with-upload.lineage_id
  staging_generation_id = (
    var.ccm_activate_latest_on_staging
    ? akamai_cloudcertificates_upload.ecdsa-only-with-upload.generation_id
    : akamai_cloudcertificates_lineage.ecdsa-only-with-upload.current_staging.generation_id
  )
  production_generation_id = (
    var.ccm_activate_latest_on_production
    ? akamai_cloudcertificates_upload.ecdsa-only-with-upload.generation_id
    : akamai_cloudcertificates_lineage.ecdsa-only-with-upload.current_production.generation_id
  )
}
*/
