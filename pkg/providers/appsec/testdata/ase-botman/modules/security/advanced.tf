// Global Advanced
resource "akamai_appsec_advanced_settings_logging" "logging" {
  config_id = akamai_appsec_configuration.config.config_id
  logging = jsonencode(
    {
      "allowSampling" : true,
      "cookies" : {
        "type" : "all"
      },
      "customHeaders" : {
        "type" : "all"
      },
      "standardHeaders" : {
        "type" : "all"
      }
    }
  )
}

resource "akamai_appsec_advanced_settings_prefetch" "prefetch" {
  config_id            = akamai_appsec_configuration.config.config_id
  enable_app_layer     = true
  all_extensions       = false
  enable_rate_controls = false
  extensions           = ["cgi", "jsp", "aspx", "EMPTY_STRING", "php", "py", "asp"]
}

resource "akamai_appsec_advanced_settings_pragma_header" "pragma_header" {
  config_id = akamai_appsec_configuration.config.config_id
  pragma_header = jsonencode(
    {
      "action" : "REMOVE"
    }
  )
}

resource "akamai_appsec_advanced_settings_pii_learning" "pii_learning" {
  config_id           = akamai_appsec_configuration.config.config_id
  enable_pii_learning = true
}

resource "akamai_appsec_advanced_settings_attack_payload_logging" "attack_payload_logging" {
  config_id = akamai_appsec_configuration.config.config_id
  attack_payload_logging = jsonencode(
    {
      "enabled" : true,
      "requestBody" : {
        "type" : "NONE"
      },
      "responseBody" : {
        "type" : "ATTACK_PAYLOAD"
      }
    }
  )
}

resource "akamai_appsec_advanced_settings_request_body" "config_settings" {
  config_id                     = akamai_appsec_configuration.config.config_id
  request_body_inspection_limit = "default"
}

// RequestBody Overrides
resource "akamai_appsec_advanced_settings_request_body" "default_policy" {
  config_id                              = akamai_appsec_configuration.config.config_id
  security_policy_id                     = akamai_appsec_security_policy.default_policy.security_policy_id
  request_body_inspection_limit          = "default"
  request_body_inspection_limit_override = true
}
