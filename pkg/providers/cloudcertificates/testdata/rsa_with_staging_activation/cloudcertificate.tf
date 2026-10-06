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

resource "akamai_cloudcertificates_lineage" "rsa-with-staging-activation" {
  contract_id    = "ctr_1-ABC"
  group_id       = 123
  geo_class      = "STANDARD_WORLDWIDE"
  lineage_name   = "rsa-with-staging-activation"
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

resource "akamai_cloudcertificates_upload" "rsa-with-staging-activation" {
  lineage_id = akamai_cloudcertificates_lineage.rsa-with-staging-activation.lineage_id
  algorithms = {
    "RSA" = {
      signed_certificate_pem = <<EOT
-----BEGIN CERTIFICATE-----
rsa-signed-certificate
-----END CERTIFICATE-----
EOT
      trust_chain_pem        = <<EOT
-----BEGIN CERTIFICATE-----
test-trust-chain
-----END CERTIFICATE-----
EOT
    }
  }
}
resource "akamai_cloudcertificates_activation" "rsa-with-staging-activation" {
  lineage_id = akamai_cloudcertificates_upload.rsa-with-staging-activation.lineage_id
  staging_generation_id = (
    var.ccm_activate_latest_on_staging
    ? akamai_cloudcertificates_upload.rsa-with-staging-activation.generation_id
    : akamai_cloudcertificates_lineage.rsa-with-staging-activation.current_staging.generation_id
  )
  # production_generation_id = (
  #   var.ccm_activate_latest_on_production
  #   ? akamai_cloudcertificates_upload.rsa-with-staging-activation.generation_id
  #   : akamai_cloudcertificates_lineage.rsa-with-staging-activation.current_production.generation_id
  # )
}
