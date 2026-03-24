// Enable/Disable Protections for policy policy2
resource "akamai_appsec_security_policy_protections" "policy2" {
  config_id                         = local.config_id
  security_policy_id                = akamai_appsec_security_policy.policy2.security_policy_id
  apply_account_protection_controls = false
  apply_api_constraints             = true
  apply_application_layer_controls  = true
  apply_botman_controls             = true
  apply_malware_controls            = true
  apply_network_layer_controls      = true
  apply_rate_controls               = true
  apply_reputation_controls         = true
  apply_slow_post_controls          = true
  apply_url_protection_controls     = true
}

// Enable/Disable Protections for policy andrew
resource "akamai_appsec_security_policy_protections" "andrew" {
  config_id                         = local.config_id
  security_policy_id                = akamai_appsec_security_policy.andrew.security_policy_id
  apply_account_protection_controls = false
  apply_api_constraints             = true
  apply_application_layer_controls  = true
  apply_botman_controls             = true
  apply_malware_controls            = true
  apply_network_layer_controls      = true
  apply_rate_controls               = true
  apply_reputation_controls         = true
  apply_slow_post_controls          = true
  apply_url_protection_controls     = true
}

// Enable/Disable Protections for policy policy1
resource "akamai_appsec_security_policy_protections" "policy1" {
  config_id                         = local.config_id
  security_policy_id                = akamai_appsec_security_policy.policy1.security_policy_id
  apply_account_protection_controls = false
  apply_api_constraints             = true
  apply_application_layer_controls  = true
  apply_botman_controls             = true
  apply_malware_controls            = true
  apply_network_layer_controls      = true
  apply_rate_controls               = true
  apply_reputation_controls         = true
  apply_slow_post_controls          = true
  apply_url_protection_controls     = false
}

