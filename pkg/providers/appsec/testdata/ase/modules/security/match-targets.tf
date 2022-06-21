resource "akamai_appsec_match_target" "website_4262513" {
  config_id = akamai_appsec_configuration.config.config_id
  match_target = jsonencode(
    {
      "defaultFile" : "NO_MATCH",
      "filePaths" : [
        "/*"
      ],
      "hostnames" : [
        "www.rlw7w.uk"
      ],
      "isNegativeFileExtensionMatch" : false,
      "isNegativePathMatch" : false,
      "securityPolicy" : {
        "policyId" : "${akamai_appsec_security_policy.default_policy.security_policy_id}"
      },
      "sequence" : 0,
      "type" : "website"
    }
  )
}
