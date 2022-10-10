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

data "akamai_gtm_default_datacenter" "default_datacenter_5401" {
  domain     = akamai_gtm_domain.test_name.name
  datacenter = 5401
}

data "akamai_gtm_default_datacenter" "default_datacenter_5402" {
  domain     = akamai_gtm_domain.test_name.name
  datacenter = 5402
}

