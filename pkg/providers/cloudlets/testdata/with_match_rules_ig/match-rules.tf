data "akamai_cloudlets_request_control_match_rule" "match_rules_ig" {
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
    allow_deny     = "allow"
    matches_always = false
    disabled       = false
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
        name                = "Accept"
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
    allow_deny     = "allow"
    matches_always = false
    disabled       = false
  }

  match_rules {
    name           = "rule_empty"
    start          = 0
    end            = 0
    allow_deny     = "deny"
    matches_always = true
    disabled       = true
  }
}
