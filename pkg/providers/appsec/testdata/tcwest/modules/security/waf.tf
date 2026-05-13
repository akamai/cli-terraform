resource "akamai_appsec_waf_mode" "policy2" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.policy2.security_policy_id
  mode               = "KRS"
}

resource "akamai_appsec_waf_mode" "andrew" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.andrew.security_policy_id
  mode               = "KRS"
}

resource "akamai_appsec_waf_mode" "policy1" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.policy1.security_policy_id
  mode               = "KRS"
}


resource "akamai_appsec_custom_rule_action" "policy1_60088542" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.policy1.security_policy_id
  custom_rule_id     = akamai_appsec_custom_rule.custom_rule_1_60088542.custom_rule_id
  custom_rule_action = "deny"
}
