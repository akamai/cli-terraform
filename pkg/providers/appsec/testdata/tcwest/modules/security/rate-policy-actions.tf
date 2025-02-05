// Rate Policy Actions
resource "akamai_appsec_rate_policy_action" "policy2_high_rate" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_rate_protection.policy2.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.high_rate.rate_policy_id
  ipv4_action        = "deny"
  ipv6_action        = "deny"
}

resource "akamai_appsec_rate_policy_action" "policy2_low_rate" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_rate_protection.policy2.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.low_rate.rate_policy_id
  ipv4_action        = "deny"
  ipv6_action        = "deny"
}

// Rate Policy Actions
resource "akamai_appsec_rate_policy_action" "policy1_high_rate" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_rate_protection.policy1.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.high_rate.rate_policy_id
  ipv4_action        = "deny_custom_78842"
  ipv6_action        = "deny_custom_78842"
}

resource "akamai_appsec_rate_policy_action" "policy1_low_rate" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_rate_protection.policy1.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.low_rate.rate_policy_id
  ipv4_action        = "alert"
  ipv6_action        = "alert"
}

