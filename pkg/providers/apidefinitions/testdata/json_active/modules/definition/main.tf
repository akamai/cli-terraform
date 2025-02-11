terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 6.5.0"
    }
  }
  required_version = ">= 1.0"
}

resource "akamai_apidefinitions_api" "pet_store" {
  api = file("${path.module}/api.json")
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