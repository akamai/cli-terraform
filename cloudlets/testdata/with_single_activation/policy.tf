terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cloudlets_policy" "policy" {
  name = "test_policy_export"
  cloudlet_code = "ER"
  description = "Testing exported policy"
  group_id = "12345"
  match_rule_format = "1.0"
  match_rules = data.akamai_cloudlets_edge_redirector_match_rule.match_rules_er.json
}

data "akamai_cloudlets_edge_redirector_match_rule" "match_rules_er" {
  match_rules {
    name = "r1"
    start = 1
    end = 2
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
      match_type = "extension"
      match_value = "txt"
      match_operator = "equals"
      case_sensitive = false
      negate = false
      check_ips = ""
    }
    matches {
      match_type = "cookie"
      match_value = "cookie=cookievalue"
      match_operator = "equals"
      case_sensitive = true
      negate = false
      check_ips = ""
    }
    matches {
      match_type = "hostname"
      match_value = "3333.dom"
      match_operator = "equals"
      case_sensitive = true
      negate = true
      check_ips = ""
    }
    use_relative_url = "copy_scheme_hostname"
    status_code = 307
    redirect_url = "/abc/sss"
    match_url = "test.url"
    use_incoming_query_string = false
    use_incoming_scheme_and_host = true
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
        name = "ALB"
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
    use_relative_url = "copy_scheme_hostname"
    status_code = 301
    redirect_url = "/ddd"
    match_url = "abc.com"
    use_incoming_query_string = false
    use_incoming_scheme_and_host = true
  }
}

resource "akamai_cloudlets_policy_activation" "policy_activation_prod" {
  policy_id = 2
  network = "prod"
  version = 1
  associated_properties = [ "prp_0" ]
}
