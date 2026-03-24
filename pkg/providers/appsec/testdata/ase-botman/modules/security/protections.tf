// Enable/Disable Protections for policy default_policy
resource "akamai_appsec_security_policy_protections" "default_policy" {
  config_id                         = local.config_id
  security_policy_id                = akamai_appsec_security_policy.default_policy.security_policy_id
  apply_account_protection_controls = false
  apply_api_constraints             = false
  apply_application_layer_controls  = true
  apply_botman_controls             = false
  apply_malware_controls            = true
  apply_network_layer_controls      = true
  apply_rate_controls               = true
  apply_reputation_controls         = false
  apply_slow_post_controls          = true
  apply_url_protection_controls     = false
}

resource "akamai_botman_bot_management_settings" "default_policy" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy_protections.default_policy.security_policy_id
  bot_management_settings = jsonencode(
    {
      "arrayKey" : [
        "arrayValueA1",
        "arrayValueA2"
      ],
      "objectKey" : {
        "innerKey" : "innerValueA"
      },
      "primitiveKey" : "primitiveValueA"
    }
  )
}
