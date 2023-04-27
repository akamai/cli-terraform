resource "akamai_botman_custom_bot_category" "category_a_dae597b8-b552-4c95-ab8b-066a3fef2f75" {
  config_id = akamai_appsec_configuration.config.config_id
  custom_bot_category = jsonencode(
    {
      "arrayKey" : [
        "arrayValueA1",
        "arrayValueA2"
      ],
      "categoryName" : "category a",
      "objectKey" : {
        "innerKey" : "innerValueA"
      },
      "primitiveKey" : "primitiveValueA"
    }
  )
}

resource "akamai_botman_recategorized_akamai_defined_bot" "akamai_defined_bot_a_eceac3f9-871b-4c57-9a24-c25b0237949a" {
  config_id   = akamai_appsec_configuration.config.config_id
  bot_id      = "eceac3f9-871b-4c57-9a24-c25b0237949a"
  category_id = akamai_botman_custom_bot_category.category_a_dae597b8-b552-4c95-ab8b-066a3fef2f75.category_id
}

resource "akamai_botman_recategorized_akamai_defined_bot" "akamai_defined_bot_b_c590d2e5-a041-4f05-8fda-71608f42d720" {
  config_id   = akamai_appsec_configuration.config.config_id
  bot_id      = "c590d2e5-a041-4f05-8fda-71608f42d720"
  category_id = akamai_botman_custom_bot_category.category_a_dae597b8-b552-4c95-ab8b-066a3fef2f75.category_id
}

resource "akamai_botman_custom_bot_category" "category_b_c3362371-4b98-40fe-a7b9-cd7fab93eec5" {
  config_id = akamai_appsec_configuration.config.config_id
  custom_bot_category = jsonencode(
    {
      "arrayKey" : [
        "arrayValueB1",
        "arrayValueB2"
      ],
      "categoryName" : "category b",
      "objectKey" : {
        "innerKey" : "innerValueB"
      },
      "primitiveKey" : "primitiveValueB"
    }
  )
}

resource "akamai_botman_custom_bot_category_sequence" "custom_bot_category_sequence" {
  config_id    = akamai_appsec_configuration.config.config_id
  category_ids = [akamai_botman_custom_bot_category.category_b_c3362371-4b98-40fe-a7b9-cd7fab93eec5.category_id, akamai_botman_custom_bot_category.category_a_dae597b8-b552-4c95-ab8b-066a3fef2f75.category_id]
}

resource "akamai_botman_custom_defined_bot" "bot_a_50789280-ba99-4f8f-b4c6-ad9c1c69569a" {
  config_id = akamai_botman_custom_bot_category_sequence.custom_bot_category_sequence.config_id
  custom_defined_bot = jsonencode(
    {
      "arrayKey" : [
        "arrayValueA1",
        "arrayValueA2"
      ],
      "botName" : "Bot A",
      "objectKey" : {
        "innerKey" : "innerValueA"
      },
      "primitiveKey" : "primitiveValueA"
    }
  )
}

resource "akamai_botman_custom_defined_bot" "bot_b_da1de35e-deda-4273-933d-3131291fa3d4" {
  config_id = akamai_botman_custom_bot_category_sequence.custom_bot_category_sequence.config_id
  custom_defined_bot = jsonencode(
    {
      "arrayKey" : [
        "arrayValueB1",
        "arrayValueB2"
      ],
      "botName" : "Bot B",
      "objectKey" : {
        "innerKey" : "innerValueB"
      },
      "primitiveKey" : "primitiveValueB"
    }
  )
}

