
data "akamai_property_rules_builder" "test-edgesuite-net_rule_default" {
  rules_v2024_08_13 {
    name      = "default"
    is_secure = false
    uuid      = "default"
    behavior {
      datastream {
        beacon_stream_title      = "test"
        collect_midgress_traffic = false
        datastream_ids           = "77-85-6"
        enabled                  = true
        log_enabled              = false
        log_stream_name          = ["60", ]
        log_stream_title         = "test"
        sampling_percentage      = 100
        stream_type              = "LOG"
      }
    }
  }
}
