// SIEM Settings
resource "akamai_appsec_siem_settings" "siem" {
  config_id               = local.config_id
  enable_siem             = true
  enable_for_all_policies = false
  enable_botman_siem      = false
  siem_id                 = 1
  security_policy_ids     = ["ASE1_156138"]
}
