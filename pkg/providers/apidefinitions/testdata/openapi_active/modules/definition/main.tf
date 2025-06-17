terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 8.1.0"
    }
  }
  required_version = ">= 1.0"
}

data "akamai_apidefinitions_openapi" "pet_store_openapi" {
  file_path = "${path.module}/api.yml"
}

resource "akamai_apidefinitions_api" "pet_store" {
  api         = data.akamai_apidefinitions_openapi.pet_store_openapi.api
  contract_id = var.contract_id
  group_id    = var.group_id
}

output "api_id" {
  value = akamai_apidefinitions_api.pet_store.id
}

output "api_latest_version" {
  value = akamai_apidefinitions_api.pet_store.latest_version
}

output "api_staging_version" {
  value = akamai_apidefinitions_api.pet_store.staging_version
}

output "api_production_version" {
  value = akamai_apidefinitions_api.pet_store.production_version
}

resource "akamai_apidefinitions_resource_operations" "pet_store" {
  api_id              = akamai_apidefinitions_api.pet_store.id
  resource_operations = file("${path.module}/operations-api.json")
}