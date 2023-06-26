// IP/GEO Firewall
resource "akamai_appsec_ip_geo" "policy2" {
  config_id                  = akamai_appsec_configuration.config.config_id
  security_policy_id         = akamai_appsec_ip_geo_protection.policy2.security_policy_id
  mode                       = "block"
  geo_network_lists          = ["32346_JTONGGEOTEST"]
  ip_network_lists           = ["16360_AISHTEST"]
  exception_ip_network_lists = ["16656_CPISERVERS"]
}

// IP/GEO Firewall
resource "akamai_appsec_ip_geo" "andrew" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_ip_geo_protection.andrew.security_policy_id
  mode               = "block"
}

// IP/GEO Firewall
resource "akamai_appsec_ip_geo" "policy1" {
  config_id                  = akamai_appsec_configuration.config.config_id
  security_policy_id         = akamai_appsec_ip_geo_protection.policy1.security_policy_id
  mode                       = "block"
  geo_network_lists          = ["113698_CUSTOMERGEOBLOCK"]
  ip_network_lists           = ["19843_TESTLIST"]
  exception_ip_network_lists = ["9132_MYTESTLIST"]
  ukraine_geo_control_action = "deny"
}

