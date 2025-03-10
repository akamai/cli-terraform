resource "akamai_apidefinitions_resource_operations" "pet_store" {
  api_id              = var.api_id
  resource_operations = jsonencode(file("${path.module}/operations-api.json"))
}