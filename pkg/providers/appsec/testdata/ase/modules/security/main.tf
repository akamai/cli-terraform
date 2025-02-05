data "akamai_appsec_configuration" "config" {
  name = var.name
}

data "akamai_group" "group" {
  group_name  = var.group_name
  contract_id = var.contract_id
}

locals {
  config_id = data.akamai_appsec_configuration.config.config_id
}

output "config_id" {
  value = local.config_id
}
