terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 5.6.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}


data "akamai_property_rules_template" "rules_test_include" {
  template_file = abspath("${path.module}/property-snippets/test_include.json")
}

/*
data "akamai_property_include_parents" "test_include" {
  contract_id = "test_contract"
  group_id    = "test_group"
  include_id  = "inc_123456"
}
*/

resource "akamai_property_include" "test_include" {
  contract_id = "test_contract"
  group_id    = "test_group"
  name        = "test_include"
  product_id  = "test_product"
  type        = "MICROSERVICES"
  rule_format = "v2020-11-02"
  rules       = data.akamai_property_rules_template.rules_test_include.json
}

#resource "akamai_property_include_activation" "test_include_staging" {
#  contract_id = akamai_property_include.test_include.contract_id
#  group_id = akamai_property_include.test_include.group_id
#  include_id = akamai_property_include.test_include.id
#  network = "STAGING"
#  auto_acknowledge_rule_warnings = false
#  version = var.activate_latest_on_staging ? akamai_property_include.test_include.latest_version : akamai_property_include.test_include.staging_version
#  notify_emails = []
#}

#resource "akamai_property_include_activation" "test_include_production" {
#  contract_id = akamai_property_include.test_include.contract_id
#  group_id = akamai_property_include.test_include.group_id
#  include_id = akamai_property_include.test_include.id
#  network = "PRODUCTION"
#  auto_acknowledge_rule_warnings = false
#  version = var.activate_latest_on_production ? akamai_property_include.test_include.latest_version : akamai_property_include.test_include.production_version
#  notify_emails = []
#}