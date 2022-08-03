terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_gtm_domain" "test_name" {
  contract                  = var.contractid
  group                     = var.groupid
  name                      = "test.name.akadns.net"
  type                      = "basic"
  comment                   = "test"
  email_notification_list   = ["john@akamai.com", "jdoe@akamai.com"]
  default_timeout_penalty   = 10
  load_imbalance_percentage = 50
  default_error_penalty     = 90
  cname_coalescing_enabled  = true
  load_feedback             = true
  end_user_mapping_enabled  = false
}
