
data "akamai_property_rules_template" "rules_test_include" {
  template_file = abspath("${path.module}/property-snippets/test_include.json")
}

/*
data "akamai_property_include_parents" "include_parents" {
  contract_id = "test_contract"
  group_id    = "test_group"
  include_id  = "inc_123456"
}
*/

resource "akamai_property_include" "test_include" {
  contract_id = "test_contract"
  group_id    = "test_group"
  name        = "test_include"
  rule_format = "v2020-11-02"
  type        = "MICROSERVICES"
  rules       = data.akamai_property_rules_template.rules_test_include.json
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
  compliance_record {
    noncompliance_reason = "NONE"
    customer_email       = "customer@email.com"
    peer_reviewed_by     = "John Doe"
    unit_tested          = true
  }
}
data "akamai_property_rules_template" "rules_test_include_1" {
  template_file = abspath("${path.module}/property-snippets/test_include_1.json")
}

/*
data "akamai_property_include_parents" "include_parents" {
  contract_id = "test_contract"
  group_id    = "test_group"
  include_id  = "inc_78910"
}
*/

resource "akamai_property_include" "test_include_1" {
  contract_id = "test_contract"
  group_id    = "test_group"
  name        = "test_include_1"
  rule_format = "v2020-11-02"
  type        = "MICROSERVICES"
  rules       = data.akamai_property_rules_template.rules_test_include_1.json
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