// IP/GEO/ASN Firewall
resource "akamai_appsec_ip_geo" "default_policy" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_ip_geo_protection.default_policy.security_policy_id
  mode               = "block"
  asn_controls {
    action            = "deny"
    asn_network_lists = ["119711_ASNLIST"]
  }
  ip_controls {
    action           = "deny"
    ip_network_lists = ["118736_TFDEMOLISTATUL"]
  }
  ukraine_geo_control_action = "deny"
}

