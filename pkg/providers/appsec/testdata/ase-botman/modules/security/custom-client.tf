resource "akamai_botman_custom_client" "custom_client_a_a7fe489d-0354-43bd-b81c-8cabbe850cdd" {
  config_id = akamai_appsec_configuration.config.config_id
  custom_client = jsonencode(
    {
      "arrayKey" : [
        "arrayValueA1",
        "arrayValueA2"
      ],
      "customClientName" : "Custom Client A",
      "objectKey" : {
        "innerKey" : "innerValueA"
      },
      "primitiveKey" : "primitiveValueA"
    }
  )
}

resource "akamai_botman_custom_client" "custom_client_b_60374346-2d1d-444d-91c1-90373e3f804a" {
  config_id = akamai_appsec_configuration.config.config_id
  custom_client = jsonencode(
    {
      "arrayKey" : [
        "arrayValueB1",
        "arrayValueB2"
      ],
      "customClientName" : "Custom Client B",
      "objectKey" : {
        "innerKey" : "innerValueB"
      },
      "primitiveKey" : "primitiveValueB"
    }
  )
}

resource "akamai_botman_custom_client_sequence" "sequence" {
  config_id = akamai_appsec_configuration.config.config_id
  custom_client_ids = [
    akamai_botman_custom_client.custom_client_a_a7fe489d-0354-43bd-b81c-8cabbe850cdd.custom_client_id,
    akamai_botman_custom_client.custom_client_b_60374346-2d1d-444d-91c1-90373e3f804a.custom_client_id
  ]
}
