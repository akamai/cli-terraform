terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 7.1.0"
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

resource "akamai_property" "test-edgesuite-net" {
  name        = "test.edgesuite.net"
  contract_id = var.contract_id
  group_id    = var.group_id
  product_id  = "prd_HTTP_Content_Del"
  rule_format = "latest"
  hostnames {
    cname_from             = "test.edgesuite.net"
    cname_to               = "test.edgesuite.net"
    cert_provisioning_type = "CPS_MANAGED"
  }
  rules = data.akamai_property_rules_template.rules.json
}

# NOTE: Be careful when removing this resource as you can disable traffic
#resource "akamai_property_activation" "test-edgesuite-net-staging" {
#  property_id                    = akamai_property.test-edgesuite-net.id
#  contact                        = ["jsmith@akamai.com"]
#  version                        = var.activate_latest_on_staging ? akamai_property.test-edgesuite-net.latest_version : akamai_property.test-edgesuite-net.staging_version
#  network                        = staging
#  auto_acknowledge_rule_warnings = false
#}

# NOTE: Be careful when removing this resource as you can disable traffic
#resource "akamai_property_activation" "test-edgesuite-net-production" {
#  property_id                    = akamai_property.test-edgesuite-net.id
#  contact                        = ["jsmith@akamai.com"]
#  version                        = var.activate_latest_on_production ? akamai_property.test-edgesuite-net.latest_version : akamai_property.test-edgesuite-net.production_version
#  network                        = staging
#  auto_acknowledge_rule_warnings = false
#}
