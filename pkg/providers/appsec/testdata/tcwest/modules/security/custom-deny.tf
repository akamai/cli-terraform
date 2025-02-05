resource "akamai_appsec_custom_deny" "deny_message_deny_custom_78842" {
  config_id = local.config_id
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

resource "akamai_appsec_custom_deny" "deny_message_2_deny_custom_80270" {
  config_id = local.config_id
  custom_deny = jsonencode(
    {
      "name" : "deny message 2",
      "parameters" : [
        {
          "name" : "prevent_browser_cache",
          "value" : "true"
        },
        {
          "name" : "response_body_content",
          "value" : "bad actor 2"
        },
        {
          "name" : "response_content_type",
          "value" : "text/html"
        },
        {
          "name" : "response_status_code",
          "value" : "403"
        }
      ]
    }
  )
}

