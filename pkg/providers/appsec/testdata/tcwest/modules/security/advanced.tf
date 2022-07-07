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
  extensions           = ["cgi", "jsp", "EMPTY_STRING", "aspx", "py", "php", "asp"]
}

resource "akamai_appsec_advanced_settings_pragma_header" "pragma_header" {
  config_id = akamai_appsec_configuration.config.config_id
  pragma_header = jsonencode(
    {
      "action" : "REMOVE",
      "conditionOperator" : "OR",
      "excludeCondition" : [
        {
          "type" : "ipMatch",
          "positiveMatch" : true,
          "header" : "",
          "value" : [
            "3.3.3.3"
          ],
          "name" : "",
          "valueCase" : false,
          "valueWildcard" : false,
          "useHeaders" : true
        }
      ]
    }
  )
}

// Logging Overides
resource "akamai_appsec_advanced_settings_logging" "policy1" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.policy1.security_policy_id
  logging = jsonencode(
    {
      "allowSampling" : true,
      "cookies" : {
        "type" : "all"
      },
      "customHeaders" : {
        "type" : "all"
      },
      "override" : true,
      "standardHeaders" : {
        "type" : "exclude",
        "values" : [
          "Accept-Charset"
        ]
      }
    }
  )
}
