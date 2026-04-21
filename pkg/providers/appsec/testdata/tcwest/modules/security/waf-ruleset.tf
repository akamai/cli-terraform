resource "akamai_appsec_waf_ruleset" "policy2" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.policy2.security_policy_id
  rules = [
    {
      rule_id     = "699989"
      rule_action = "alert"
      condition_exception = jsonencode(merge(
        { "exception" : {
          "specificHeaderCookieParamXmlOrJsonNames" : [
            {
              "names" : [
                "value",
                "body%5B0%5D%5Bvalue%5D"
              ],
              "selector" : "ARGS"
            }
          ]
        } },
        { "conditions" : [
          {
            "type" : "requestHeaderMatch",
            "name" : "User-Agent",
            "positiveMatch" : true,
            "value" : "*Slackbot-LinkExpanding*",
            "valueWildcard" : true
          }
        ] }
      ))
    },
    {
      rule_id     = "699996"
      rule_action = "alert"
    },
    {
      rule_id     = "950000"
      rule_action = "alert"
    },
    {
      rule_id     = "950001"
      rule_action = "alert"
    },
    {
      rule_id     = "950002"
      rule_action = "alert"
    },
    {
      rule_id     = "950003"
      rule_action = "alert"
    },
    {
      rule_id     = "950005"
      rule_action = "alert"
    },
    {
      rule_id     = "950006"
      rule_action = "alert"
    },
    {
      rule_id     = "950007"
      rule_action = "alert"
    },
    {
      rule_id     = "950008"
      rule_action = "alert"
    },
    {
      rule_id     = "950009"
      rule_action = "alert"
    },
    {
      rule_id     = "950010"
      rule_action = "alert"
    },
    {
      rule_id     = "950011"
      rule_action = "alert"
    },
    {
      rule_id     = "950018"
      rule_action = "alert"
    },
    {
      rule_id     = "950019"
      rule_action = "alert"
    },
    {
      rule_id     = "950103"
      rule_action = "alert"
    },
    {
      rule_id     = "950107"
      rule_action = "alert"
    },
    {
      rule_id     = "950108"
      rule_action = "alert"
    },
    {
      rule_id     = "950109"
      rule_action = "alert"
    },
    {
      rule_id     = "950110"
      rule_action = "alert"
    },
    {
      rule_id     = "950116"
      rule_action = "alert"
    },
    {
      rule_id     = "950117"
      rule_action = "alert"
    },
    {
      rule_id     = "950118"
      rule_action = "alert"
    },
    {
      rule_id     = "950119"
      rule_action = "alert"
    },
    {
      rule_id     = "950120"
      rule_action = "alert"
    },
    {
      rule_id     = "950901"
      rule_action = "alert"
    },
    {
      rule_id     = "950908"
      rule_action = "alert"
    },
    {
      rule_id     = "950910"
      rule_action = "alert"
    },
    {
      rule_id     = "950911"
      rule_action = "alert"
    },
    {
      rule_id     = "950921"
      rule_action = "alert"
    },
    {
      rule_id     = "958000"
      rule_action = "alert"
    },
    {
      rule_id     = "958001"
      rule_action = "alert"
    },
    {
      rule_id     = "958002"
      rule_action = "alert"
    },
    {
      rule_id     = "958003"
      rule_action = "alert"
    },
    {
      rule_id     = "958004"
      rule_action = "alert"
    },
    {
      rule_id     = "958005"
      rule_action = "alert"
    },
    {
      rule_id     = "958006"
      rule_action = "alert"
    },
    {
      rule_id     = "958007"
      rule_action = "alert"
    },
    {
      rule_id     = "958008"
      rule_action = "alert"
    },
    {
      rule_id     = "958009"
      rule_action = "alert"
    },
    {
      rule_id     = "958010"
      rule_action = "alert"
    },
    {
      rule_id     = "958011"
      rule_action = "alert"
    },
    {
      rule_id     = "958012"
      rule_action = "alert"
    },
    {
      rule_id     = "958013"
      rule_action = "alert"
    },
    {
      rule_id     = "958016"
      rule_action = "alert"
    },
    {
      rule_id     = "958017"
      rule_action = "alert"
    },
    {
      rule_id     = "958018"
      rule_action = "alert"
    },
    {
      rule_id     = "958019"
      rule_action = "alert"
    },
    {
      rule_id     = "958020"
      rule_action = "alert"
    },
    {
      rule_id     = "958022"
      rule_action = "alert"
    },
    {
      rule_id     = "958023"
      rule_action = "alert"
    },
    {
      rule_id     = "958024"
      rule_action = "alert"
    },
    {
      rule_id     = "958025"
      rule_action = "alert"
    },
    {
      rule_id     = "958026"
      rule_action = "alert"
    },
    {
      rule_id     = "958027"
      rule_action = "alert"
    },
    {
      rule_id     = "958028"
      rule_action = "alert"
    },
    {
      rule_id     = "958030"
      rule_action = "alert"
    },
    {
      rule_id     = "958031"
      rule_action = "alert"
    },
    {
      rule_id     = "958032"
      rule_action = "alert"
    },
    {
      rule_id     = "958033"
      rule_action = "alert"
    },
    {
      rule_id     = "958034"
      rule_action = "alert"
    },
    {
      rule_id     = "958036"
      rule_action = "alert"
    },
    {
      rule_id     = "958037"
      rule_action = "alert"
    },
    {
      rule_id     = "958038"
      rule_action = "alert"
    },
    {
      rule_id     = "958039"
      rule_action = "alert"
    },
    {
      rule_id     = "958040"
      rule_action = "alert"
    },
    {
      rule_id     = "958041"
      rule_action = "alert"
    },
    {
      rule_id     = "958045"
      rule_action = "alert"
    },
    {
      rule_id     = "958046"
      rule_action = "alert"
    },
    {
      rule_id     = "958047"
      rule_action = "alert"
    },
    {
      rule_id     = "958049"
      rule_action = "alert"
    },
    {
      rule_id     = "958051"
      rule_action = "alert"
    },
    {
      rule_id     = "958052"
      rule_action = "alert"
    },
    {
      rule_id     = "958054"
      rule_action = "alert"
    },
    {
      rule_id     = "958056"
      rule_action = "alert"
    },
    {
      rule_id     = "958057"
      rule_action = "alert"
    },
    {
      rule_id     = "958059"
      rule_action = "alert"
    },
    {
      rule_id     = "958230"
      rule_action = "alert"
    },
    {
      rule_id     = "958231"
      rule_action = "alert"
    },
    {
      rule_id     = "958291"
      rule_action = "alert"
    },
    {
      rule_id     = "958295"
      rule_action = "alert"
    },
    {
      rule_id     = "958404"
      rule_action = "alert"
    },
    {
      rule_id     = "958405"
      rule_action = "alert"
    },
    {
      rule_id     = "958406"
      rule_action = "alert"
    },
    {
      rule_id     = "958407"
      rule_action = "alert"
    },
    {
      rule_id     = "958408"
      rule_action = "alert"
    },
    {
      rule_id     = "958409"
      rule_action = "alert"
    },
    {
      rule_id     = "958410"
      rule_action = "alert"
    },
    {
      rule_id     = "958411"
      rule_action = "alert"
    },
    {
      rule_id     = "958412"
      rule_action = "alert"
    },
    {
      rule_id     = "958413"
      rule_action = "alert"
    },
    {
      rule_id     = "958414"
      rule_action = "alert"
    },
    {
      rule_id     = "958415"
      rule_action = "alert"
    },
    {
      rule_id     = "958416"
      rule_action = "alert"
    },
    {
      rule_id     = "958417"
      rule_action = "alert"
    },
    {
      rule_id     = "958418"
      rule_action = "alert"
    },
    {
      rule_id     = "958419"
      rule_action = "alert"
    },
    {
      rule_id     = "958420"
      rule_action = "alert"
    },
    {
      rule_id     = "958421"
      rule_action = "alert"
    },
    {
      rule_id     = "958422"
      rule_action = "alert"
    },
    {
      rule_id     = "958423"
      rule_action = "alert"
    },
    {
      rule_id     = "958976"
      rule_action = "alert"
    },
    {
      rule_id     = "958977"
      rule_action = "alert"
    },
    {
      rule_id     = "959070"
      rule_action = "alert"
    },
    {
      rule_id     = "959071"
      rule_action = "alert"
    },
    {
      rule_id     = "959072"
      rule_action = "alert"
    },
    {
      rule_id     = "959073"
      rule_action = "alert"
    },
    {
      rule_id     = "959151"
      rule_action = "alert"
    },
    {
      rule_id     = "960010"
      rule_action = "alert"
    },
    {
      rule_id     = "960011"
      rule_action = "alert"
    },
    {
      rule_id     = "960012"
      rule_action = "alert"
    },
    {
      rule_id     = "960016"
      rule_action = "alert"
    },
    {
      rule_id     = "960022"
      rule_action = "alert"
    },
    {
      rule_id     = "960034"
      rule_action = "alert"
    },
    {
      rule_id     = "960035"
      rule_action = "alert"
    },
    {
      rule_id     = "960208"
      rule_action = "alert"
    },
    {
      rule_id     = "960209"
      rule_action = "alert"
    },
    {
      rule_id     = "960335"
      rule_action = "alert"
    },
    {
      rule_id     = "960341"
      rule_action = "alert"
    },
    {
      rule_id     = "960901"
      rule_action = "alert"
    },
    {
      rule_id     = "960902"
      rule_action = "alert"
    },
    {
      rule_id     = "960904"
      rule_action = "alert"
    },
    {
      rule_id     = "960912"
      rule_action = "alert"
    },
    {
      rule_id     = "960913"
      rule_action = "alert"
    },
    {
      rule_id     = "960914"
      rule_action = "alert"
    },
    {
      rule_id     = "973300"
      rule_action = "alert"
    },
    {
      rule_id     = "973301"
      rule_action = "alert"
    },
    {
      rule_id     = "973302"
      rule_action = "alert"
    },
    {
      rule_id     = "973303"
      rule_action = "alert"
    },
    {
      rule_id     = "973304"
      rule_action = "alert"
    },
    {
      rule_id     = "973305"
      rule_action = "alert"
    },
    {
      rule_id     = "973306"
      rule_action = "alert"
    },
    {
      rule_id     = "973307"
      rule_action = "alert"
    },
    {
      rule_id     = "973308"
      rule_action = "alert"
    },
    {
      rule_id     = "973309"
      rule_action = "alert"
    },
    {
      rule_id     = "973310"
      rule_action = "alert"
    },
    {
      rule_id     = "973311"
      rule_action = "alert"
    },
    {
      rule_id     = "973312"
      rule_action = "alert"
    },
    {
      rule_id     = "973313"
      rule_action = "alert"
    },
    {
      rule_id     = "973314"
      rule_action = "alert"
    },
    {
      rule_id     = "973315"
      rule_action = "alert"
    },
    {
      rule_id     = "973316"
      rule_action = "alert"
    },
    {
      rule_id     = "973317"
      rule_action = "alert"
    },
    {
      rule_id     = "973318"
      rule_action = "alert"
    },
    {
      rule_id     = "973319"
      rule_action = "alert"
    },
    {
      rule_id     = "973320"
      rule_action = "alert"
    },
    {
      rule_id     = "973321"
      rule_action = "alert"
    },
    {
      rule_id     = "973322"
      rule_action = "alert"
    },
    {
      rule_id     = "973323"
      rule_action = "alert"
    },
    {
      rule_id     = "973324"
      rule_action = "alert"
    },
    {
      rule_id     = "973325"
      rule_action = "alert"
    },
    {
      rule_id     = "973326"
      rule_action = "alert"
    },
    {
      rule_id     = "973327"
      rule_action = "alert"
    },
    {
      rule_id     = "973328"
      rule_action = "alert"
    },
    {
      rule_id     = "973329"
      rule_action = "alert"
    },
    {
      rule_id     = "973330"
      rule_action = "alert"
    },
    {
      rule_id     = "973331"
      rule_action = "alert"
    },
    {
      rule_id     = "973332"
      rule_action = "alert"
    },
    {
      rule_id     = "973333"
      rule_action = "alert"
    },
    {
      rule_id     = "973334"
      rule_action = "alert"
    },
    {
      rule_id     = "973335"
      rule_action = "alert"
    },
    {
      rule_id     = "973336"
      rule_action = "alert"
    },
    {
      rule_id     = "973337"
      rule_action = "alert"
    },
    {
      rule_id     = "981173"
      rule_action = "alert"
    },
    {
      rule_id     = "981241"
      rule_action = "alert"
    },
    {
      rule_id     = "981242"
      rule_action = "alert"
    },
    {
      rule_id     = "981243"
      rule_action = "alert"
    },
    {
      rule_id     = "981244"
      rule_action = "alert"
    },
    {
      rule_id     = "981245"
      rule_action = "alert"
    },
    {
      rule_id     = "981246"
      rule_action = "alert"
    },
    {
      rule_id     = "981247"
      rule_action = "alert"
    },
    {
      rule_id     = "981248"
      rule_action = "alert"
    },
    {
      rule_id     = "981249"
      rule_action = "alert"
    },
    {
      rule_id     = "981250"
      rule_action = "alert"
    },
    {
      rule_id     = "981251"
      rule_action = "alert"
    },
    {
      rule_id     = "981252"
      rule_action = "alert"
    },
    {
      rule_id     = "981253"
      rule_action = "alert"
    },
    {
      rule_id     = "981254"
      rule_action = "alert"
    },
    {
      rule_id     = "981255"
      rule_action = "alert"
    },
    {
      rule_id     = "981256"
      rule_action = "alert"
    },
    {
      rule_id     = "981260"
      rule_action = "alert"
    },
    {
      rule_id     = "981270"
      rule_action = "alert"
    },
    {
      rule_id     = "981272"
      rule_action = "alert"
    },
    {
      rule_id     = "981276"
      rule_action = "alert"
    },
    {
      rule_id     = "981277"
      rule_action = "alert"
    },
    {
      rule_id     = "981300"
      rule_action = "alert"
    },
    {
      rule_id     = "981318"
      rule_action = "alert"
    },
    {
      rule_id     = "981319"
      rule_action = "alert"
    },
    {
      rule_id     = "981320"
      rule_action = "alert"
    },
    {
      rule_id     = "990002"
      rule_action = "alert"
    },
    {
      rule_id     = "990012"
      rule_action = "alert"
    },
    {
      rule_id     = "990901"
      rule_action = "alert"
    },
    {
      rule_id     = "990902"
      rule_action = "alert"
    },
    {
      rule_id     = "3000000"
      rule_action = "alert"
    },
    {
      rule_id     = "3000001"
      rule_action = "alert"
    },
    {
      rule_id     = "3000002"
      rule_action = "alert"
    },
    {
      rule_id     = "3000003"
      rule_action = "alert"
    },
    {
      rule_id     = "3000004"
      rule_action = "alert"
    },
    {
      rule_id     = "3000005"
      rule_action = "alert"
    },
    {
      rule_id     = "3000006"
      rule_action = "alert"
    },
    {
      rule_id     = "3000007"
      rule_action = "alert"
    },
    {
      rule_id     = "3000008"
      rule_action = "alert"
    },
    {
      rule_id     = "3000009"
      rule_action = "alert"
    },
    {
      rule_id     = "3000010"
      rule_action = "alert"
    },
    {
      rule_id     = "3000011"
      rule_action = "alert"
    },
    {
      rule_id     = "3000012"
      rule_action = "alert"
    },
    {
      rule_id     = "3000013"
      rule_action = "alert"
    },
    {
      rule_id     = "3000014"
      rule_action = "alert"
    },
    {
      rule_id     = "3000015"
      rule_action = "alert"
    },
    {
      rule_id     = "3000016"
      rule_action = "alert"
    },
    {
      rule_id     = "3000017"
      rule_action = "alert"
    },
    {
      rule_id     = "3000018"
      rule_action = "alert"
    },
    {
      rule_id     = "3000019"
      rule_action = "alert"
    },
    {
      rule_id     = "3000020"
      rule_action = "alert"
    },
    {
      rule_id     = "3000021"
      rule_action = "alert"
    },
    {
      rule_id     = "3000022"
      rule_action = "alert"
    },
    {
      rule_id     = "3000023"
      rule_action = "alert"
    },
    {
      rule_id     = "3000024"
      rule_action = "alert"
    },
    {
      rule_id     = "3000025"
      rule_action = "alert"
    },
    {
      rule_id     = "3000027"
      rule_action = "alert"
    },
    {
      rule_id     = "3000029"
      rule_action = "alert"
    },
    {
      rule_id     = "3000030"
      rule_action = "alert"
    },
    {
      rule_id     = "3000031"
      rule_action = "alert"
    },
    {
      rule_id     = "3000032"
      rule_action = "alert"
    },
    {
      rule_id     = "3000033"
      rule_action = "alert"
    },
    {
      rule_id     = "3000034"
      rule_action = "alert"
    },
    {
      rule_id     = "3000035"
      rule_action = "alert"
    },
    {
      rule_id     = "3000036"
      rule_action = "alert"
    },
    {
      rule_id     = "3000037"
      rule_action = "alert"
    },
    {
      rule_id     = "3000038"
      rule_action = "alert"
    },
    {
      rule_id     = "3000039"
      rule_action = "alert"
    },
    {
      rule_id     = "3000040"
      rule_action = "alert"
    },
    {
      rule_id     = "3000041"
      rule_action = "alert"
    },
    {
      rule_id     = "3000042"
      rule_action = "alert"
    },
    {
      rule_id     = "3000043"
      rule_action = "alert"
    },
    {
      rule_id     = "3000044"
      rule_action = "alert"
    },
    {
      rule_id     = "3000045"
      rule_action = "alert"
    },
    {
      rule_id     = "3000046"
      rule_action = "alert"
    },
    {
      rule_id     = "3000047"
      rule_action = "alert"
    },
    {
      rule_id     = "3000048"
      rule_action = "alert"
    },
    {
      rule_id     = "3000049"
      rule_action = "alert"
    },
    {
      rule_id     = "3000050"
      rule_action = "alert"
    },
    {
      rule_id     = "3000051"
      rule_action = "alert"
    },
    {
      rule_id     = "3000052"
      rule_action = "alert"
    },
    {
      rule_id     = "3000053"
      rule_action = "alert"
    },
    {
      rule_id     = "3000054"
      rule_action = "alert"
    },
    {
      rule_id     = "3000055"
      rule_action = "alert"
    },
    {
      rule_id     = "3000056"
      rule_action = "alert"
    },
    {
      rule_id     = "3000057"
      rule_action = "alert"
    },
    {
      rule_id     = "3000058"
      rule_action = "alert"
    },
    {
      rule_id     = "3000059"
      rule_action = "alert"
    },
    {
      rule_id     = "3000060"
      rule_action = "alert"
    },
    {
      rule_id     = "3000061"
      rule_action = "alert"
    },
    {
      rule_id     = "3000062"
      rule_action = "alert"
    },
    {
      rule_id     = "3000063"
      rule_action = "alert"
    },
    {
      rule_id     = "3000064"
      rule_action = "alert"
    },
    {
      rule_id     = "3000065"
      rule_action = "alert"
    },
    {
      rule_id     = "3000066"
      rule_action = "alert"
    },
    {
      rule_id     = "3000067"
      rule_action = "alert"
    },
    {
      rule_id     = "3000068"
      rule_action = "alert"
    },
    {
      rule_id     = "3000071"
      rule_action = "alert"
    },
    {
      rule_id     = "3000072"
      rule_action = "alert"
    },
    {
      rule_id     = "3000073"
      rule_action = "alert"
    },
    {
      rule_id     = "3000074"
      rule_action = "alert"
    },
    {
      rule_id     = "3000075"
      rule_action = "alert"
    },
    {
      rule_id     = "3000076"
      rule_action = "alert"
    },
    {
      rule_id     = "3000077"
      rule_action = "alert"
    },
    {
      rule_id     = "3000079"
      rule_action = "alert"
    },
    {
      rule_id     = "3000080"
      rule_action = "alert"
    },
    {
      rule_id     = "3000081"
      rule_action = "alert"
    },
    {
      rule_id     = "3000082"
      rule_action = "alert"
    },
    {
      rule_id     = "3000083"
      rule_action = "alert"
    },
    {
      rule_id     = "3000084"
      rule_action = "alert"
    },
    {
      rule_id     = "3000999"
      rule_action = "alert"
    },
  ]
  attack_groups = [
    {
      attack_group        = "SQL"
      attack_group_action = "alert"
    },
    {
      attack_group        = "XSS"
      attack_group_action = "deny"
      condition_exception = jsonencode(merge(
        { "advancedExceptions" : {
          "conditionOperator" : "AND",
          "conditions" : [
            {
              "type" : "extensionMatch",
              "extensions" : [
                "pdf"
              ],
              "positiveMatch" : true
            },
            {
              "type" : "filenameMatch",
              "filenames" : [
                "upload.txt"
              ],
              "positiveMatch" : true
            },
            {
              "type" : "pathMatch",
              "paths" : [
                "/test"
              ],
              "positiveMatch" : true
            }
          ],
          "specificHeaderCookieOrParamNameValue" : [
            {
              "namesValues" : [
                {
                  "names" : [
                    "test"
                  ],
                  "values" : [
                    "test"
                  ]
                }
              ],
              "selector" : "ARGS",
              "valueWildcard" : true,
              "wildcard" : true
            }
          ],
          "specificHeaderCookieParamXmlOrJsonNames" : [
            {
              "names" : [
                "testgroup*"
              ],
              "selector" : "ARGS_NAMES",
              "wildcard" : true
            },
            {
              "names" : [
                "group*"
              ],
              "selector" : "ARGS",
              "wildcard" : true
            }
          ]
        } }
      ))
    },
    {
      attack_group        = "CMD"
      attack_group_action = "alert"
    },
    {
      attack_group        = "HTTP"
      attack_group_action = "alert"
    },
    {
      attack_group        = "RFI"
      attack_group_action = "alert"
    },
    {
      attack_group        = "PHP"
      attack_group_action = "alert"
    },
    {
      attack_group        = "TROJAN"
      attack_group_action = "alert"
    },
    {
      attack_group        = "DDOS"
      attack_group_action = "alert"
    },
    {
      attack_group        = "IN"
      attack_group_action = "deny_custom_78842"
    },
    {
      attack_group        = "OUT"
      attack_group_action = "alert"
    },
  ]
}

