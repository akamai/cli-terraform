
data "akamai_property_rules_builder" "test-edgesuite-net_rule_default" {
  rules_v2023_01_05 {
    name      = "default"
    is_secure = false
    children = [
      data.akamai_property_rules_builder.test-edgesuite-net_rule_new_rule.json,
    ]
  }
}
