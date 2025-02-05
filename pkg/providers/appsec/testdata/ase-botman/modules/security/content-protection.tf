resource "akamai_botman_content_protection_rule" "default_policy_new_rule_fakeba52-5c9e-4aa0-b5f7-5b88601d0d76" {
  config_id          = local.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  content_protection_rule = jsonencode(
    {
      "arrayKey" : [
        "arrayValueA1",
        "arrayValueA2"
      ],
      "contentProtectionRuleName" : "New Rule",
      "objectKey" : {
        "innerKey" : "innerValueA"
      },
      "primitiveKey" : "primitiveValueA"
    }
  )
}

resource "akamai_botman_content_protection_rule_sequence" "default_policy_sequence" {
  config_id          = local.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  content_protection_rule_ids = [
    akamai_botman_content_protection_rule.default_policy_new_rule_fakeba52-5c9e-4aa0-b5f7-5b88601d0d76.content_protection_rule_id
  ]
}

resource "akamai_botman_content_protection_javascript_injection_rule" "default_policy_new_injection_rule_fakeb37c-15ce-4ec8-ad99-0252d8a4580b" {
  config_id          = local.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  content_protection_javascript_injection_rule = jsonencode(
    {
      "arrayKey" : [
        "arrayValueA1",
        "arrayValueA2"
      ],
      "contentProtectionJavaScriptInjectionRuleName" : "New Injection Rule",
      "objectKey" : {
        "innerKey" : "innerValueA"
      },
      "primitiveKey" : "primitiveValueA"
    }
  )
}