resource "akamai_appsec_waf_ruleset" "andrew" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.andrew.security_policy_id
  rules = [
    {
      rule_id     = "950002"
      rule_action = "alert"
    },
    {
      rule_id     = "950006"
      rule_action = "alert"
    },
    {
      rule_id     = "950007"
      rule_action = "alert"
    },
    {
      rule_id     = "950011"
      rule_action = "alert"
    },
    {
      rule_id     = "950118"
      rule_action = "alert"
    },
    {
      rule_id     = "950203"
      rule_action = "alert"
    },
    {
      rule_id     = "950216"
      rule_action = "alert"
    },
    {
      rule_id     = "950220"
      rule_action = "alert"
    },
    {
      rule_id     = "950901"
      rule_action = "alert"
    },
    {
      rule_id     = "951910"
      rule_action = "alert"
    },
    {
      rule_id     = "958003"
      rule_action = "alert"
    },
    {
      rule_id     = "958008"
      rule_action = "alert"
    },
    {
      rule_id     = "958023"
      rule_action = "alert"
    },
    {
      rule_id     = "958034"
      rule_action = "alert"
    },
    {
      rule_id     = "958051"
      rule_action = "alert"
    },
    {
      rule_id     = "958052"
      rule_action = "alert"
    },
    {
      rule_id     = "959070"
      rule_action = "alert"
    },
    {
      rule_id     = "959073"
      rule_action = "alert"
    },
    {
      rule_id     = "959976"
      rule_action = "alert"
    },
    {
      rule_id     = "959977"
      rule_action = "alert"
    },
    {
      rule_id     = "961011"
      rule_action = "alert"
    },
    {
      rule_id     = "961012"
      rule_action = "alert"
    },
    {
      rule_id     = "961034"
      rule_action = "alert"
    },
    {
      rule_id     = "961904"
      rule_action = "alert"
    },
    {
      rule_id     = "961912"
      rule_action = "alert"
    },
    {
      rule_id     = "968230"
      rule_action = "alert"
    },
    {
      rule_id     = "969151"
      rule_action = "alert"
    },
    {
      rule_id     = "970003"
      rule_action = "alert"
    },
    {
      rule_id     = "970004"
      rule_action = "alert"
    },
    {
      rule_id     = "970009"
      rule_action = "alert"
    },
    {
      rule_id     = "970011"
      rule_action = "alert"
    },
    {
      rule_id     = "970013"
      rule_action = "alert"
    },
    {
      rule_id     = "970014"
      rule_action = "alert"
    },
    {
      rule_id     = "970015"
      rule_action = "alert"
    },
    {
      rule_id     = "970118"
      rule_action = "alert"
    },
    {
      rule_id     = "970901"
      rule_action = "alert"
    },
    {
      rule_id     = "970902"
      rule_action = "alert"
    },
    {
      rule_id     = "970903"
      rule_action = "alert"
    },
    {
      rule_id     = "970904"
      rule_action = "alert"
    },
    {
      rule_id     = "973305"
      rule_action = "alert"
    },
    {
      rule_id     = "973307"
      rule_action = "alert"
    },
    {
      rule_id     = "973311"
      rule_action = "alert"
    },
    {
      rule_id     = "973312"
      rule_action = "alert"
    },
    {
      rule_id     = "973335"
      rule_action = "alert"
    },
    {
      rule_id     = "981240"
      rule_action = "alert"
    },
    {
      rule_id     = "981242"
      rule_action = "alert"
    },
    {
      rule_id     = "981243"
      rule_action = "alert"
    },
    {
      rule_id     = "981244"
      rule_action = "alert"
    },
    {
      rule_id     = "981247"
      rule_action = "alert"
    },
    {
      rule_id     = "981248"
      rule_action = "alert"
    },
    {
      rule_id     = "981251"
      rule_action = "alert"
    },
    {
      rule_id     = "981252"
      rule_action = "alert"
    },
    {
      rule_id     = "981253"
      rule_action = "alert"
    },
    {
      rule_id     = "981254"
      rule_action = "alert"
    },
    {
      rule_id     = "981255"
      rule_action = "alert"
    },
    {
      rule_id     = "981256"
      rule_action = "alert"
    },
    {
      rule_id     = "981260"
      rule_action = "alert"
    },
    {
      rule_id     = "981270"
      rule_action = "alert"
    },
    {
      rule_id     = "981276"
      rule_action = "alert"
    },
    {
      rule_id     = "981300"
      rule_action = "alert"
    },
    {
      rule_id     = "981320"
      rule_action = "alert"
    },
    {
      rule_id     = "999002"
      rule_action = "alert"
    },
    {
      rule_id     = "999901"
      rule_action = "alert"
    },
    {
      rule_id     = "999902"
      rule_action = "alert"
    },
    {
      rule_id     = "3000000"
      rule_action = "alert"
    },
    {
      rule_id     = "3000004"
      rule_action = "alert"
    },
    {
      rule_id     = "3000005"
      rule_action = "alert"
    },
    {
      rule_id     = "3000006"
      rule_action = "alert"
    },
    {
      rule_id     = "3000007"
      rule_action = "alert"
    },
    {
      rule_id     = "3000012"
      rule_action = "alert"
    },
    {
      rule_id     = "3000013"
      rule_action = "alert"
    },
    {
      rule_id     = "3000014"
      rule_action = "alert"
    },
    {
      rule_id     = "3000015"
      rule_action = "alert"
    },
    {
      rule_id     = "3000017"
      rule_action = "alert"
    },
    {
      rule_id     = "3000022"
      rule_action = "alert"
    },
    {
      rule_id     = "3000023"
      rule_action = "alert"
    },
    {
      rule_id     = "3000025"
      rule_action = "alert"
    },
    {
      rule_id     = "3000029"
      rule_action = "alert"
    },
    {
      rule_id     = "3000033"
      rule_action = "alert"
    },
    {
      rule_id     = "3000034"
      rule_action = "alert"
    },
    {
      rule_id     = "3000037"
      rule_action = "alert"
    },
    {
      rule_id     = "3000038"
      rule_action = "alert"
    },
    {
      rule_id     = "3000039"
      rule_action = "alert"
    },
    {
      rule_id     = "3000041"
      rule_action = "alert"
    },
    {
      rule_id     = "3000056"
      rule_action = "alert"
    },
    {
      rule_id     = "3000057"
      rule_action = "alert"
    },
    {
      rule_id     = "3000058"
      rule_action = "alert"
    },
    {
      rule_id     = "3000061"
      rule_action = "alert"
    },
    {
      rule_id     = "3000065"
      rule_action = "alert"
    },
    {
      rule_id     = "3000072"
      rule_action = "alert"
    },
    {
      rule_id     = "3000080"
      rule_action = "alert"
    },
    {
      rule_id     = "3000081"
      rule_action = "alert"
    },
    {
      rule_id     = "3000100"
      rule_action = "alert"
    },
    {
      rule_id     = "3000101"
      rule_action = "alert"
    },
    {
      rule_id     = "3000102"
      rule_action = "alert"
    },
    {
      rule_id     = "3000108"
      rule_action = "alert"
    },
    {
      rule_id     = "3000109"
      rule_action = "alert"
    },
    {
      rule_id     = "3000110"
      rule_action = "alert"
    },
    {
      rule_id     = "3000111"
      rule_action = "alert"
    },
    {
      rule_id     = "3000112"
      rule_action = "alert"
    },
    {
      rule_id     = "3000113"
      rule_action = "alert"
    },
    {
      rule_id     = "3000114"
      rule_action = "alert"
    },
    {
      rule_id     = "3000115"
      rule_action = "alert"
    },
    {
      rule_id     = "3000116"
      rule_action = "alert"
    },
    {
      rule_id     = "3000117"
      rule_action = "alert"
    },
    {
      rule_id     = "3000118"
      rule_action = "alert"
    },
    {
      rule_id     = "3000119"
      rule_action = "alert"
    },
    {
      rule_id     = "3000120"
      rule_action = "alert"
    },
    {
      rule_id     = "3000121"
      rule_action = "alert"
    },
    {
      rule_id     = "3000122"
      rule_action = "alert"
    },
    {
      rule_id     = "3000123"
      rule_action = "alert"
    },
    {
      rule_id     = "3000124"
      rule_action = "alert"
    },
    {
      rule_id     = "3000125"
      rule_action = "alert"
    },
    {
      rule_id     = "3000126"
      rule_action = "alert"
    },
    {
      rule_id     = "3000127"
      rule_action = "alert"
    },
    {
      rule_id     = "3000128"
      rule_action = "alert"
    },
    {
      rule_id     = "3000129"
      rule_action = "alert"
    },
    {
      rule_id     = "3000130"
      rule_action = "alert"
    },
    {
      rule_id     = "3000131"
      rule_action = "alert"
    },
    {
      rule_id     = "3000132"
      rule_action = "alert"
    },
    {
      rule_id     = "3000133"
      rule_action = "alert"
    },
    {
      rule_id     = "3000134"
      rule_action = "alert"
    },
    {
      rule_id     = "3000135"
      rule_action = "alert"
    },
    {
      rule_id     = "3000136"
      rule_action = "alert"
    },
    {
      rule_id     = "3000137"
      rule_action = "alert"
    },
    {
      rule_id     = "3000138"
      rule_action = "alert"
    },
    {
      rule_id     = "3000139"
      rule_action = "alert"
    },
    {
      rule_id     = "3000140"
      rule_action = "alert"
    },
    {
      rule_id     = "3000141"
      rule_action = "alert"
    },
    {
      rule_id     = "3000142"
      rule_action = "alert"
    },
    {
      rule_id     = "3000143"
      rule_action = "alert"
    },
    {
      rule_id     = "3000144"
      rule_action = "alert"
    },
    {
      rule_id     = "3000145"
      rule_action = "alert"
    },
    {
      rule_id     = "3000146"
      rule_action = "alert"
    },
    {
      rule_id     = "3000147"
      rule_action = "alert"
    },
    {
      rule_id     = "3000148"
      rule_action = "alert"
    },
    {
      rule_id     = "3000149"
      rule_action = "alert"
    },
    {
      rule_id     = "3000150"
      rule_action = "alert"
    },
    {
      rule_id     = "3000151"
      rule_action = "alert"
    },
    {
      rule_id     = "3000152"
      rule_action = "alert"
    },
    {
      rule_id     = "3000153"
      rule_action = "alert"
    },
    {
      rule_id     = "3000154"
      rule_action = "alert"
    },
    {
      rule_id     = "3000155"
      rule_action = "alert"
    },
    {
      rule_id     = "3000156"
      rule_action = "alert"
    },
    {
      rule_id     = "3000157"
      rule_action = "alert"
    },
    {
      rule_id     = "3000160"
      rule_action = "alert"
    },
    {
      rule_id     = "3000162"
      rule_action = "alert"
    },
    {
      rule_id     = "3000164"
      rule_action = "alert"
    },
    {
      rule_id     = "3000166"
      rule_action = "alert"
    },
    {
      rule_id     = "3000167"
      rule_action = "alert"
    },
    {
      rule_id     = "3000168"
      rule_action = "alert"
    },
    {
      rule_id     = "3000171"
      rule_action = "alert"
    },
    {
      rule_id     = "3000173"
      rule_action = "alert"
    },
    {
      rule_id     = "3000174"
      rule_action = "alert"
    },
    {
      rule_id     = "3000175"
      rule_action = "alert"
    },
    {
      rule_id     = "3000176"
      rule_action = "alert"
    },
    {
      rule_id     = "3000177"
      rule_action = "alert"
    },
    {
      rule_id     = "3000179"
      rule_action = "alert"
    },
    {
      rule_id     = "3000180"
      rule_action = "alert"
    },
    {
      rule_id     = "3000183"
      rule_action = "alert"
    },
    {
      rule_id     = "3000184"
      rule_action = "alert"
    },
    {
      rule_id     = "3000185"
      rule_action = "alert"
    },
    {
      rule_id     = "3000186"
      rule_action = "alert"
    },
    {
      rule_id     = "3000187"
      rule_action = "alert"
    },
    {
      rule_id     = "3000188"
      rule_action = "alert"
    },
  ]
  attack_groups = [
    {
      attack_group        = "POLICY"
      attack_group_action = "alert"
    },
    {
      attack_group        = "WAT"
      attack_group_action = "alert"
    },
    {
      attack_group        = "PROTOCOL"
      attack_group_action = "alert"
    },
    {
      attack_group        = "SQL"
      attack_group_action = "alert"
    },
    {
      attack_group        = "XSS"
      attack_group_action = "alert"
    },
    {
      attack_group        = "CMD"
      attack_group_action = "alert"
    },
    {
      attack_group        = "LFI"
      attack_group_action = "alert"
    },
    {
      attack_group        = "RFI"
      attack_group_action = "alert"
    },
    {
      attack_group        = "PLATFORM"
      attack_group_action = "alert"
    },
  ]
}

