resource "akamai_appsec_waf_mode" "default_policy" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.default_policy.security_policy_id
  mode               = "KRS"
}


