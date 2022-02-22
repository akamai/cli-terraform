data "akamai_cloudlets_audience_segmentation_match_rule" "match_rules_as" {
  match_rules {
    name  = "rule1"
    start = 0
    end   = 0
    matches {
      match_type     = "method"
      match_value    = ""
      match_operator = "equals"
      case_sensitive = true
      negate         = false
      check_ips      = ""
      object_match_value {
        type  = "simple"
        value = ["GET"]
      }
    }
    match_url = "test.url"
    forward_settings {
      origin_id                 = ""
      path_and_qs               = "some_path"
      use_incoming_query_string = false
    }
    disabled = false
  }

  match_rules {
    name  = "rule2"
    start = 0
    end   = 0
    matches {
      match_type     = "header"
      match_value    = ""
      match_operator = "equals"
      case_sensitive = false
      negate         = false
      check_ips      = ""
      object_match_value {
        name                = "AS"
        type                = "object"
        name_case_sensitive = false
        name_has_wildcard   = false
        options {
          value                = ["y"]
          value_has_wildcard   = true
          value_case_sensitive = false
          value_escaped        = false
        }
      }
    }
    match_url = "abc.com"
    forward_settings {
      origin_id                 = ""
      path_and_qs               = ""
      use_incoming_query_string = true
    }
    disabled = false
  }

  match_rules {
    name  = "rule3"
    start = 1
    end   = 2
    matches {
      match_type     = "range"
      match_value    = ""
      match_operator = "equals"
      case_sensitive = false
      negate         = false
      check_ips      = ""
      object_match_value {
        type  = "range"
        value = [1, 50]
      }
    }
    match_url = "test.url"
    forward_settings {
      origin_id                 = "test_origin"
      path_and_qs               = ""
      use_incoming_query_string = false
    }
    disabled = false
  }

  match_rules {
    name      = "rule_empty"
    start     = 0
    end       = 0
    match_url = ""
    forward_settings {
      origin_id                 = ""
      path_and_qs               = ""
      use_incoming_query_string = false
    }
    disabled = true
  }
}
