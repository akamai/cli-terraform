data "akamai_cloudlets_forward_rewrite_match_rule" "match_rules_fr" {
  match_rules {
    name = "r1"
    start = 0
    end = 0
    matches {
      match_type = "cookie"
      match_value = "cookie=cookievalue"
      match_operator = "equals"
      case_sensitive = true
      negate = false
      check_ips = ""
      object_match_value {
        type = "simple"
        value = ["GET"]
      }
    }
    matches {
      match_type = "hostname"
      match_value = "3333.dom"
      match_operator = "equals"
      case_sensitive = true
      negate = true
      check_ips = ""
    }
    match_url = "test.url"
    forward_settings {
      origin_id = "test_origin"
      path_and_qs = "/test"
      use_incoming_query_string = "false"
    }
    disabled = false
  }

  match_rules {
    name = "r2"
    start = 0
    end = 0
    matches {
      match_type = "header"
      match_value = ""
      match_operator = "equals"
      case_sensitive = false
      negate = false
      check_ips = ""
      object_match_value {
        name = "test_omv"
        type = "object"
        name_case_sensitive = false
        name_has_wildcard = false
        options {
          value = ["y"]
          value_has_wildcard = true
          value_case_sensitive = false
          value_escaped = false
        }
      }
    }
    match_url = "abc.com"
    forward_settings {
      origin_id = "test_origin"
      path_and_qs = ""
      use_incoming_query_string = "false"
    }
    disabled = true
  }
}
