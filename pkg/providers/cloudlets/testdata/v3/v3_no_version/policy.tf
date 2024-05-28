terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 6.2.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cloudlets_policy" "policy" {
  name          = "test_policy_export"
  cloudlet_code = "ER"
  description   = ""
  group_id      = "234"
  is_shared     = true
}

/*
resource "akamai_cloudlets_policy_activation" "policy_activation" {
  policy_id = tonumber(akamai_cloudlets_policy.policy.id)
  network = var.env
  version = akamai_cloudlets_policy.policy.version
}
*/
