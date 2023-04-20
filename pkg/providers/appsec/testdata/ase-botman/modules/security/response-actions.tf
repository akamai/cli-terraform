resource "akamai_botman_serve_alternate_action" "serve_alternate_action_a_action_A" {
  config_id = akamai_appsec_configuration.config.config_id
  serve_alternate_action = jsonencode(
    {
      "actionName" : "Serve Alternate Action A",
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

resource "akamai_botman_serve_alternate_action" "serve_alternate_action_b_action_B" {
  config_id = akamai_appsec_configuration.config.config_id
  serve_alternate_action = jsonencode(
    {
      "actionName" : "Serve Alternate Action B",
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

resource "akamai_botman_challenge_action" "challenge_action_a_action_A" {
  config_id = akamai_appsec_configuration.config.config_id
  challenge_action = jsonencode(
    {
      "actionName" : "Challenge Action A",
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

resource "akamai_botman_challenge_action" "challenge_action_b_action_B" {
  config_id = akamai_appsec_configuration.config.config_id
  challenge_action = jsonencode(
    {
      "actionName" : "Challenge Action B",
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

resource "akamai_botman_conditional_action" "conditional_action_a_action_A" {
  config_id = akamai_appsec_configuration.config.config_id
  conditional_action = jsonencode(
    {
      "actionName" : "Conditional Action A",
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

resource "akamai_botman_conditional_action" "conditional_action_b_action_B" {
  config_id = akamai_appsec_configuration.config.config_id
  conditional_action = jsonencode(
    {
      "actionName" : "Conditional Action B",
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

resource "akamai_botman_challenge_interception_rules" "challenge_interception_rules" {
  config_id = akamai_appsec_configuration.config.config_id
  challenge_interception_rules = jsonencode(
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

