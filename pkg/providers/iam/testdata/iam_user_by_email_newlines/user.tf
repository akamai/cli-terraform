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
  first_name      = "Terraform\n newline"
  last_name       = "Test\n newline"
  email           = "terraform8@akamai.com"
  country         = "Canada"
  phone           = "(617) 444-4649"
  enable_tfa      = true
  contact_type    = "Technical Decision Maker"
  job_title       = "job\n title "
  time_zone       = "GMT"
  secondary_email = "secondary-email-a@akamai.net"
  mobile_phone    = "(617) 444-4649"
  address = trimsuffix(<<EOT
123
A
Street
EOT
  , "\n")
  city               = "A-Town"
  state              = "TBD"
  zip_code           = "34567"
  preferred_language = "English"
  session_timeout    = 900
  auth_grants_json   = "[{\"groupId\":56789,\"groupName\":\"Custom group\",\"isBlocked\":false,\"roleId\":12345}]"
  lock               = false
}

