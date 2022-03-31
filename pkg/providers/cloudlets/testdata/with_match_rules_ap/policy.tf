terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cloudlets_policy" "policy" {
  name              = "test_policy_export"
  cloudlet_code     = "AP"
  description       = "Testing exported policy"
  group_id          = "12345"
  match_rule_format = "1.0"
  match_rules       = data.akamai_cloudlets_api_prioritization_match_rule.match_rules_ap.json
}

/*
resource "akamai_cloudlets_policy_activation" "policy_activation" {
  policy_id = tonumber(akamai_cloudlets_policy.policy.id)
  network = var.env
  version = akamai_cloudlets_policy.policy.version
  associated_properties = [ "UNKNOWN_CHANGE_ME" ]
}
*/
