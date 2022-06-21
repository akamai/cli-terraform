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

resource "akamai_iam_user" "iam_user_123" {
  first_name         = "Terraform"
  last_name          = "Test"
  email              = "terraform8@akamai.com"
  country            = "Canada"
  phone              = "(617) 444-4649"
  enable_tfa         = true
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
  auth_grants_json   = "[{\"groupId\":56789,\"groupName\":\"Custom group 56789\",\"isBlocked\":false,\"roleDescription\":\"Custom role description\",\"roleId\":12345,\"roleName\":\"Custom role 12345\"},{\"groupId\":987,\"groupName\":\"Custom group 987\",\"isBlocked\":false,\"roleDescription\":\"Custom role description\",\"roleId\":54321,\"roleName\":\"Custom role 54321\"}]"
  lock               = false
}

