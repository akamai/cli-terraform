terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 5.4.0"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_property" "test-edgesuite-net" {
  name        = "test.edgesuite.net"
  contract_id = var.contract_id
  group_id    = var.group_id
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
