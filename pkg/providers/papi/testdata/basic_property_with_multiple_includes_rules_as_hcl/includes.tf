

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
  rule_format = data.akamai_property_rules_builder.test_include_rule_default.rule_format
  rules       = data.akamai_property_rules_builder.test_include_rule_default.json
}

resource "akamai_property_include_activation" "test_include_staging" {
  contract_id                    = akamai_property_include.test_include.contract_id
  group_id                       = akamai_property_include.test_include.group_id
  include_id                     = akamai_property_include.test_include.id
  network                        = "STAGING"
  auto_acknowledge_rule_warnings = false
  version                        = "1"
  note                           = "test staging activation"
  notify_emails                  = ["test@example.com"]
}

resource "akamai_property_include_activation" "test_include_production" {
  contract_id                    = akamai_property_include.test_include.contract_id
  group_id                       = akamai_property_include.test_include.group_id
  include_id                     = akamai_property_include.test_include.id
  network                        = "PRODUCTION"
  auto_acknowledge_rule_warnings = false
  version                        = "1"
  note                           = "test production activation"
  notify_emails                  = ["test@example.com", "test1@example.com"]
}

/*
data "akamai_property_include_parents" "test_include_1" {
  contract_id = "test_contract"
  group_id    = "test_group"
  include_id  = "inc_78910"
}
*/

resource "akamai_property_include" "test_include_1" {
  contract_id = "test_contract"
  group_id    = "test_group"
  name        = "test_include_1"
  product_id  = "test_product2"
  type        = "MICROSERVICES"
  rule_format = data.akamai_property_rules_builder.test_include_1_rule_default.rule_format
  rules       = data.akamai_property_rules_builder.test_include_1_rule_default.json
}

resource "akamai_property_include_activation" "test_include_1_staging" {
  contract_id                    = akamai_property_include.test_include_1.contract_id
  group_id                       = akamai_property_include.test_include_1.group_id
  include_id                     = akamai_property_include.test_include_1.id
  network                        = "STAGING"
  auto_acknowledge_rule_warnings = false
  version                        = "1"
  note                           = "test staging activation"
  notify_emails                  = ["test@example.com"]
}

#resource "akamai_property_include_activation" "test_include_1_production" {
#  contract_id = akamai_property_include.test_include_1.contract_id
#  group_id = akamai_property_include.test_include_1.group_id
#  include_id = akamai_property_include.test_include_1.id
#  network = "PRODUCTION"
#  auto_acknowledge_rule_warnings = false
#  version = "0"
#  notify_emails = []
#}