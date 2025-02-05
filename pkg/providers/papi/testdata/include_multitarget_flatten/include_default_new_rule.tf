
data "akamai_property_rules_builder" "include_rule_new_rule" {
  rules_v2023_01_05 {
    name = "New Rule"
    children = [
      data.akamai_property_rules_builder.include_rule_new_rule_1.json,
    ]
  }
}
