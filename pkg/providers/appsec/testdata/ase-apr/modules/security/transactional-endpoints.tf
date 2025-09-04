resource "akamai_botman_transactional_endpoint" "default_policy_061429d0-a709-418e-9311-1c3b4ee28792" {
  config_id          = akamai_botman_transactional_endpoint_protection.transactional_endpoint_protection.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  operation_id       = "061429d0-a709-418e-9311-1c3b4ee28792"
  transactional_endpoint = jsonencode(
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

resource "akamai_botman_transactional_endpoint" "default_policy_c2b20de9-bb2e-4da8-ac99-f5b6e9ae4e10" {
  config_id          = akamai_botman_transactional_endpoint_protection.transactional_endpoint_protection.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  operation_id       = "c2b20de9-bb2e-4da8-ac99-f5b6e9ae4e10"
  transactional_endpoint = jsonencode(
    {
      "arrayKey" : [
        "arrayValueB1",
        "arrayValueB2"
      ],
      "objectKey" : {
        "innerKey" : "innerValueB"
      },
      "primitiveKey" : "primitiveValueB"
    }
  )
}

resource "akamai_botman_bot_category_exception" "default_policy" {
  config_id          = akamai_botman_custom_bot_category_sequence.custom_bot_category_sequence.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  bot_category_exception = jsonencode(
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

