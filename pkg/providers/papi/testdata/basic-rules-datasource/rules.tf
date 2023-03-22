
data "akamai_property_rules_builder" "test-edgesuite-net_rule_default" {
  rules_v2023_01_05 {
    name                  = "default"
    is_secure             = false
    criteria_must_satisfy = "all"
    uuid                  = "default"
    variable {
      name        = "PMUSER_TESTSTR"
      description = "DSTR"
      value       = "STR"
      hidden      = false
      sensitive   = true
    }
    variable {
      name        = "PMUSER_TEST100"
      description = "D100"
      value       = "100"
      hidden      = false
      sensitive   = false
    }
    criterion {
      match_advanced {
        uuid        = "fa27bc4d-bfff-4541-8eb7-ade156a57256"
        description = ""
      }
    }
    criterion {
      content_type {
        match_case_sensitive = false
        match_operator       = "IS_ONE_OF"
        match_wildcard       = true
        values               = ["text/html*", "text/css*", "application/x-javascript*", ]
      }
    }
    behavior {
      application_load_balancer {
        all_down_net_storage_file   = ""
        all_down_status_code        = ""
        allow_cache_prefresh        = true
        enabled                     = true
        failover_attempts_threshold = 5
        failover_mode               = "MANUAL"
        failover_origin_map {
          from_origin_id = "dddd"
          to_origin_ids  = ["yyyy", "yyyy1", "yyyy2", ]
        }
        failover_origin_map {
          from_origin_id = "oooo"
          to_origin_ids  = ["xxxxx", ]
        }
        failover_origin_map {
          from_origin_id = "wwww"
          to_origin_ids  = ["zzzzzz", ]
        }
        failover_status_codes                = ["500", "501", "502", "503", "504", "505", "506", "507", "508", "509", ]
        label                                = ""
        stickiness_cookie_automatic_salt     = true
        stickiness_cookie_set_http_only_flag = true
        stickiness_cookie_type               = "ON_BROWSER_CLOSE"
      }
    }
    behavior {
      origin {
        cache_key_hostname    = "ORIGIN_HOSTNAME"
        compress              = true
        enable_true_client_ip = false
        forward_host_header   = "REQUEST_HOST_HEADER"
        hostname              = "1.2.3.4"
        http_port             = 80
        https_port            = 443
        origin_sni            = false
        origin_type           = "CUSTOMER"
        use_unique_cache_key  = false
        verification_mode     = "PLATFORM_SETTINGS"
      }
    }
    behavior {
      cp_code {
        value {
          id = 1047836
        }
      }
    }
    behavior {
      caching {
        behavior = "NO_STORE"
      }
    }
    behavior {
      allow_post {
        allow_without_content_length = false
        enabled                      = true
      }
    }
    behavior {
      report {
        log_accept_language  = false
        log_cookies          = "OFF"
        log_custom_log_field = false
        log_host             = false
        log_referer          = false
        log_user_agent       = true
      }
    }
    behavior {
      advanced {
        uuid        = "feeaeff9-fe7e-4e27-ba0c-7b1dcecdba8b"
        description = "extract inputs"
        xml         = <<-EOT
<assign:extract-value>
   <variable-name>ENDUSER</variable-name>
   <location>Query_String</location>
   <location-id>enduser</location-id>
   <separator>=</separator>
</assign:extract-value>
<assign:extract-value>
   <variable-name>GHOST</variable-name>
   <location>Query_String</location>
   <location-id>ghost</location-id>
   <separator>=</separator>
</assign:extract-value>

<assign:variable>
   <name>DISTANCE</name>
   <transform>
      <geo-distance>
         <ip1>%(ENDUSER)</ip1>
         <ip2>%(GHOST)</ip2>
      </geo-distance>
   </transform>
</assign:variable>



<edgeservices:construct-response>
   <status>on</status>
   <http-status>200</http-status>
   <body>%(DISTANCE)</body>
   <force-cache-eviction>off</force-cache-eviction>
</edgeservices:construct-response>

<edgeservices:modify-outgoing-response.add-header>
      <name>Distance</name>
      <value>%(DISTANCE)</value>
   </edgeservices:modify-outgoing-response.add-header>
EOT
      }
    }
    behavior {
      fail_action {
        action_type = "RECREATED_NS"
        cp_code {
          id = 192729
        }
        enabled = true
        net_storage_hostname {
          cp_code              = 196797
          download_domain_name = "spm.download.akamai.com"
        }
        net_storage_path = "/pathto/sorry_page.html"
        status_code      = 200
      }
    }
    children = [
      data.akamai_property_rules_builder.test-edgesuite-net_rule_strange_characters--a-------------ą.json,
      data.akamai_property_rules_builder.test-edgesuite-net_rule_static_content.json,
      data.akamai_property_rules_builder.test-edgesuite-net_rule_dynamic_content.json,
      data.akamai_property_rules_builder.test-edgesuite-net_rule_newrule.json,
      data.akamai_property_rules_builder.test-edgesuite-net_rule_newrule1.json,
    ]
  }
}

