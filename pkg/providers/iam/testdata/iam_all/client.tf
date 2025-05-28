resource "akamai_iam_api_client" "api_client_1a2b3" {
  authorized_users           = ["mw+2"]
  can_auto_create_credential = false
  allow_account_switch       = false
  client_type                = "CLIENT"
  client_name                = "mw+2_1"
  notification_emails        = ["mw+2@example.com"]
  client_description         = "Test API Client"
  lock                       = false
  credential = {
    description = "Test API Client Credential 2"
    expires_on  = "2025-06-13T14:48:08Z"
    status      = "ACTIVE"
  }
  group_access = {
    clone_authorized_user_groups = false
    groups = [
      {
        group_id = 123
        role_id  = 340
        sub_groups = [
          {
            group_id = 333
            role_id  = 540
            sub_groups = [
              {
                group_id = 444
                role_id  = 640
              }
            ]
          }
        ]
    }]
  }
  ip_acl = {
    enable = true
    cidr   = ["128.5.6.5/24"]
  }
  api_access = {
    all_accessible_apis = false
    apis = [
      {
        api_id       = 5801
        access_level = "READ-WRITE"
      },
      {
        api_id       = 5580
        access_level = "READ-ONLY"
    }]
  }
  purge_options = {
    can_purge_by_cp_code   = true
    can_purge_by_cache_tag = true
    cp_code_access = {
      all_current_and_new_cp_codes = false
      cp_codes                     = [101]
    }
  }
}