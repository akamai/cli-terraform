resource "akamai_iam_api_client" "api_client_1a2b3" {
  authorized_users           = ["mw+2"]
  can_create_auto_credential = false
  allow_account_switch       = false
  client_type                = "CLIENT"
  client_name                = "mw+2_1"
  notification_emails        = ["mw+2@example.com"]
  client_description         = "Test API Client"
  lock                       = false
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
        api_id            = 5801
        api_name          = "EdgeWorkers"
        description       = "EdgeWorkers"
        endpoint          = "/edgeworkers/"
        documentation_url = "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html"
        access_level      = "READ-WRITE"
      },
      {
        api_id            = 5580
        api_name          = "Search Data Feed"
        description       = "Search Data Feed"
        endpoint          = "/search-portal-data-feed-api/"
        documentation_url = "/"
        access_level      = "READ-ONLY"
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