terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cps_dv_enrollment" "enrollment_id_1" {
  common_name    = "test.akamai.com"
  sans           = []
  secure_network = "enhanced-tls"
  sni_only       = true
  admin_contact {
    first_name       = "R1"
    last_name        = "D1"
    organization     = "Akamai"
    email            = "r1d1@akamai.com"
    phone            = "123123123"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    region           = "MA"
    postal_code      = "12345"
    country_code     = "US"
  }
  certificate_chain_type = "default"
  csr {
    country_code        = "US"
    city                = "Cambridge"
    organization        = "Akamai"
    organizational_unit = "WebEx"
    state               = "MA"
  }
  network_configuration {
    disallowed_tls_versions = ["TLSv1", "TLSv1_1", ]
    geography               = "core"
    must_have_ciphers       = "ak-akamai-default"
    ocsp_stapling           = "on"
    preferred_ciphers       = "ak-akamai-default"
  }
  signature_algorithm = "SHA-256"
  tech_contact {
    first_name       = "R2"
    last_name        = "D2"
    organization     = "Akamai"
    email            = "r2d2@akamai.com"
    phone            = "123123123"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    region           = "MA"
    postal_code      = "12345"
    country_code     = "US"
  }
  organization {
    name             = "Akamai"
    phone            = "321321321"
    address_line_one = "150 Broadway"
    city             = "Cambridge"
    region           = "MA"
    postal_code      = "12345"
    country_code     = "US"
  }
  contract_id = "ctr_1"
}