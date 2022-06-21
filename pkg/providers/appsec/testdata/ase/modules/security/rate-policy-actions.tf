// Rate Policy Actions
resource "akamai_appsec_rate_policy_action" "default_policy_page_view_requests" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_rate_protection.default_policy.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.page_view_requests.rate_policy_id
  ipv4_action        = "alert"
  ipv6_action        = "alert"
}

resource "akamai_appsec_rate_policy_action" "default_policy_origin_error" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_rate_protection.default_policy.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.origin_error.rate_policy_id
  ipv4_action        = "alert"
  ipv6_action        = "alert"
}

resource "akamai_appsec_rate_policy_action" "default_policy_post_page_requests" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_rate_protection.default_policy.security_policy_id
  rate_policy_id     = akamai_appsec_rate_policy.post_page_requests.rate_policy_id
  ipv4_action        = "alert"
  ipv6_action        = "alert"
}

