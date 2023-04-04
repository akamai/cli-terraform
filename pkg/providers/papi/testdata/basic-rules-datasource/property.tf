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

data "akamai_group" "group" {
  group_name  = "test_group"
  contract_id = "test_contract"
}

data "akamai_contract" "contract" {
  group_name = data.akamai_group.group.group_name
}

resource "akamai_edge_hostname" "test-edgesuite-net" {
  product_id    = "prd_HTTP_Content_Del"
  contract_id   = data.akamai_contract.contract.id
  group_id      = data.akamai_group.group.id
  ip_behavior   = "IPV6_COMPLIANCE"
  edge_hostname = "test.edgesuite.net"
}

resource "akamai_property" "test-edgesuite-net" {
  name        = "test.edgesuite.net"
  contract_id = data.akamai_contract.contract.id
  group_id    = data.akamai_group.group.id
  product_id  = "prd_HTTP_Content_Del"
  rule_format = "v2023-01-05"
  hostnames {
    cname_from             = "test.edgesuite.net"
    cname_to               = akamai_edge_hostname.test-edgesuite-net.edge_hostname
    cert_provisioning_type = "CPS_MANAGED"
  }
  rules = data.akamai_property_rules_builder.test-edgesuite-net_rule_default.json
}

resource "akamai_property_activation" "test-edgesuite-net" {
  property_id = akamai_property.test-edgesuite-net.id
  contact     = ["jsmith@akamai.com"]
  version     = akamai_property.test-edgesuite-net.latest_version
  network     = upper(var.env)
}