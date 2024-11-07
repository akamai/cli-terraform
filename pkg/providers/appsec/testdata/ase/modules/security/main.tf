data "akamai_appsec_configuration" "config" {
  name = var.name
}

data "akamai_group" "group" {
  group_name  = var.group_name
  contract_id = var.contract_id
}

output "config_id" {
  value = akamai_appsec_configuration.config.config_id
}
