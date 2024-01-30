terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 5.6.0"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cloudlets_policy" "policy" {
  name          = "test_policy_export"
  cloudlet_code = "AP"
  description   = "Testing exported policy"
  group_id      = "12345"
  match_rules   = data.akamai_cloudlets_api_prioritization_match_rule.match_rules_ap.json
  is_shared     = true
}

/*
resource "akamai_cloudlets_policy_activation" "policy_activation" {
  policy_id = tonumber(akamai_cloudlets_policy.policy.id)
  network = var.env
  version = akamai_cloudlets_policy.policy.version
}
*/
