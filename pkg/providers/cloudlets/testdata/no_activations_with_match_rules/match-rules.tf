data "akamai_cloudlets_edge_redirector_match_rule" "match_rules_er" {
  match_rules {
    name  = "r1"
    start = 1
    end   = 2
    matches {
      match_type     = "extension"
      match_value    = "txt"
      match_operator = "equals"
      case_sensitive = false
      negate         = false
      check_ips      = ""
    }
    matches {
      match_type     = "cookie"
      match_value    = "cookie=cookievalue"
      match_operator = "equals"
      case_sensitive = true
      negate         = false
      check_ips      = ""
    }
    matches {
      match_type     = "hostname"
      match_value    = "3333.dom"
      match_operator = "equals"
      case_sensitive = true
      negate         = true
      check_ips      = ""
    }
    use_relative_url          = "copy_scheme_hostname"
    status_code               = 307
    redirect_url              = "/abc/sss"
    match_url                 = "test.url"
    use_incoming_query_string = false
    disabled                  = false
  }

  match_rules {
    name                      = "r2"
    start                     = 0
    end                       = 0
    use_relative_url          = "copy_scheme_hostname"
    status_code               = 301
    redirect_url              = "/ddd"
    match_url                 = "abc.com"
    use_incoming_query_string = false
    disabled                  = false
  }
}
