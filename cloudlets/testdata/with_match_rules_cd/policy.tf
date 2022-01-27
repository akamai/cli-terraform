terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cloudlets_policy" "policy" {
  name = "test_policy_export"
  cloudlet_code = "CD"
  description = "Testing exported policy"
  group_id = "12345"
  match_rule_format = "1.0"
  match_rules = data.akamai_cloudlets_phased_release_match_rule.match_rules_cd.json
}
