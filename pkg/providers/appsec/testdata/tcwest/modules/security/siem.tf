// SIEM Settings
resource "akamai_appsec_siem_settings" "siem" {
  config_id               = akamai_appsec_configuration.config.config_id
  enable_siem             = true
  enable_for_all_policies = false
  enable_botman_siem      = true
  siem_id                 = 1
  security_policy_ids     = ["easy_80433"]
}
