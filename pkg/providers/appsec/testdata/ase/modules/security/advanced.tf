// Global Advanced
resource "akamai_appsec_advanced_settings_logging" "logging" {
  config_id = local.config_id
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
  config_id            = local.config_id
  enable_app_layer     = true
  all_extensions       = false
  enable_rate_controls = false
  extensions           = ["cgi", "jsp", "aspx", "EMPTY_STRING", "php", "py", "asp"]
}

resource "akamai_appsec_advanced_settings_pragma_header" "pragma_header" {
  config_id = local.config_id
  pragma_header = jsonencode(
    {
      "action" : "REMOVE"
    }
  )
}

resource "akamai_appsec_advanced_settings_evasive_path_match" "evasive_path_match" {
  config_id         = local.config_id
  enable_path_match = false
}

resource "akamai_appsec_advanced_settings_url_evasion_defense" "url_evasion_defense" {
  depends_on   = [akamai_appsec_advanced_settings_evasive_path_match.evasive_path_match]
  config_id    = local.config_id
  status       = "enabled"
  bypass_lists = ["123456_BYPASSLIST"]
  rules = [
    {
      rule_id            = 3002500
      action             = "deny"
      condition_operator = "AND"
      conditions = [
        {
          type           = "hostMatch"
          hosts          = ["host.name"]
          positive_match = true
        },
        {
          type           = "extensionMatch"
          extensions     = ["extension"]
          positive_match = true
        },
        {
          type      = "filenameMatch"
          filenames = ["filename"]
        },
        {
          type           = "ipMatch"
          ips            = ["192.168.0.1"]
          positive_match = true
        },
        {
          type                 = "uriQueryMatch"
          name                 = "query"
          value                = "queryValue"
          positive_match       = true
          name_case_sensitive  = true
          value_case_sensitive = true
          value_wildcard       = true
        },
        {
          type                 = "requestHeaderMatch"
          header               = "Accept"
          value                = "XML"
          positive_match       = true
          value_case_sensitive = true
          value_wildcard       = true
        },
        {
          type           = "requestMethodMatch"
          methods        = ["DELETE"]
          positive_match = true
        },
        {
          type           = "pathMatch"
          paths          = ["/path"]
          positive_match = true
        },
        {
          type           = "clientListMatch"
          client_lists   = ["1234567_LIST"]
          positive_match = true
          use_headers    = true
        },
      ]
    },
    {
      rule_id = 3002501
      action  = "none"
    },
  ]
}

resource "akamai_appsec_advanced_settings_pii_learning" "pii_learning" {
  config_id           = local.config_id
  enable_pii_learning = false
}

resource "akamai_appsec_advanced_settings_attack_payload_logging" "attack_payload_logging" {
  config_id = local.config_id
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
  config_id                              = local.config_id
  request_body_inspection_limit          = "default"
  request_body_inspection_limit_override = false
}

// RequestBody Overrides
resource "akamai_appsec_advanced_settings_request_body" "default_policy" {
  config_id                              = local.config_id
  security_policy_id                     = akamai_appsec_security_policy.default_policy.security_policy_id
  request_body_inspection_limit          = "default"
  request_body_inspection_limit_override = true
}
