terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 6.4.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

data "akamai_property_rules_template" "rules" {
  template_file = abspath("${path.module}/property-snippets/main.json")
}

resource "akamai_edge_hostname" "test-edgesuite-net" {
  contract_id   = var.contract_id
  group_id      = var.group_id
  ip_behavior   = "IPV6_COMPLIANCE"
  edge_hostname = "test.edgesuite.net"
}

resource "akamai_property" "test-edgesuite-net" {
  name        = "test.edgesuite.net"
  contract_id = var.contract_id
  group_id    = var.group_id
  product_id  = "prd_HTTP_Content_Del"
  hostnames {
    cname_from             = "test.edgesuite.net"
    cname_to               = akamai_edge_hostname.test-edgesuite-net.edge_hostname
    cert_provisioning_type = "CPS_MANAGED"
  }
  rule_format = "latest"
  rules       = data.akamai_property_rules_template.rules.json
}

# NOTE: Be careful when removing this resource as you can disable traffic
resource "akamai_property_activation" "test-edgesuite-net-staging" {
  property_id                    = akamai_property.test-edgesuite-net.id
  contact                        = ["jsmith@akamai.com"]
  version                        = var.activate_latest_on_staging ? akamai_property.test-edgesuite-net.latest_version : akamai_property.test-edgesuite-net.staging_version
  network                        = "STAGING"
  auto_acknowledge_rule_warnings = false
}

# NOTE: Be careful when removing this resource as you can disable traffic
#resource "akamai_property_activation" "test-edgesuite-net-production" {
#  property_id                    = akamai_property.test-edgesuite-net.id
#  contact                        = []
#  version                        = var.activate_latest_on_production ? akamai_property.test-edgesuite-net.latest_version : akamai_property.test-edgesuite-net.production_version
#  network                        = "PRODUCTION"
#  auto_acknowledge_rule_warnings = false
#}
