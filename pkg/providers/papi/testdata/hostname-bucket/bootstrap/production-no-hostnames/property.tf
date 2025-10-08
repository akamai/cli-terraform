terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 9.1.0"
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

resource "akamai_property_bootstrap" "test-edgesuite-net" {
  name                = "test.edgesuite.net"
  contract_id         = var.contract_id
  group_id            = var.group_id
  product_id          = "prd_HTTP_Content_Del"
  use_hostname_bucket = true
}

resource "akamai_property" "test-edgesuite-net" {
  property_id         = akamai_property_bootstrap.test-edgesuite-net.id
  name                = akamai_property_bootstrap.test-edgesuite-net.name
  contract_id         = akamai_property_bootstrap.test-edgesuite-net.contract_id
  group_id            = akamai_property_bootstrap.test-edgesuite-net.group_id
  product_id          = akamai_property_bootstrap.test-edgesuite-net.product_id
  use_hostname_bucket = akamai_property_bootstrap.test-edgesuite-net.use_hostname_bucket
  rule_format         = "latest"
  rules               = data.akamai_property_rules_template.rules.json
}

# NOTE: Be careful when removing this resource as you can disable traffic
#resource "akamai_property_activation" "test-edgesuite-net-staging" {
#  property_id                    = akamai_property.test-edgesuite-net.id
#  contact                        = []
#  version                        = var.activate_latest_on_staging ? akamai_property.test-edgesuite-net.latest_version : akamai_property.test-edgesuite-net.staging_version
#  network                        = "STAGING"
#  auto_acknowledge_rule_warnings = false
#}

# NOTE: Be careful when removing this resource as you can disable traffic
resource "akamai_property_activation" "test-edgesuite-net-production" {
  property_id                    = akamai_property.test-edgesuite-net.id
  contact                        = ["jsmith@akamai.com"]
  version                        = var.activate_latest_on_production ? akamai_property.test-edgesuite-net.latest_version : akamai_property.test-edgesuite-net.production_version
  network                        = "PRODUCTION"
  auto_acknowledge_rule_warnings = false
}

resource "akamai_property_hostname_bucket" "test-edgesuite-net-hostname-bucket-production" {
  property_id   = akamai_property_activation.test-edgesuite-net-production.property_id
  contract_id   = var.contract_id
  group_id      = var.group_id
  network       = "PRODUCTION"
  note          = "production note"
  notify_emails = ["test@mail.com"]
  hostnames = {
    for hostname in local.hostnames :
    hostname.cname_from => {
      cert_provisioning_type = hostname.cert_provisioning_type
      edge_hostname_id       = hostname.edge_hostname_id
    }
    if hostname.production == true
  }
}