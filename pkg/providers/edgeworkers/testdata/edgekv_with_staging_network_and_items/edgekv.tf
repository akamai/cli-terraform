terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 3.6.0"
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
  network              = "staging"
  group_id             = 123
  retention_in_seconds = 0
  geo_location         = "EU"
}

resource "akamai_edgekv_group_items" "group1" {
  namespace_name = akamai_edgekv.edgekv.namespace_name
  network        = "staging"
  group_name     = "group1"
  items = {
    "item1.1" = "value1.1"
    "item1.2" = "value1.2"
  }
}

resource "akamai_edgekv_group_items" "group2" {
  namespace_name = akamai_edgekv.edgekv.namespace_name
  network        = "staging"
  group_name     = "group2"
  items = {
    "item2.1" = "value2.1"
    "item2.2" = "value\n2.2"
  }
}
