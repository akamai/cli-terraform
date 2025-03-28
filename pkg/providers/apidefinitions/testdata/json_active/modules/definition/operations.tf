resource "akamai_apidefinitions_resource_operations" "pet_store" {
  api_id              = akamai_apidefinitions_api.pet_store.id
  resource_operations = file("${path.module}/definition-operations-api.json")
}