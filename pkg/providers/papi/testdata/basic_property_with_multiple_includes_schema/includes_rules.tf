
data "akamai_property_rules_builder" "test_include_rule_default" {
  rules_v2023_01_05 {
    name      = "default"
    is_secure = false
    uuid      = "default"
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
          created_date = 1506429558000
          description  = "Test-NewHire"
          id           = 626358
          name         = "Test-NewHire"
          products     = ["Site_Defender", ]
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
        xml = trimsuffix(<<EOT
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
        , "\n")
      }
    }
    children = [
      data.akamai_property_rules_builder.test_include_rule_content_compression.json,
      data.akamai_property_rules_builder.test_include_rule_static_content.json,
      data.akamai_property_rules_builder.test_include_rule_dynamic_content.json,
    ]
  }
}

data "akamai_property_rules_builder" "test_include_rule_content_compression" {
  rules_v2023_01_05 {
    name                  = "Content Compression"
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
  }
}

data "akamai_property_rules_builder" "test_include_rule_static_content" {
  rules_v2023_01_05 {
    name                  = "Static Content"
    criteria_must_satisfy = "all"
    criterion {
      file_extension {
        match_case_sensitive = false
        match_operator       = "IS_ONE_OF"
        values               = ["aif", "aiff", "au", "avi", "bin", "bmp", "cab", "carb", "cct", "cdf", "class", "css", "doc", "dcr", "dtd", "exe", "flv", "gcf", "gff", "gif", "grv", "hdml", "hqx", "ico", "ini", "jpeg", "jpg", "js", "mov", "mp3", "nc", "pct", "pdf", "png", "ppc", "pws", "swa", "swf", "txt", "vbs", "w32", "wav", "wbmp", "wml", "wmlc", "wmls", "wmlsc", "xsd", "zip", "webp", "jxr", "hdp", "wdp", "pict", "tif", "tiff", "mid", "midi", "ttf", "eot", "woff", "otf", "svg", "svgz", "jar", "woff2", ]
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

data "akamai_property_rules_builder" "test_include_rule_dynamic_content" {
  rules_v2023_01_05 {
    name                  = "Dynamic Content"
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

data "akamai_property_rules_builder" "test_include_1_rule_default" {
  rules_v2023_01_05 {
    name      = "default"
    is_secure = false
    uuid      = "default"
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
          created_date = 1506429558000
          description  = "Test-NewHire"
          id           = 626358
          name         = "Test-NewHire"
          products     = ["Site_Defender", ]
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
        xml = trimsuffix(<<EOT
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
        , "\n")
      }
    }
    children = [
      data.akamai_property_rules_builder.test_include_1_rule_content_compression.json,
      data.akamai_property_rules_builder.test_include_1_rule_static_content.json,
      data.akamai_property_rules_builder.test_include_1_rule_dynamic_content.json,
    ]
  }
}

data "akamai_property_rules_builder" "test_include_1_rule_content_compression" {
  rules_v2023_01_05 {
    name                  = "Content Compression"
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
  }
}

data "akamai_property_rules_builder" "test_include_1_rule_static_content" {
  rules_v2023_01_05 {
    name                  = "Static Content"
    criteria_must_satisfy = "all"
    criterion {
      file_extension {
        match_case_sensitive = false
        match_operator       = "IS_ONE_OF"
        values               = ["aif", "aiff", "au", "avi", "bin", "bmp", "cab", "carb", "cct", "cdf", "class", "css", "doc", "dcr", "dtd", "exe", "flv", "gcf", "gff", "gif", "grv", "hdml", "hqx", "ico", "ini", "jpeg", "jpg", "js", "mov", "mp3", "nc", "pct", "pdf", "png", "ppc", "pws", "swa", "swf", "txt", "vbs", "w32", "wav", "wbmp", "wml", "wmlc", "wmls", "wmlsc", "xsd", "zip", "webp", "jxr", "hdp", "wdp", "pict", "tif", "tiff", "mid", "midi", "ttf", "eot", "woff", "otf", "svg", "svgz", "jar", "woff2", ]
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

data "akamai_property_rules_builder" "test_include_1_rule_dynamic_content" {
  rules_v2023_01_05 {
    name                  = "Dynamic Content"
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
