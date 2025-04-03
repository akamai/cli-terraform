locals {
  hostnames = [
    {
      cname_from             = "www.test.cname_from.0.com"
      cert_provisioning_type = "CPS_MANAGED"
      edge_hostname_id       = "ehn_12345"
      staging                = true
      production             = false
    },
    {
      cname_from             = "www.test.cname_from.0.com"
      cert_provisioning_type = "DEFAULT"
      edge_hostname_id       = "ehn_12345"
      staging                = false
      production             = true
    },
  ]
}