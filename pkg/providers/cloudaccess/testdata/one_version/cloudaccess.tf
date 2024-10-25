terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 6.3.0"
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
  authentication_method = "AWS4_HMAC_SHA256"
  group_id              = 1234
  contract_id           = "C-Contract123"
  network_configuration = {
    additional_cdn   = "CHINA_CDN"
    security_network = "ENHANCED_TLS"
  }
  credentials_a = {
    cloud_access_key_id     = "testAccessKey1"
    cloud_secret_access_key = var.secret_access_key_a
    primary_key             = false
  }
}