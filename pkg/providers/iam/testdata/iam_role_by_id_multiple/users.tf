resource "akamai_iam_user" "iam_user_123" {
  first_name         = "John"
  last_name          = "Smith"
  email              = "terraform@akamai.com"
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
  auth_grants_json   = "[{\"groupId\":56789,\"groupName\":\"Custom group 1\",\"isBlocked\":false,\"roleDescription\":\"Custom role description\",\"roleId\":12345,\"roleName\":\"Custom role\"}]"
  lock               = false
}

resource "akamai_iam_user" "iam_user_321" {
  first_name         = "Steve"
  last_name          = "Smith"
  email              = "terraform_1@akamai.com"
  country            = "Canada"
  phone              = "(617) 444-4650"
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
  auth_grants_json   = "[{\"groupId\":56789,\"groupName\":\"Custom group 1\",\"isBlocked\":false,\"roleDescription\":\"Custom role description\",\"roleId\":12345,\"roleName\":\"Custom role\"},{\"groupId\":98765,\"groupName\":\"Custom group 2\",\"isBlocked\":false,\"roleDescription\":\"Other custom role description\",\"roleId\":54321,\"roleName\":\"Other custom role\"}]"
  lock               = false
}

