// API Request Constraints
resource "akamai_appsec_api_request_constraints" "andrew_767805" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_api_constraints_protection.andrew.security_policy_id
  api_endpoint_id    = 767805 // Note: We don't have an API Endpoint Definitions in our provider yet so can't reference this ID to another resource
  action             = "alert"
}
