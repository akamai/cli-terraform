resource "akamai_botman_bot_analytics_cookie" "bot_analytics_cookie" {
  config_id = local.config_id
  bot_analytics_cookie = jsonencode(
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

resource "akamai_botman_client_side_security" "client_side_security" {
  config_id = local.config_id
  client_side_security = jsonencode(
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

resource "akamai_botman_transactional_endpoint_protection" "transactional_endpoint_protection" {
  config_id = local.config_id
  transactional_endpoint_protection = jsonencode(
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

resource "akamai_apr_user_risk_response_strategy" "user_risk_response_strategy" {
  config_id = local.config_id
  user_risk_response_strategy = jsonencode(
    {
      "traffic" : {
        "inline" : {
          "aggressive" : {
            "threshold" : 76
          },
          "cautious" : {
            "threshold" : 0
          },
          "strict" : {
            "threshold" : 51
          }
        },
        "nativeSdkAndroid" : {
          "aggressive" : {
            "threshold" : 76
          },
          "cautious" : {
            "threshold" : 0
          },
          "strict" : {
            "threshold" : 51
          }
        },
        "nativeSdkIos" : {
          "aggressive" : {
            "threshold" : 76
          },
          "cautious" : {
            "threshold" : 0
          },
          "strict" : {
            "threshold" : 51
          }
        },
        "standard" : {
          "aggressive" : {
            "threshold" : 76
          },
          "cautious" : {
            "threshold" : 0
          },
          "strict" : {
            "threshold" : 51
          }
        }
      }
    }
  )
}
resource "akamai_apr_user_allow_list" "user_allow_list" {
  config_id = local.config_id
  user_allow_list = jsonencode(
    {
      "userAllowListId" : "82004_APRE2ENLDONOTDELETE"
    }
  )
}

