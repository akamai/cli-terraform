
data "akamai_property_rules_builder" "test-edgesuite-net_rule_default" {
  rules_v2024_01_09 {
    name      = "default"
    is_secure = false
    uuid      = "default"
    behavior {
      origin {
        cache_key_hostname = "ORIGIN_HOSTNAME"
        compress           = true
        custom_certificate_authorities {}
        custom_certificates {}
        enable_true_client_ip            = false
        forward_host_header              = "REQUEST_HOST_HEADER"
        hostname                         = "1.2.3.4"
        http_port                        = 80
        https_port                       = 443
        origin_certs_to_honor            = "COMBO"
        origin_sni                       = false
        origin_type                      = "CUSTOMER"
        standard_certificate_authorities = []
        use_unique_cache_key             = false
        verification_mode                = "CUSTOM"
      }
    }
  }
}
