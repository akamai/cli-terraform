resource "akamai_appsec_security_policy" "policy2" {
  config_id              = local.config_id
  default_settings       = true
  security_policy_name   = "policy2"
  security_policy_prefix = "hard"
}

resource "akamai_appsec_security_policy" "andrew" {
  config_id              = local.config_id
  default_settings       = true
  security_policy_name   = "andrew"
  security_policy_prefix = "last"
}

resource "akamai_appsec_security_policy" "policy1" {
  config_id              = local.config_id
  default_settings       = true
  security_policy_name   = "policy1"
  security_policy_prefix = "easy"
}

