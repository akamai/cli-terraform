// Enable/Disable Protections for policy policy2
resource "akamai_appsec_waf_protection" "policy2" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.policy2.security_policy_id
  enabled            = true
}

resource "akamai_appsec_api_constraints_protection" "policy2" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  enabled            = true
}

resource "akamai_appsec_ip_geo_protection" "policy2" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_api_constraints_protection.policy2.security_policy_id
  enabled            = true
}

resource "akamai_appsec_malware_protection" "policy2" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_ip_geo_protection.policy2.security_policy_id
  enabled            = true
}

resource "akamai_appsec_rate_protection" "policy2" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_malware_protection.policy2.security_policy_id
  enabled            = true
}

resource "akamai_appsec_reputation_protection" "policy2" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_rate_protection.policy2.security_policy_id
  enabled            = true
}

resource "akamai_appsec_slowpost_protection" "policy2" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_reputation_protection.policy2.security_policy_id
  enabled            = true
}

// Enable/Disable Protections for policy andrew
resource "akamai_appsec_waf_protection" "andrew" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.andrew.security_policy_id
  enabled            = true
}

resource "akamai_appsec_api_constraints_protection" "andrew" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  enabled            = true
}

resource "akamai_appsec_ip_geo_protection" "andrew" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_api_constraints_protection.andrew.security_policy_id
  enabled            = true
}

resource "akamai_appsec_malware_protection" "andrew" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_ip_geo_protection.andrew.security_policy_id
  enabled            = true
}

resource "akamai_appsec_rate_protection" "andrew" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_malware_protection.andrew.security_policy_id
  enabled            = true
}

resource "akamai_appsec_reputation_protection" "andrew" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_rate_protection.andrew.security_policy_id
  enabled            = true
}

resource "akamai_appsec_slowpost_protection" "andrew" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_reputation_protection.andrew.security_policy_id
  enabled            = true
}

// Enable/Disable Protections for policy policy1
resource "akamai_appsec_waf_protection" "policy1" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.policy1.security_policy_id
  enabled            = true
}

resource "akamai_appsec_api_constraints_protection" "policy1" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  enabled            = true
}

resource "akamai_appsec_ip_geo_protection" "policy1" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_api_constraints_protection.policy1.security_policy_id
  enabled            = true
}

resource "akamai_appsec_malware_protection" "policy1" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_ip_geo_protection.policy1.security_policy_id
  enabled            = true
}

resource "akamai_appsec_rate_protection" "policy1" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_malware_protection.policy1.security_policy_id
  enabled            = true
}

resource "akamai_appsec_reputation_protection" "policy1" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_rate_protection.policy1.security_policy_id
  enabled            = true
}

resource "akamai_appsec_slowpost_protection" "policy1" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_reputation_protection.policy1.security_policy_id
  enabled            = true
}

