terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 8.0.0"
    }
  }
  required_version = ">= 1.0"
}

output "rules" {
  value = data.akamai_property_rules_builder.test-edgesuite-net.json
}

output "rule_format" {
  value = data.akamai_property_rules_builder.test-edgesuite-net.rule_format
}