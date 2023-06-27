terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 5.0.0"
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

resource "akamai_property" "test-edgesuite-net" {
  name        = "test.edgesuite.net"
  contract_id = data.akamai_contract.contract.id
  group_id    = data.akamai_group.group.id
  product_id  = "prd_HTTP_Content_Del"
  rule_format = data.akamai_property_rules_builder.test-edgesuite-net_rule_default.rule_format
  rules       = data.akamai_property_rules_builder.test-edgesuite-net_rule_default.json
}

#resource "akamai_property_activation" "test-edgesuite-net-staging" {
#  property_id                    = akamai_property.test-edgesuite-net.id
#  contact                        = []
#  version                        = akamai_property.test-edgesuite-net.latest_version
#  network                        = "STAGING"
#  auto_acknowledge_rule_warnings = false
#}

#resource "akamai_property_activation" "test-edgesuite-net-production" {
#  property_id                    = akamai_property.test-edgesuite-net.id
#  contact                        = []
#  version                        = akamai_property.test-edgesuite-net.latest_version
#  network                        = "PRODUCTION"
#  auto_acknowledge_rule_warnings = false
#}
