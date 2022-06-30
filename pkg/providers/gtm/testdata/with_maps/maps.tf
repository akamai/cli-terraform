resource "akamai_gtm_cidrmap" "test_cidrmap" {
  domain = akamai_gtm_domain.test_name.name
  default_datacenter {
    nickname      = "default"
    datacenter_id = akamai_gtm_datacenter.TEST2.datacenter_id
  }
  name = "test_cidrmap"
  depends_on = [
    akamai_gtm_domain.test_name
  ]
}

resource "akamai_gtm_geomap" "test_geomap" {
  domain = akamai_gtm_domain.test_name.name
  default_datacenter {
    nickname      = "default"
    datacenter_id = akamai_gtm_datacenter.TEST2.datacenter_id
  }
  assignment {
    nickname      = "TEST1"
    datacenter_id = akamai_gtm_datacenter.TEST1.datacenter_id
    countries     = ["US"]
  }
  name = "test_geomap"
  depends_on = [
    akamai_gtm_datacenter.TEST1,
    akamai_gtm_domain.test_name
  ]
}

resource "akamai_gtm_asmap" "test_asmap" {
  domain = akamai_gtm_domain.test_name.name
  default_datacenter {
    nickname      = "default"
    datacenter_id = akamai_gtm_datacenter.TEST1.datacenter_id
  }
  assignment {
    nickname      = "TEST1"
    datacenter_id = akamai_gtm_datacenter.TEST1.datacenter_id
    as_numbers    = [1, 2, 3]
  }
  name = "test_asmap"
  depends_on = [
    akamai_gtm_datacenter.TEST1,
    akamai_gtm_domain.test_name
  ]
}
