// API Request Constraints
resource "akamai_appsec_api_request_constraints" "default_policy_12345" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_api_constraints_protection.default_policy.security_policy_id
  api_endpoint_id    = 12345 // Note: We don't have an API Endpoint Definitions in our provider yet so can't reference this ID to another resource
  action             = "alert"
}

resource "akamai_appsec_api_request_constraints" "default_policy_12346" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_api_constraints_protection.default_policy.security_policy_id
  api_endpoint_id    = 12346 // Note: We don't have an API Endpoint Definitions in our provider yet so can't reference this ID to another resource
  action             = "alert"
}
