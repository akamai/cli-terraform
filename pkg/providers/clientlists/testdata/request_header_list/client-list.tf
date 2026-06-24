terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 10.3.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_clientlist_list" "list_123_ABC" {
  name  = "Test Request Header Client List"
  type  = "REQUEST_HEADER_NAME_VALUE"
  notes = "Request header notes"
  tags  = ["ListTag3"]

  contract_id = var.contract_id
  group_id    = var.group_id

  dynamic "items" {
    for_each = jsondecode(file("./123_ABC.json"))

    content {
      key             = items.value.key
      values          = items.value.values
      description     = items.value.description
      tags            = items.value.tags
      expiration_date = items.value.expirationDate
    }
  }
}

# resource "akamai_clientlist_activation" "activation_123_ABC_STAGING" {
#   list_id                 = akamai_clientlist_list.list_123_ABC.list_id
#   network                 = "STAGING"
#   comments                = ""
#   notification_recipients = []
#   siebel_ticket_id        = ""
# }

# resource "akamai_clientlist_activation" "activation_123_ABC_PRODUCTION" {
#   list_id                 = akamai_clientlist_list.list_123_ABC.list_id
#   network                 = "PRODUCTION"
#   comments                = ""
#   notification_recipients = []
#   siebel_ticket_id        = ""
# }