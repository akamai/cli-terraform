terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 6.6.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cloudaccess_key" "TestKeyName" {
  access_key_name       = "TestKeyName"
  authentication_method = "VP_QUEUE_IT"
  group_id              = 1234
  contract_id           = "C-Contract123"
  network_configuration = {
    security_network = "ENHANCED_TLS"
  }
  credentials_a = {
    cloud_secret_access_key = var.secret_access_key_a
    primary_key             = false
  }
  credentials_b = {
    cloud_secret_access_key = var.secret_access_key_b
    primary_key             = false
  }
}