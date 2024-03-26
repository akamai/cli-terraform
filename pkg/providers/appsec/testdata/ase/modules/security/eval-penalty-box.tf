// Eval Penalty Box
resource "akamai_appsec_eval_penalty_box" "default_policy" {
  config_id              = akamai_appsec_configuration.config.config_id
  security_policy_id     = akamai_appsec_security_policy.default_policy.security_policy_id
  penalty_box_protection = true
  penalty_box_action     = "alert"
}

// Eval Penalty Box Conditions
resource "akamai_appsec_eval_penalty_box_conditions" "eval_penalty_box_conditions" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.default_policy.security_policy_id
  penalty_box_conditions = jsonencode(
    {
      "conditionOperator" : "AND",
      "conditions" : [
        {
          "type" : "filenameMatch",
          "filenames" : [
            "appptest45"
          ],
          "positiveMatch" : true
        }
      ]
    }
  )
}
