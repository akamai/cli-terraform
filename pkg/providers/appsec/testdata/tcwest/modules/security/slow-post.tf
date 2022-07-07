// Slow Post Protection
resource "akamai_appsec_slow_post" "policy2" {
  config_id                  = akamai_appsec_configuration.config.config_id
  security_policy_id         = akamai_appsec_slowpost_protection.policy2.security_policy_id
  slow_rate_action           = "alert"
  slow_rate_threshold_rate   = 10
  slow_rate_threshold_period = 60
}

// Slow Post Protection
resource "akamai_appsec_slow_post" "andrew" {
  config_id                  = akamai_appsec_configuration.config.config_id
  security_policy_id         = akamai_appsec_slowpost_protection.andrew.security_policy_id
  slow_rate_action           = "alert"
  slow_rate_threshold_rate   = 10
  slow_rate_threshold_period = 60
}

// Slow Post Protection
resource "akamai_appsec_slow_post" "policy1" {
  config_id                  = akamai_appsec_configuration.config.config_id
  security_policy_id         = akamai_appsec_slowpost_protection.policy1.security_policy_id
  slow_rate_action           = "abort"
  slow_rate_threshold_rate   = 10
  slow_rate_threshold_period = 60
  duration_threshold_timeout = 60
}

