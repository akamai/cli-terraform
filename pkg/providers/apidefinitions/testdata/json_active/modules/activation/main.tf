terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 8.1.0"
    }
  }
  required_version = ">= 1.0"
}

resource "akamai_apidefinitions_activation" "pet_store" {
  api_id                    = var.api_id
  version                   = var.api_version
  network                   = var.network
  notification_recipients   = var.notification_recipients
  notes                     = var.notes
  auto_acknowledge_warnings = var.auto_acknowledge_warnings
}
