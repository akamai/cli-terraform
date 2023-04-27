resource "akamai_appsec_custom_rule" "custom_rule_1_60088542" {
  config_id = akamai_appsec_configuration.config.config_id
  custom_rule = jsonencode(
    {
      "conditions" : [
        {
          "positiveMatch" : true,
          "type" : "requestMethodMatch",
          "value" : [
            "POST"
          ]
        },
        {
          "positiveMatch" : true,
          "type" : "pathMatch",
          "value" : [
            "/login"
          ]
        },
        {
          "name" : "email",
          "positiveMatch" : true,
          "type" : "argsPostMatch",
          "value" : [
            "me@email.com"
          ]
        }
      ],
      "name" : "Custom Rule 1",
      "tag" : [
        "Login"
      ]
    }
  )
}

