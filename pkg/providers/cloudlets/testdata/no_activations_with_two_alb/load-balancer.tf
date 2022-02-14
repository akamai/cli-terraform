resource "akamai_cloudlets_application_load_balancer" "load_balancer_test_origin" {
  origin_id      = "test_origin"
  description    = "test description"
  balancing_type = "WEIGHTED"

  data_centers {
    latitude                          = 102.78108
    longitude                         = -116.07064
    continent                         = "NA"
    country                           = "US"
    origin_id                         = "test_origin"
    percent                           = 10
    cloud_service                     = true
    liveness_hosts                    = ["tf1.test", "tf2.test"]
    hostname                          = "test-hostname"
    state_or_province                 = "MA"
    city                              = "Boston"
    cloud_server_host_header_override = false
  }

  liveness_settings {
    port        = 1234
    protocol    = "HTTP"
    path        = "status"
    host_header = "header"
    additional_headers = {
      abc = "123"
    }
    interval                      = 10
    peer_certificate_verification = false
    request_string                = "test_request_string"
    response_string               = "test_response_string"
    status_3xx_failure            = false
    status_4xx_failure            = false
    status_5xx_failure            = false
    timeout                       = 60
  }
}

resource "akamai_cloudlets_application_load_balancer" "load_balancer_test_origin_2" {
  origin_id      = "test_origin_2"
  description    = "test description"
  balancing_type = "WEIGHTED"

  data_centers {
    latitude                          = 102.78108
    longitude                         = -116.07064
    continent                         = "NA"
    country                           = "US"
    origin_id                         = "test_origin"
    percent                           = 10
    cloud_service                     = true
    liveness_hosts                    = ["tf1.test", "tf2.test"]
    hostname                          = "test-hostname"
    state_or_province                 = "MA"
    city                              = "Boston"
    cloud_server_host_header_override = false
  }

  liveness_settings {
    port        = 1234
    protocol    = "HTTP"
    path        = "status"
    host_header = "header"
    additional_headers = {
      abc = "123"
    }
    interval                      = 10
    peer_certificate_verification = false
    request_string                = "test_request_string"
    response_string               = "test_response_string"
    status_3xx_failure            = false
    status_4xx_failure            = false
    status_5xx_failure            = false
    timeout                       = 60
  }
}

resource "akamai_cloudlets_application_load_balancer_activation" "load_balancer_activation_test_origin" {
  origin_id = akamai_cloudlets_application_load_balancer.load_balancer_test_origin.origin_id
  network   = var.env
  version   = akamai_cloudlets_application_load_balancer.load_balancer_test_origin.version
}

resource "akamai_cloudlets_application_load_balancer_activation" "load_balancer_activation_test_origin_2" {
  origin_id = akamai_cloudlets_application_load_balancer.load_balancer_test_origin_2.origin_id
  network   = var.env
  version   = akamai_cloudlets_application_load_balancer.load_balancer_test_origin_2.version
}

