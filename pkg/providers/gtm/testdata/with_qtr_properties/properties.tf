resource "akamai_gtm_property" "test_property1" {
  domain                      = akamai_gtm_domain.test_name.name
  name                        = "test property1"
  type                        = "qtr"
  ipv6                        = false
  score_aggregation_type      = "worst"
  stickiness_bonus_percentage = 0
  stickiness_bonus_constant   = 0
  use_computed_targets        = false
  balance_by_download_score   = false
  dynamic_ttl                 = 60
  handout_limit               = 8
  handout_mode                = "normal"
  failover_delay              = 0
  failback_delay              = 0
  ghost_demand_reporting      = false
  traffic_target {
    datacenter_id = data.akamai_gtm_default_datacenter.default_datacenter_5401.datacenter_id
    enabled       = true
    weight        = 1
    servers       = ["1.2.3.4"]
  }
  liveness_test {
    name                             = "HTTP"
    peer_certificate_verification    = false
    test_interval                    = 60
    test_object                      = "/"
    http_error3xx                    = true
    http_error4xx                    = true
    http_error5xx                    = true
    disabled                         = false
    test_object_protocol             = "HTTP"
    test_object_port                 = 80
    disable_nonstandard_port_warning = false
    test_timeout                     = 10
    answers_required                 = false
    recursion_requested              = false
  }
  depends_on = [
    data.akamai_gtm_default_datacenter.default_datacenter_5401,
    akamai_gtm_domain.test_name
  ]
}

resource "akamai_gtm_property" "test_property2" {
  domain                      = akamai_gtm_domain.test_name.name
  name                        = "test property2"
  type                        = "qtr"
  ipv6                        = false
  score_aggregation_type      = "worst"
  stickiness_bonus_percentage = 0
  stickiness_bonus_constant   = 0
  use_computed_targets        = false
  balance_by_download_score   = false
  dynamic_ttl                 = 60
  handout_limit               = 8
  handout_mode                = "normal"
  failover_delay              = 0
  failback_delay              = 0
  ghost_demand_reporting      = false
  traffic_target {
    datacenter_id = akamai_gtm_datacenter.TEST1.datacenter_id
    enabled       = true
    weight        = 1
    servers       = ["1.2.3.4"]
  }
  traffic_target {
    datacenter_id = data.akamai_gtm_default_datacenter.default_datacenter_5402.datacenter_id
    enabled       = true
    weight        = 1
    servers       = ["7.6.5.4"]
  }
  liveness_test {
    name                             = "HTTP"
    peer_certificate_verification    = false
    test_interval                    = 60
    test_object                      = "/"
    http_error3xx                    = true
    http_error4xx                    = true
    http_error5xx                    = true
    disabled                         = false
    test_object_protocol             = "HTTP"
    test_object_port                 = 80
    disable_nonstandard_port_warning = false
    test_timeout                     = 10
    answers_required                 = false
    recursion_requested              = false
  }
  depends_on = [
    akamai_gtm_datacenter.TEST1,
    data.akamai_gtm_default_datacenter.default_datacenter_5402,
    akamai_gtm_domain.test_name
  ]
}

