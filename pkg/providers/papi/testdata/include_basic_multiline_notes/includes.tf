terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 5.3.0"
    }
  }
  required_version = ">= 0.13"
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
  type        = "MICROSERVICES"
  rule_format = "v2020-11-02"
  rules       = data.akamai_property_rules_template.rules_test_include.json
}

resource "akamai_property_include_activation" "test_include_staging" {
  contract_id                    = akamai_property_include.test_include.contract_id
  group_id                       = akamai_property_include.test_include.group_id
  include_id                     = akamai_property_include.test_include.id
  network                        = "STAGING"
  auto_acknowledge_rule_warnings = false
  version                        = "1"
  note = trimsuffix(<<EOT
first
second

last
EOT
  , "\n")
  notify_emails = ["test@example.com"]
}

resource "akamai_property_include_activation" "test_include_production" {
  contract_id                    = akamai_property_include.test_include.contract_id
  group_id                       = akamai_property_include.test_include.group_id
  include_id                     = akamai_property_include.test_include.id
  network                        = "PRODUCTION"
  auto_acknowledge_rule_warnings = false
  version                        = "1"
  note                           = <<EOT
first
second
EOT
  notify_emails                  = ["test@example.com", "test1@example.com"]
}