resource "akamai_appsec_waf_ruleset" "policy1" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.policy1.security_policy_id
  rules = [
    {
      rule_id     = "699989"
      rule_action = "deny"
    },
    {
      rule_id     = "699996"
      rule_action = "alert"
    },
    {
      rule_id     = "950000"
      rule_action = "alert"
    },
    {
      rule_id     = "950001"
      rule_action = "alert"
    },
    {
      rule_id     = "950002"
      rule_action = "alert"
      condition_exception = jsonencode(merge(
        { "advancedExceptions" : {
          "conditionOperator" : "OR",
          "conditions" : [
            {
              "type" : "hostMatch",
              "hosts" : [
                "test.com"
              ],
              "positiveMatch" : true
            },
            {
              "type" : "pathMatch",
              "paths" : [
                "/test"
              ],
              "positiveMatch" : true
            },
            {
              "type" : "requestHeaderMatch",
              "header" : "Test",
              "positiveMatch" : true,
              "value" : "test*"
            }
          ]
        } }
      ))
    },
    {
      rule_id     = "950003"
      rule_action = "alert"
    },
    {
      rule_id     = "950005"
      rule_action = "alert"
    },
    {
      rule_id     = "950006"
      rule_action = "alert"
    },
    {
      rule_id     = "950007"
      rule_action = "alert"
    },
    {
      rule_id     = "950008"
      rule_action = "alert"
    },
    {
      rule_id     = "950009"
      rule_action = "alert"
    },
    {
      rule_id     = "950010"
      rule_action = "alert"
    },
    {
      rule_id     = "950011"
      rule_action = "alert"
    },
    {
      rule_id     = "950018"
      rule_action = "alert"
    },
    {
      rule_id     = "950019"
      rule_action = "alert"
    },
    {
      rule_id     = "950103"
      rule_action = "alert"
    },
    {
      rule_id     = "950107"
      rule_action = "alert"
    },
    {
      rule_id     = "950108"
      rule_action = "alert"
    },
    {
      rule_id     = "950109"
      rule_action = "alert"
    },
    {
      rule_id     = "950110"
      rule_action = "alert"
    },
    {
      rule_id     = "950116"
      rule_action = "alert"
    },
    {
      rule_id     = "950117"
      rule_action = "alert"
    },
    {
      rule_id     = "950118"
      rule_action = "alert"
    },
    {
      rule_id     = "950119"
      rule_action = "alert"
    },
    {
      rule_id     = "950120"
      rule_action = "alert"
    },
    {
      rule_id     = "950901"
      rule_action = "alert"
    },
    {
      rule_id     = "950908"
      rule_action = "alert"
    },
    {
      rule_id     = "950910"
      rule_action = "alert"
    },
    {
      rule_id     = "950911"
      rule_action = "alert"
    },
    {
      rule_id     = "950921"
      rule_action = "alert"
    },
    {
      rule_id     = "958000"
      rule_action = "alert"
    },
    {
      rule_id     = "958001"
      rule_action = "alert"
    },
    {
      rule_id     = "958002"
      rule_action = "alert"
    },
    {
      rule_id     = "958003"
      rule_action = "alert"
    },
    {
      rule_id     = "958004"
      rule_action = "alert"
    },
    {
      rule_id     = "958005"
      rule_action = "alert"
    },
    {
      rule_id     = "958006"
      rule_action = "alert"
    },
    {
      rule_id     = "958007"
      rule_action = "alert"
    },
    {
      rule_id     = "958008"
      rule_action = "alert"
    },
    {
      rule_id     = "958009"
      rule_action = "alert"
    },
    {
      rule_id     = "958010"
      rule_action = "alert"
    },
    {
      rule_id     = "958011"
      rule_action = "alert"
    },
    {
      rule_id     = "958012"
      rule_action = "alert"
    },
    {
      rule_id     = "958013"
      rule_action = "alert"
    },
    {
      rule_id     = "958016"
      rule_action = "alert"
    },
    {
      rule_id     = "958017"
      rule_action = "alert"
    },
    {
      rule_id     = "958018"
      rule_action = "alert"
    },
    {
      rule_id     = "958019"
      rule_action = "alert"
    },
    {
      rule_id     = "958020"
      rule_action = "alert"
    },
    {
      rule_id     = "958022"
      rule_action = "alert"
    },
    {
      rule_id     = "958023"
      rule_action = "alert"
    },
    {
      rule_id     = "958024"
      rule_action = "alert"
    },
    {
      rule_id     = "958025"
      rule_action = "alert"
    },
    {
      rule_id     = "958026"
      rule_action = "alert"
    },
    {
      rule_id     = "958027"
      rule_action = "alert"
    },
    {
      rule_id     = "958028"
      rule_action = "alert"
    },
    {
      rule_id     = "958030"
      rule_action = "alert"
    },
    {
      rule_id     = "958031"
      rule_action = "alert"
    },
    {
      rule_id     = "958032"
      rule_action = "alert"
    },
    {
      rule_id     = "958033"
      rule_action = "alert"
    },
    {
      rule_id     = "958034"
      rule_action = "alert"
    },
    {
      rule_id     = "958036"
      rule_action = "alert"
    },
    {
      rule_id     = "958037"
      rule_action = "alert"
    },
    {
      rule_id     = "958038"
      rule_action = "alert"
    },
    {
      rule_id     = "958039"
      rule_action = "alert"
    },
    {
      rule_id     = "958040"
      rule_action = "alert"
    },
    {
      rule_id     = "958041"
      rule_action = "alert"
    },
    {
      rule_id     = "958045"
      rule_action = "alert"
    },
    {
      rule_id     = "958046"
      rule_action = "alert"
    },
    {
      rule_id     = "958047"
      rule_action = "alert"
    },
    {
      rule_id     = "958049"
      rule_action = "alert"
    },
    {
      rule_id     = "958051"
      rule_action = "alert"
    },
    {
      rule_id     = "958052"
      rule_action = "alert"
    },
    {
      rule_id     = "958054"
      rule_action = "alert"
    },
    {
      rule_id     = "958056"
      rule_action = "alert"
    },
    {
      rule_id     = "958057"
      rule_action = "alert"
    },
    {
      rule_id     = "958059"
      rule_action = "alert"
    },
    {
      rule_id     = "958230"
      rule_action = "alert"
    },
    {
      rule_id     = "958231"
      rule_action = "alert"
    },
    {
      rule_id     = "958291"
      rule_action = "alert"
    },
    {
      rule_id     = "958295"
      rule_action = "alert"
    },
    {
      rule_id     = "958404"
      rule_action = "alert"
    },
    {
      rule_id     = "958405"
      rule_action = "alert"
    },
    {
      rule_id     = "958406"
      rule_action = "alert"
    },
    {
      rule_id     = "958407"
      rule_action = "alert"
    },
    {
      rule_id     = "958408"
      rule_action = "alert"
    },
    {
      rule_id     = "958409"
      rule_action = "alert"
    },
    {
      rule_id     = "958410"
      rule_action = "alert"
    },
    {
      rule_id     = "958411"
      rule_action = "alert"
    },
    {
      rule_id     = "958412"
      rule_action = "alert"
    },
    {
      rule_id     = "958413"
      rule_action = "alert"
    },
    {
      rule_id     = "958414"
      rule_action = "alert"
    },
    {
      rule_id     = "958415"
      rule_action = "alert"
    },
    {
      rule_id     = "958416"
      rule_action = "alert"
    },
    {
      rule_id     = "958417"
      rule_action = "alert"
    },
    {
      rule_id     = "958418"
      rule_action = "alert"
    },
    {
      rule_id     = "958419"
      rule_action = "alert"
    },
    {
      rule_id     = "958420"
      rule_action = "alert"
    },
    {
      rule_id     = "958421"
      rule_action = "alert"
    },
    {
      rule_id     = "958422"
      rule_action = "alert"
    },
    {
      rule_id     = "958423"
      rule_action = "alert"
    },
    {
      rule_id     = "958976"
      rule_action = "alert"
    },
    {
      rule_id     = "958977"
      rule_action = "alert"
    },
    {
      rule_id     = "959070"
      rule_action = "alert"
    },
    {
      rule_id     = "959071"
      rule_action = "alert"
    },
    {
      rule_id     = "959072"
      rule_action = "alert"
    },
    {
      rule_id     = "959073"
      rule_action = "alert"
    },
    {
      rule_id     = "959151"
      rule_action = "alert"
    },
    {
      rule_id     = "960010"
      rule_action = "alert"
    },
    {
      rule_id     = "960011"
      rule_action = "alert"
    },
    {
      rule_id     = "960012"
      rule_action = "alert"
    },
    {
      rule_id     = "960016"
      rule_action = "alert"
    },
    {
      rule_id     = "960022"
      rule_action = "alert"
    },
    {
      rule_id     = "960034"
      rule_action = "alert"
    },
    {
      rule_id     = "960035"
      rule_action = "alert"
    },
    {
      rule_id     = "960208"
      rule_action = "alert"
    },
    {
      rule_id     = "960209"
      rule_action = "alert"
    },
    {
      rule_id     = "960335"
      rule_action = "alert"
    },
    {
      rule_id     = "960341"
      rule_action = "alert"
    },
    {
      rule_id     = "960901"
      rule_action = "alert"
    },
    {
      rule_id     = "960902"
      rule_action = "alert"
    },
    {
      rule_id     = "960904"
      rule_action = "alert"
    },
    {
      rule_id     = "960912"
      rule_action = "alert"
    },
    {
      rule_id     = "960913"
      rule_action = "alert"
    },
    {
      rule_id     = "960914"
      rule_action = "alert"
    },
    {
      rule_id     = "973300"
      rule_action = "alert"
    },
    {
      rule_id     = "973301"
      rule_action = "alert"
    },
    {
      rule_id     = "973302"
      rule_action = "alert"
    },
    {
      rule_id     = "973303"
      rule_action = "alert"
    },
    {
      rule_id     = "973304"
      rule_action = "alert"
    },
    {
      rule_id     = "973305"
      rule_action = "alert"
    },
    {
      rule_id     = "973306"
      rule_action = "alert"
    },
    {
      rule_id     = "973307"
      rule_action = "alert"
    },
    {
      rule_id     = "973308"
      rule_action = "alert"
    },
    {
      rule_id     = "973309"
      rule_action = "alert"
    },
    {
      rule_id     = "973310"
      rule_action = "alert"
    },
    {
      rule_id     = "973311"
      rule_action = "alert"
    },
    {
      rule_id     = "973312"
      rule_action = "alert"
    },
    {
      rule_id     = "973313"
      rule_action = "alert"
    },
    {
      rule_id     = "973314"
      rule_action = "alert"
    },
    {
      rule_id     = "973315"
      rule_action = "alert"
    },
    {
      rule_id     = "973316"
      rule_action = "alert"
    },
    {
      rule_id     = "973317"
      rule_action = "alert"
    },
    {
      rule_id     = "973318"
      rule_action = "alert"
    },
    {
      rule_id     = "973319"
      rule_action = "alert"
    },
    {
      rule_id     = "973320"
      rule_action = "alert"
    },
    {
      rule_id     = "973321"
      rule_action = "alert"
    },
    {
      rule_id     = "973322"
      rule_action = "alert"
    },
    {
      rule_id     = "973323"
      rule_action = "alert"
    },
    {
      rule_id     = "973324"
      rule_action = "alert"
    },
    {
      rule_id     = "973325"
      rule_action = "alert"
    },
    {
      rule_id     = "973326"
      rule_action = "alert"
    },
    {
      rule_id     = "973327"
      rule_action = "alert"
    },
    {
      rule_id     = "973328"
      rule_action = "alert"
    },
    {
      rule_id     = "973329"
      rule_action = "alert"
    },
    {
      rule_id     = "973330"
      rule_action = "alert"
    },
    {
      rule_id     = "973331"
      rule_action = "alert"
    },
    {
      rule_id     = "973332"
      rule_action = "alert"
    },
    {
      rule_id     = "973333"
      rule_action = "alert"
    },
    {
      rule_id     = "973334"
      rule_action = "alert"
    },
    {
      rule_id     = "973335"
      rule_action = "alert"
    },
    {
      rule_id     = "973336"
      rule_action = "alert"
    },
    {
      rule_id     = "973337"
      rule_action = "alert"
    },
    {
      rule_id     = "981173"
      rule_action = "alert"
    },
    {
      rule_id     = "981241"
      rule_action = "alert"
    },
    {
      rule_id     = "981242"
      rule_action = "alert"
    },
    {
      rule_id     = "981243"
      rule_action = "alert"
    },
    {
      rule_id     = "981244"
      rule_action = "alert"
    },
    {
      rule_id     = "981245"
      rule_action = "alert"
    },
    {
      rule_id     = "981246"
      rule_action = "alert"
    },
    {
      rule_id     = "981247"
      rule_action = "alert"
    },
    {
      rule_id     = "981248"
      rule_action = "alert"
    },
    {
      rule_id     = "981249"
      rule_action = "alert"
    },
    {
      rule_id     = "981250"
      rule_action = "alert"
    },
    {
      rule_id     = "981251"
      rule_action = "alert"
    },
    {
      rule_id     = "981252"
      rule_action = "alert"
    },
    {
      rule_id     = "981253"
      rule_action = "alert"
    },
    {
      rule_id     = "981254"
      rule_action = "alert"
    },
    {
      rule_id     = "981255"
      rule_action = "alert"
    },
    {
      rule_id     = "981256"
      rule_action = "alert"
    },
    {
      rule_id     = "981260"
      rule_action = "alert"
    },
    {
      rule_id     = "981270"
      rule_action = "alert"
    },
    {
      rule_id     = "981272"
      rule_action = "alert"
    },
    {
      rule_id     = "981276"
      rule_action = "alert"
    },
    {
      rule_id     = "981277"
      rule_action = "alert"
    },
    {
      rule_id     = "981300"
      rule_action = "alert"
    },
    {
      rule_id     = "981318"
      rule_action = "alert"
    },
    {
      rule_id     = "981319"
      rule_action = "alert"
    },
    {
      rule_id     = "981320"
      rule_action = "alert"
    },
    {
      rule_id     = "990002"
      rule_action = "alert"
    },
    {
      rule_id     = "990012"
      rule_action = "alert"
    },
    {
      rule_id     = "990901"
      rule_action = "alert"
    },
    {
      rule_id     = "990902"
      rule_action = "alert"
    },
    {
      rule_id     = "3000000"
      rule_action = "alert"
    },
    {
      rule_id     = "3000001"
      rule_action = "alert"
    },
    {
      rule_id     = "3000002"
      rule_action = "alert"
    },
    {
      rule_id     = "3000003"
      rule_action = "alert"
    },
    {
      rule_id     = "3000004"
      rule_action = "alert"
    },
    {
      rule_id     = "3000005"
      rule_action = "alert"
    },
    {
      rule_id     = "3000006"
      rule_action = "alert"
    },
    {
      rule_id     = "3000007"
      rule_action = "alert"
    },
    {
      rule_id     = "3000008"
      rule_action = "alert"
    },
    {
      rule_id     = "3000009"
      rule_action = "alert"
    },
    {
      rule_id     = "3000010"
      rule_action = "alert"
    },
    {
      rule_id     = "3000011"
      rule_action = "alert"
    },
    {
      rule_id     = "3000012"
      rule_action = "alert"
    },
    {
      rule_id     = "3000013"
      rule_action = "alert"
    },
    {
      rule_id     = "3000014"
      rule_action = "alert"
    },
    {
      rule_id     = "3000015"
      rule_action = "alert"
    },
    {
      rule_id     = "3000016"
      rule_action = "alert"
    },
    {
      rule_id     = "3000017"
      rule_action = "alert"
    },
    {
      rule_id     = "3000018"
      rule_action = "alert"
    },
    {
      rule_id     = "3000019"
      rule_action = "alert"
    },
    {
      rule_id     = "3000020"
      rule_action = "alert"
    },
    {
      rule_id     = "3000021"
      rule_action = "alert"
    },
    {
      rule_id     = "3000022"
      rule_action = "alert"
    },
    {
      rule_id     = "3000023"
      rule_action = "alert"
    },
    {
      rule_id     = "3000024"
      rule_action = "alert"
    },
    {
      rule_id     = "3000025"
      rule_action = "alert"
    },
    {
      rule_id     = "3000027"
      rule_action = "alert"
    },
    {
      rule_id     = "3000029"
      rule_action = "alert"
    },
    {
      rule_id     = "3000030"
      rule_action = "alert"
    },
    {
      rule_id     = "3000031"
      rule_action = "alert"
    },
    {
      rule_id     = "3000032"
      rule_action = "alert"
    },
    {
      rule_id     = "3000033"
      rule_action = "alert"
    },
    {
      rule_id     = "3000034"
      rule_action = "alert"
    },
    {
      rule_id     = "3000035"
      rule_action = "alert"
    },
    {
      rule_id     = "3000036"
      rule_action = "alert"
    },
    {
      rule_id     = "3000037"
      rule_action = "alert"
    },
    {
      rule_id     = "3000038"
      rule_action = "alert"
    },
    {
      rule_id     = "3000039"
      rule_action = "alert"
    },
    {
      rule_id     = "3000040"
      rule_action = "alert"
    },
    {
      rule_id     = "3000041"
      rule_action = "alert"
    },
    {
      rule_id     = "3000042"
      rule_action = "alert"
    },
    {
      rule_id     = "3000043"
      rule_action = "alert"
    },
    {
      rule_id     = "3000044"
      rule_action = "alert"
    },
    {
      rule_id     = "3000045"
      rule_action = "alert"
    },
    {
      rule_id     = "3000046"
      rule_action = "alert"
    },
    {
      rule_id     = "3000047"
      rule_action = "alert"
    },
    {
      rule_id     = "3000048"
      rule_action = "alert"
    },
    {
      rule_id     = "3000049"
      rule_action = "alert"
    },
    {
      rule_id     = "3000050"
      rule_action = "alert"
    },
    {
      rule_id     = "3000051"
      rule_action = "alert"
    },
    {
      rule_id     = "3000052"
      rule_action = "alert"
    },
    {
      rule_id     = "3000053"
      rule_action = "alert"
    },
    {
      rule_id     = "3000054"
      rule_action = "alert"
    },
    {
      rule_id     = "3000055"
      rule_action = "alert"
    },
    {
      rule_id     = "3000056"
      rule_action = "alert"
    },
    {
      rule_id     = "3000057"
      rule_action = "alert"
    },
    {
      rule_id     = "3000058"
      rule_action = "alert"
    },
    {
      rule_id     = "3000059"
      rule_action = "alert"
    },
    {
      rule_id     = "3000060"
      rule_action = "alert"
    },
    {
      rule_id     = "3000061"
      rule_action = "alert"
    },
    {
      rule_id     = "3000062"
      rule_action = "alert"
    },
    {
      rule_id     = "3000063"
      rule_action = "alert"
    },
    {
      rule_id     = "3000064"
      rule_action = "alert"
    },
    {
      rule_id     = "3000065"
      rule_action = "alert"
    },
    {
      rule_id     = "3000066"
      rule_action = "alert"
    },
    {
      rule_id     = "3000067"
      rule_action = "alert"
    },
    {
      rule_id     = "3000068"
      rule_action = "alert"
    },
    {
      rule_id     = "3000071"
      rule_action = "alert"
    },
    {
      rule_id     = "3000072"
      rule_action = "alert"
    },
    {
      rule_id     = "3000073"
      rule_action = "alert"
    },
    {
      rule_id     = "3000074"
      rule_action = "alert"
    },
    {
      rule_id     = "3000075"
      rule_action = "alert"
    },
    {
      rule_id     = "3000076"
      rule_action = "alert"
    },
    {
      rule_id     = "3000077"
      rule_action = "alert"
    },
    {
      rule_id     = "3000079"
      rule_action = "alert"
    },
    {
      rule_id     = "3000080"
      rule_action = "alert"
    },
    {
      rule_id     = "3000081"
      rule_action = "alert"
    },
    {
      rule_id     = "3000082"
      rule_action = "alert"
    },
    {
      rule_id     = "3000083"
      rule_action = "alert"
    },
    {
      rule_id     = "3000084"
      rule_action = "alert"
    },
    {
      rule_id     = "3000999"
      rule_action = "alert"
    },
  ]
  attack_groups = [
    {
      attack_group        = "SQL"
      attack_group_action = "alert"
    },
    {
      attack_group        = "XSS"
      attack_group_action = "alert"
    },
    {
      attack_group        = "CMD"
      attack_group_action = "deny"
    },
    {
      attack_group        = "HTTP"
      attack_group_action = "alert"
    },
    {
      attack_group        = "RFI"
      attack_group_action = "alert"
    },
    {
      attack_group        = "PHP"
      attack_group_action = "deny_custom_78842"
    },
    {
      attack_group        = "TROJAN"
      attack_group_action = "alert"
    },
    {
      attack_group        = "IN"
      attack_group_action = "deny"
    },
    {
      attack_group        = "OUT"
      attack_group_action = "alert"
    },
  ]
}

