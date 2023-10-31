resource "akamai_gtm_property" "test_property1" {
  domain                      = akamai_gtm_domain.test_name.name
  name                        = "test property1"
  type                        = "static"
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
  comments                    = <<EOT
first
second
EOT
  ghost_demand_reporting      = false
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
    akamai_gtm_domain.test_name
  ]
}

