data "akamai_cloudlets_visitor_prioritization_match_rule" "match_rules_vp" {
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
    pass_through_percent = 100
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
        name = "VP"
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
    pass_through_percent = -1
    disabled = false
  }

  match_rules {
    name = "r3"
    start = 0
    end = 0
    match_url = ""
    pass_through_percent = 50.55
    disabled = true
  }
}
