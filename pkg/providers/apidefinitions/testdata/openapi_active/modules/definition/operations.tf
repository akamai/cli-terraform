terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 6.5.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = "~/.edgerc"
  config_section = "DEFAULT"
}

resource "akamai_apidefinitions_resource_operations" "pet_store" {
  endpoint_id         = 1
  resource_operations = jsonencode(file("${path.module}/operations-api.json"))
}