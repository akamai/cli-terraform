// Enable/Disable Protections for policy default_policy
resource "akamai_appsec_waf_protection" "default_policy" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.default_policy.security_policy_id
  enabled            = true
}

resource "akamai_appsec_api_constraints_protection" "default_policy" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.default_policy.security_policy_id
  enabled            = false
}

resource "akamai_appsec_ip_geo_protection" "default_policy" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_api_constraints_protection.default_policy.security_policy_id
  enabled            = true
}

resource "akamai_appsec_malware_protection" "default_policy" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_ip_geo_protection.default_policy.security_policy_id
  enabled            = true
}

resource "akamai_appsec_rate_protection" "default_policy" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_malware_protection.default_policy.security_policy_id
  enabled            = true
}

resource "akamai_appsec_reputation_protection" "default_policy" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_rate_protection.default_policy.security_policy_id
  enabled            = false
}

resource "akamai_appsec_slowpost_protection" "default_policy" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_reputation_protection.default_policy.security_policy_id
  enabled            = true
}

resource "akamai_botman_bot_management_settings" "default_policy" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_slowpost_protection.default_policy.security_policy_id
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
