
data "akamai_property_rules_builder" "include2_rule_default" {
  rules_v2023_01_05 {
    name      = "default"
    is_secure = false
  }
}
