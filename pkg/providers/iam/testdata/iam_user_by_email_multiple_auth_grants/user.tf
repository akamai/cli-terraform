terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_iam_user" "iam_user_123" {
  first_name         = "Terraform"
  last_name          = "Test"
  email              = "terraform8@akamai.com"
  country            = "Canada"
  phone              = "(617) 444-4649"
  enable_tfa         = true
  enable_mfa         = false
  contact_type       = "Technical Decision Maker"
  job_title          = "job title "
  time_zone          = "GMT"
  secondary_email    = "secondary-email-a@akamai.net"
  mobile_phone       = "(617) 444-4649"
  address            = "123 A Street"
  city               = "A-Town"
  state              = "TBD"
  zip_code           = "34567"
  preferred_language = "English"
  session_timeout    = 900
  auth_grants_json   = "[{\"groupId\":56789,\"isBlocked\":false,\"roleId\":12345},{\"groupId\":987,\"isBlocked\":false,\"roleId\":54321}]"
  lock               = false
  user_notifications {
    api_client_credential_expiry_notification = false
    new_user_notification                     = false
    password_expiry                           = false
    proactive                                 = []
    upgrade                                   = []
    enable_email_notifications                = false
  }
}

