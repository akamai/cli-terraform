resource "akamai_apr_general_settings" "default_policy" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.default_policy.security_policy_id
  general_settings = jsonencode(
    {
      "accountProtection" : true,
      "originSignalHeader" : false,
      "originUserIdInRequestHeader" : false,
      "usernameInRequestHeader" : false
    }
  )
}
