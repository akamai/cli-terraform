resource "akamai_botman_javascript_injection" "default_policy" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  javascript_injection = jsonencode(
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

