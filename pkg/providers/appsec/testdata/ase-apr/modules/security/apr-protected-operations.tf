resource "akamai_apr_protected_operations" "default_policy_c2b20de9-bb2e-4da8-ac99-f5b6e9ae4e10" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.default_policy.security_policy_id
  operation_id       = "c2b20de9-bb2e-4da8-ac99-f5b6e9ae4e10"
  protected_operation = jsonencode(
    {
      "apiEndPointId" : 650201,
      "telemetryTypeStates" : {
        "inline" : {
          "ajaxSupportEnabled" : false,
          "disabledAction" : "none",
          "enabled" : false
        },
        "nativeSdk" : {
          "ajaxSupportEnabled" : false,
          "disabledAction" : "none",
          "enabled" : false
        },
        "standard" : {
          "ajaxSupportEnabled" : false,
          "disabledAction" : "monitor",
          "enabled" : true
        }
      },
      "traffic" : {
        "inline" : {
          "aggressive" : {
            "action" : "monitor"
          },
          "cautious" : {
            "action" : "monitor"
          },
          "overrideThresholds" : false,
          "strict" : {
            "action" : "monitor"
          }
        },
        "standard" : {
          "aggressive" : {
            "action" : "monitor",
            "threshold" : 76
          },
          "cautious" : {
            "action" : "monitor",
            "threshold" : 22
          },
          "overrideThresholds" : true,
          "strict" : {
            "action" : "deny",
            "threshold" : 51
          }
        }
      }
    }
  )
}
