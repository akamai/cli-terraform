terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 5.4.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_clientlist_list" "list_123_ABC" {
  name  = "Test Client List"
  type  = "IP"
  notes = "Some Notes"
  tags  = ["tag1", "tag2"]

  contract_id = var.contract_id
  group_id    = var.group_id

  dynamic "items" {
    for_each = jsondecode(file("./123_ABC.json"))

    content {
      value           = items.value.value
      description     = items.value.description
      tags            = items.value.tags
      expiration_date = items.value.expirationDate
    }
  }
}

# resource "akamai_clientlist_activation" "activation_123_ABC_STAGING" {
#   list_id                 = akamai_clientlist_list.list_123_ABC.list_id
#   version                 = akamai_clientlist_list.list_123_ABC.version
#   network                 = "STAGING"
#   comments                = ""
#   notification_recipients = []
#   siebel_ticket_id        = ""
# }

# resource "akamai_clientlist_activation" "activation_123_ABC_PRODUCTION" {
#   list_id                 = akamai_clientlist_list.list_123_ABC.list_id
#   version                 = akamai_clientlist_list.list_123_ABC.version
#   network                 = "PRODUCTION"
#   comments                = ""
#   notification_recipients = []
#   siebel_ticket_id        = ""
# }