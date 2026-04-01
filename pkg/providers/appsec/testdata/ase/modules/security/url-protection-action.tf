// URL Protection Actions
// URL Protection Actions for policy default_policy
resource "akamai_appsec_url_protection_action" "default_policy_50411" {
  config_id                 = local.config_id
  security_policy_id        = akamai_appsec_security_policy.default_policy.security_policy_id
  url_protection_policy_id  = 50411
  max_rate_threshold_action = "alert"
  load_shedding_action      = "deny"
}

resource "akamai_appsec_url_protection_action" "default_policy_50412" {
  config_id                 = local.config_id
  security_policy_id        = akamai_appsec_security_policy.default_policy.security_policy_id
  url_protection_policy_id  = 50412
  max_rate_threshold_action = "deny"
  load_shedding_action      = "none"
}