data "akamai_property_rules_builder" "test-edgesuite-net_rule_strange_characters--a-------------ą" {
  rules_v2023_01_05 {
    name                  = "Strange Characters${a}\"\\||$%&*@#|!ą"
    is_secure             = false
    criteria_must_satisfy = "all"
    criterion {
      content_type {
        match_case_sensitive = false
        match_operator       = "IS_ONE_OF"
        match_wildcard       = true
        values               = ["text/html*", "text/css*", "application/x-javascript*", ]
      }
    }
    behavior {
      gzip_response {
        behavior = "ALWAYS"
      }
    }
    children = [
      data.akamai_property_rules_builder.test-edgesuite-net_rule_newrule2.json,
      data.akamai_property_rules_builder.test-edgesuite-net_rule_newrule3.json,
      data.akamai_property_rules_builder.test-edgesuite-net_rule_strange_characters--a-------------ą1.json,
      data.akamai_property_rules_builder.test-edgesuite-net_rule_m_pulse.json,
    ]
  }
}

data "akamai_property_rules_builder" "test-edgesuite-net_rule_static_content" {
  rules_v2023_01_05 {
    name                  = "Static Content"
    is_secure             = false
    criteria_must_satisfy = "all"
    criterion {
      file_extension {
        match_case_sensitive = false
        match_operator       = "IS_ONE_OF"
        values               = ["au", "avi", "bin", "bmp", "cab", "carb", "cct", "cdf", "class", "css", "doc", "dcr", "dtd", "exe", "flv", "gcf", "gff", "gif", "grv", "hdml", "hqx", "ico", "ini", "jpeg", "jpg", "js", "mov", "mp3", "nc", "pct", "pdf", "png", "ppc", "pws", "swa", "swf", "txt", "vbs", "w32", "wav", "wbmp", "wml", "wmlc", "wmls", "wmlsc", "xsd", "zip", "webp", "jxr", "hdp", "wdp", "pict", "tif", "tiff", "mid", "midi", "ttf", "eot", "woff", "otf", "svg", "svgz", "jar", "woff2", ]
      }
    }
    criterion {
      file_extension {
        match_case_sensitive = false
        match_operator       = "IS_ONE_OF"
        values               = ["aif", "aiff", ]
      }
    }
    behavior {
      caching {
        behavior        = "MAX_AGE"
        must_revalidate = false
        ttl             = "1d"
      }
    }
  }
}

data "akamai_property_rules_builder" "test-edgesuite-net_rule_dynamic_content" {
  rules_v2023_01_05 {
    name                  = "Dynamic Content"
    is_secure             = false
    criteria_must_satisfy = "all"
    criterion {
      cacheability {
        match_operator = "IS_NOT"
        value          = "CACHEABLE"
      }
    }
    behavior {
      downstream_cache {
        behavior = "TUNNEL_ORIGIN"
      }
    }
  }
}

data "akamai_property_rules_builder" "test-edgesuite-net_rule_newrule" {
  rules_v2023_01_05 {
    name                  = "new rule"
    is_secure             = false
    criteria_must_satisfy = "all"
  }
}

data "akamai_property_rules_builder" "test-edgesuite-net_rule_newrule1" {
  rules_v2023_01_05 {
    name                  = "new rule"
    is_secure             = false
    criteria_must_satisfy = "any"
  }
}

data "akamai_property_rules_builder" "test-edgesuite-net_rule_newrule2" {
  rules_v2023_01_05 {
    name                  = "new rule"
    is_secure             = false
    criteria_must_satisfy = "all"
  }
}

data "akamai_property_rules_builder" "test-edgesuite-net_rule_newrule3" {
  rules_v2023_01_05 {
    name                  = "new rule"
    is_secure             = false
    criteria_must_satisfy = "all"
  }
}

data "akamai_property_rules_builder" "test-edgesuite-net_rule_strange_characters--a-------------ą1" {
  rules_v2023_01_05 {
    name                  = "Strange Characters${a}\"\\&&$%&*@#|!ą"
    is_secure             = false
    criteria_must_satisfy = "all"
  }
}

data "akamai_property_rules_builder" "test-edgesuite-net_rule_m_pulse" {
  rules_v2023_01_05 {
    name                  = "mPulse"
    is_secure             = false
    comments              = "Test mPulse"
    criteria_must_satisfy = "all"
    behavior {
      m_pulse {
        buffer_size     = ""
        config_override = <<-EOT
{"name":"John", "age":30, "car":null}
EOT
        enabled         = true
        loader_version  = "V12"
        require_pci     = true
      }
    }
  }
}
