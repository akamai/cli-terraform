// Rapid Rules
resource "akamai_appsec_rapid_rules" "default_policy" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.default_policy.security_policy_id
  default_action     = "alert"
  rule_definitions = jsonencode(
    [
      {
        "id" : 3000101,
        "action" : "deny",
        "lock" : false
      },
      {
        "id" : 3000102,
        "action" : "alert",
        "lock" : true,
        "conditionException" : {
          "exception" : {
            "headerCookieOrParamValues" : [
              "exception-value"
            ],
            "specificHeaderCookieOrParamPrefix" : {
              "prefix" : "exception-prefix",
              "selector" : "REQUEST_HEADERS"
            },
            "specificHeaderCookieParamXmlOrJsonNames" : [
              {
                "names" : [
                  "Auth"
                ],
                "selector" : "REQUEST_HEADERS",
                "wildcard" : true
              }
            ]
          }
        }
      },
      {
        "id" : 3000103,
        "action" : "deny",
        "lock" : false,
        "conditionException" : {
          "advancedExceptions" : {
            "conditionOperator" : "AND",
            "specificHeaderCookieParamXmlOrJsonNames" : [
              {
                "names" : [
                  "X-Trace-Id"
                ],
                "selector" : "REQUEST_HEADERS"
              }
            ]
          }
        }
      }
    ]
  )
}

