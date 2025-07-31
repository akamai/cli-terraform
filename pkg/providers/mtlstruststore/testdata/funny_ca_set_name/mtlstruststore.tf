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

resource "akamai_mtlstruststore_ca_set" "a_funny-set_name" {
  name                = "a_funny-set.name"
  allow_insecure_sha1 = false
  certificates = [
    {
      certificate_pem = <<EOT
-----BEGIN CERTIFICATE-----
FAKECERTSTARTSEQ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKL
MNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMN
OPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOO
PQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNEND
SEQ==
-----END CERTIFICATE-----
EOT
    }
  ]
}

resource "akamai_mtlstruststore_ca_set_activation" "a_funny-set_name-staging" {
  ca_set_id = akamai_mtlstruststore_ca_set.a_funny-set_name.id
  version   = var.activate_latest_on_staging ? akamai_mtlstruststore_ca_set.a_funny-set_name.latest_version : akamai_mtlstruststore_ca_set.a_funny-set_name.staging_version
  network   = "STAGING"
}

resource "akamai_mtlstruststore_ca_set_activation" "a_funny-set_name-production" {
  ca_set_id = akamai_mtlstruststore_ca_set.a_funny-set_name.id
  version   = var.activate_latest_on_production ? akamai_mtlstruststore_ca_set.a_funny-set_name.latest_version : akamai_mtlstruststore_ca_set.a_funny-set_name.production_version
  network   = "PRODUCTION"
}
