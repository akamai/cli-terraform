terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 8.0.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_iam_user" "iam_user_001" {
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
  auth_grants_json   = "[{\"groupId\":101,\"isBlocked\":false,\"roleId\":201}]"
  lock               = false
  user_notifications {
    api_client_credential_expiry_notification = true
    new_user_notification                     = true
    password_expiry                           = true
    proactive                                 = ["NetStorage", "EdgeScape"]
    upgrade                                   = ["NetStorage"]
    enable_email_notifications                = true
  }
}

resource "akamai_iam_user" "iam_user_002" {
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
  auth_grants_json   = ""
  lock               = false
  user_notifications {
    api_client_credential_expiry_notification = true
    new_user_notification                     = true
    password_expiry                           = true
    proactive                                 = ["NetStorage", "EdgeScape"]
    upgrade                                   = ["NetStorage"]
    enable_email_notifications                = true
  }
}

