resource "akamai_appsec_custom_deny" "deny_message_deny_custom_78842" {
  config_id = akamai_appsec_configuration.config.config_id
  custom_deny = jsonencode(
    {
      "name" : "Deny Message",
      "parameters" : [
        {
          "name" : "prevent_browser_cache",
          "value" : "true"
        },
        {
          "name" : "response_body_content",
          "value" : "You were denied, bad actor!"
        },
        {
          "name" : "response_content_type",
          "value" : "text/html"
        },
        {
          "name" : "response_status_code",
          "value" : "433"
        }
      ]
    }
  )
}

