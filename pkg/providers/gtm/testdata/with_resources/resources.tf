resource "akamai_gtm_resource" "test_resource1" {
  domain                         = akamai_gtm_domain.test_name.name
  name                           = "test resource1"
  host_header                    = "header"
  type                           = "XML load object via HTTP"
  aggregation_type               = "latest"
  least_squares_decay            = 30
  upper_bound                    = 20
  description                    = "some description"
  leader_string                  = "leader"
  constrained_property           = "**"
  load_imbalance_percentage      = 51
  max_u_multiplicative_increment = 10
  decay_rate                     = 5

  resource_instance {
    datacenter_id    = akamai_gtm_datacenter.TEST1.datacenter_id
    load_object      = "load"
    load_servers     = ["server"]
    load_object_port = 80
  }

  depends_on = [
    akamai_gtm_datacenter.TEST1,
    akamai_gtm_domain.test_name
  ]
}

resource "akamai_gtm_resource" "test_resource2" {
  domain           = akamai_gtm_domain.test_name.name
  name             = "test resource2"
  type             = ""
  aggregation_type = ""

  depends_on = [
    akamai_gtm_domain.test_name
  ]
}

