
data "akamai_property_rules_builder" "test-edgesuite-net_rule_default" {
  rules_v2023_01_05 {
    name      = "default"
    is_secure = false
    children = [
      data.akamai_property_rules_builder.test-edgesuite-net_rule_new_rule.json,
    ]
  }
}

data "akamai_property_rules_builder" "test-edgesuite-net_rule_new_rule" {
  rules_v2023_01_05 {
    name = "New Rule"
    children = [
      data.akamai_property_rules_builder.test-edgesuite-net_rule_new_rule_1.json,
    ]
  }
}

data "akamai_property_rules_builder" "test-edgesuite-net_rule_new_rule_1" {
  rules_v2023_01_05 {
    name = "New Rule 1"
  }
}
