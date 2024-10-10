resource "akamai_iam_user" "iam_user_123" {
  first_name         = "John"
  last_name          = "Smith"
  email              = "terraform@akamai.com"
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
  auth_grants_json   = "[{\"groupId\":56789,\"isBlocked\":false,\"roleId\":12345}]"
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

