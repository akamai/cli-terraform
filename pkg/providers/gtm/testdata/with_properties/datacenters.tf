resource "akamai_gtm_datacenter" "TEST1" {
  domain                            = akamai_gtm_domain.test_name.name
  nickname                          = "TEST1"
  city                              = "New York"
  state_or_province                 = "NY"
  country                           = "US"
  latitude                          = 40.71305
  longitude                         = -74.00723
  cloud_server_host_header_override = false
  cloud_server_targeting            = false
  default_load_object {
    load_object      = "test load object"
    load_object_port = 111
    load_servers     = ["loadServer1", "loadServer2", "loadServer3"]
  }
  depends_on = [
    akamai_gtm_domain.test_name
  ]
}

resource "akamai_gtm_datacenter" "TEST2" {
  domain                            = akamai_gtm_domain.test_name.name
  nickname                          = "TEST2"
  city                              = "Chicago"
  state_or_province                 = "IL"
  country                           = "US"
  latitude                          = 41.88323
  longitude                         = -87.6324
  cloud_server_host_header_override = false
  cloud_server_targeting            = false
  depends_on = [
    akamai_gtm_domain.test_name
  ]
}

data "akamai_gtm_default_datacenter" "default_datacenter_5400" {
  domain     = akamai_gtm_domain.test_name.name
  datacenter = 5400
}

