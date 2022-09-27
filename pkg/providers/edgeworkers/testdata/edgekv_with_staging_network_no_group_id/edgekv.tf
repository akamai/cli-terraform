terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = "~> 2.0.0"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_edgekv" "edgekv" {
  namespace_name       = "test_namespace"
  network              = "staging"
  group_id             = 0
  retention_in_seconds = 0
  geo_location         = "EU"
}
