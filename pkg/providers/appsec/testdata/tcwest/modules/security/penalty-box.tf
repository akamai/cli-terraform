// Penalty Box
resource "akamai_appsec_penalty_box" "andrew" {
  config_id              = local.config_id
  security_policy_id     = akamai_appsec_security_policy.andrew.security_policy_id
  penalty_box_protection = true
  penalty_box_action     = "alert"
}
// Penalty Box Conditions
resource "akamai_appsec_penalty_box_conditions" "andrew" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.andrew.security_policy_id
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

