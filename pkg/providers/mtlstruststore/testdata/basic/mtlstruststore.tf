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

resource "akamai_mtlstruststore_ca_set" "test-ca-set-name" {
  name                = "test-ca-set-name"
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

resource "akamai_mtlstruststore_ca_set_activation" "test-ca-set-name-staging" {
  ca_set_id = akamai_mtlstruststore_ca_set.test-ca-set-name.id
  version   = var.activate_latest_on_staging ? akamai_mtlstruststore_ca_set.test-ca-set-name.latest_version : akamai_mtlstruststore_ca_set.test-ca-set-name.staging_version
  network   = "STAGING"
}

resource "akamai_mtlstruststore_ca_set_activation" "test-ca-set-name-production" {
  ca_set_id = akamai_mtlstruststore_ca_set.test-ca-set-name.id
  version   = var.activate_latest_on_production ? akamai_mtlstruststore_ca_set.test-ca-set-name.latest_version : akamai_mtlstruststore_ca_set.test-ca-set-name.production_version
  network   = "PRODUCTION"
}
