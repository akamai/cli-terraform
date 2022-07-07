resource "akamai_appsec_security_policy" "default_policy" {
  config_id              = akamai_appsec_configuration.config.config_id
  default_settings       = true
  security_policy_name   = "Default Policy"
  security_policy_prefix = "ASE1"
}

