resource "akamai_botman_akamai_bot_category_action" "default_policy_akamai_bot_category_a_0b116152-1d20-4715-8fa7-dcacb1c697e2" {
  config_id          = local.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  category_id        = "0b116152-1d20-4715-8fa7-dcacb1c697e2"
  akamai_bot_category_action = jsonencode(
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

resource "akamai_botman_akamai_bot_category_action" "default_policy_akamai_bot_category_b_da0596ba-2379-4657-9b84-79b460d66070" {
  config_id          = local.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  category_id        = "da0596ba-2379-4657-9b84-79b460d66070"
  akamai_bot_category_action = jsonencode(
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

resource "akamai_botman_custom_bot_category_action" "default_policy_category_a_dae597b8-b552-4c95-ab8b-066a3fef2f75" {
  config_id          = local.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  category_id        = akamai_botman_custom_bot_category.category_a_dae597b8-b552-4c95-ab8b-066a3fef2f75.category_id
  custom_bot_category_action = jsonencode(
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

resource "akamai_botman_custom_bot_category_action" "default_policy_category_b_c3362371-4b98-40fe-a7b9-cd7fab93eec5" {
  config_id          = local.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  category_id        = akamai_botman_custom_bot_category.category_b_c3362371-4b98-40fe-a7b9-cd7fab93eec5.category_id
  custom_bot_category_action = jsonencode(
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

resource "akamai_botman_bot_detection_action" "default_policy_bot_detection_a_179e6bd6-5077-4f22-9a5b-3b09ee731eca" {
  config_id          = local.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  detection_id       = "179e6bd6-5077-4f22-9a5b-3b09ee731eca"
  bot_detection_action = jsonencode(
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

resource "akamai_botman_bot_detection_action" "default_policy_bot_detection_b_c4d20de1-af7a-476f-911d-73aedd97e294" {
  config_id          = local.config_id
  security_policy_id = akamai_botman_bot_management_settings.default_policy.security_policy_id
  detection_id       = "c4d20de1-af7a-476f-911d-73aedd97e294"
  bot_detection_action = jsonencode(
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

