// IP/GEO Firewall
resource "akamai_appsec_ip_geo" "default_policy" {
  config_id                  = akamai_appsec_configuration.config.config_id
  security_policy_id         = akamai_appsec_ip_geo_protection.default_policy.security_policy_id
  mode                       = "block"
  ip_network_lists           = ["118736_TFDEMOLISTATUL"]
  ukraine_geo_control_action = "deny"
}

