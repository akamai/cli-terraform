terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 7.0.0"
    }
  }
  required_version = ">= 1.0"
}

output "rules_test_include_default" {
  value = data.akamai_property_rules_builder.test_include_default.json
}

output "rule_format_test_include_default" {
  value = data.akamai_property_rules_builder.test_include_default.rule_format
}