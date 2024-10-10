resource "akamai_appsec_wap_selected_hostnames" "default_policy" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.default_policy.security_policy_id
  protected_hosts    = ["www.rlw7w.uk"]
  evaluated_hosts    = []
}

