resource "akamai_appsec_configuration" "config" {
  name        = var.name
  description = var.description
  contract_id = var.contract_id
  group_id    = trimprefix(data.akamai_group.group.id, "grp_")
  host_names  = var.hostnames
}

data "akamai_group" "group" {
  group_name  = var.group_name
  contract_id = var.contract_id
}

locals {
  config_id = akamai_appsec_configuration.config.config_id
}

output "config_id" {
  value = local.config_id
}
