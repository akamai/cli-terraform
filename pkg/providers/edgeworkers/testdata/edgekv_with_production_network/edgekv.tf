terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 8.0.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_edgekv" "edgekv" {
  namespace_name       = "test_namespace"
  network              = "production"
  group_id             = 123
  retention_in_seconds = 0
  geo_location         = "EU"
}
