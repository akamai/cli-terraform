resource "akamai_appsec_url_protection_policy" "test_50411" {
  config_id          = local.config_id
  name               = "test"
  max_rate_threshold = 10
  bypass_conditions = [
    {
      type   = "NetworkListCondition"
      values = ["106723_1CLICKTEST", "106793_1CLICKTEST12"]
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Accept", "Accept-Datetime"]
      name_wildcard        = true
      value_case_sensitive = false
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Cache-Control", "Accept-Charset"]
      name_wildcard        = false
      value_case_sensitive = false
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Content-Type", "DNT"]
      name_wildcard        = true
      values               = ["test", "test2"]
      value_case_sensitive = true
      value_wildcard       = true
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Connection"]
      name_wildcard        = true
      values               = ["secure"]
      value_case_sensitive = true
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["DNT"]
      name_wildcard        = true
      values               = ["test"]
      value_case_sensitive = false
      value_wildcard       = true
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Pragma", "Content-MD5"]
      name_wildcard        = true
      values               = ["value"]
      value_case_sensitive = false
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Referer", "X-Forwarded-For"]
      name_wildcard        = false
      values               = ["yug", "test"]
      value_case_sensitive = true
      value_wildcard       = true
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Accept-Charset"]
      name_wildcard        = false
      values               = ["hello"]
      value_case_sensitive = true
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Content-Length"]
      name_wildcard        = false
      values               = ["test"]
      value_case_sensitive = true
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Content-Type"]
      name_wildcard        = false
      values               = ["test"]
      value_case_sensitive = false
      value_wildcard       = false
    },
  ]
  api_definitions = [
    {
      api_definition_id   = 609153
      defined_resources   = true
      resource_ids        = []
      undefined_resources = true
    },
    {
      api_definition_id   = 470134
      defined_resources   = false
      resource_ids        = [108964, 108965, 108966, 108967, 108968, 108969, 108970, 108971, 108973, 108974, 108975, 108976, 108977]
      undefined_resources = false
    },
    {
      api_definition_id   = 518153
      defined_resources   = true
      resource_ids        = []
      undefined_resources = false
    },
  ]
  intelligent_load_shedding = {
    hits_per_sec = 7
    categories = [
      "BOTS",
      "CLIENT_REPUTATIONS",
      "CLOUD_PROVIDERS",
      "PLATFORM_DDOS_INTELLIGENCE",
      "PROXIES",
      "TOR_EXIT_NODES",
    ]
    custom_criteria = [
      {
        type           = "CLIENT_LIST"
        list_ids       = ["99175_10KCL", "106791_1CLICKTEST10"]
        positive_match = true
      },
      {
        type           = "CLIENT_LIST"
        list_ids       = ["81381_123", "106857_1CLICKTEST15"]
        positive_match = false
      },
    ]
  }
}

resource "akamai_appsec_url_protection_policy" "test1_50412" {
  config_id          = local.config_id
  name               = "test1"
  max_rate_threshold = 10
  bypass_conditions = [
    {
      type   = "NetworkListCondition"
      values = ["106723_1CLICKTEST", "106793_1CLICKTEST12"]
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Accept", "Accept-Datetime"]
      name_wildcard        = true
      value_case_sensitive = false
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Cache-Control", "Accept-Charset"]
      name_wildcard        = false
      value_case_sensitive = false
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Content-Type", "DNT"]
      name_wildcard        = true
      values               = ["test", "test2"]
      value_case_sensitive = true
      value_wildcard       = true
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Connection"]
      name_wildcard        = true
      values               = ["secure"]
      value_case_sensitive = true
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["DNT"]
      name_wildcard        = true
      values               = ["test"]
      value_case_sensitive = false
      value_wildcard       = true
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Pragma", "Content-MD5"]
      name_wildcard        = true
      values               = ["value"]
      value_case_sensitive = false
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Referer", "X-Forwarded-For"]
      name_wildcard        = false
      values               = ["yug", "test"]
      value_case_sensitive = true
      value_wildcard       = true
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Accept-Charset"]
      name_wildcard        = false
      values               = ["hello"]
      value_case_sensitive = true
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Content-Length"]
      name_wildcard        = false
      values               = ["test"]
      value_case_sensitive = true
      value_wildcard       = false
    },
    {
      type                 = "RequestHeaderCondition"
      names                = ["Content-Type"]
      name_wildcard        = false
      values               = ["test"]
      value_case_sensitive = false
      value_wildcard       = false
    },
  ]
  hostname_paths = [
    {
      hostname = "www.test1.com"
      paths    = ["/*", "//"]
    },
    {
      hostname = "www.test.com"
      paths    = ["/*", "//", "/o"]
    },
  ]
  intelligent_load_shedding = {
    hits_per_sec = 7
    categories = [
      "BOTS",
      "CLIENT_REPUTATIONS",
      "CLOUD_PROVIDERS",
      "PROXIES",
      "TOR_EXIT_NODES",
    ]
    custom_criteria = [
      {
        type           = "CLIENT_LIST"
        list_ids       = ["99175_10KCL", "106791_1CLICKTEST10"]
        positive_match = true
      },
      {
        type           = "CLIENT_LIST"
        list_ids       = ["81381_123", "106857_1CLICKTEST15"]
        positive_match = false
      },
    ]
  }
}

