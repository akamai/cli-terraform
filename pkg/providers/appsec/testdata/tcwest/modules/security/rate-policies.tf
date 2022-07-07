resource "akamai_appsec_rate_policy" "high_rate" {
  config_id = akamai_appsec_configuration.config.config_id
  rate_policy = jsonencode(
    {
      "additionalMatchOptions" : null,
      "averageThreshold" : 100,
      "burstThreshold" : 500,
      "clientIdentifier" : "ip",
      "matchType" : "path",
      "name" : "High Rate",
      "pathMatchType" : "Custom",
      "pathUriPositiveMatch" : true,
      "requestType" : "ClientRequest",
      "sameActionOnIpv6" : true,
      "type" : "WAF",
      "useXForwardForHeaders" : false
    }
  )
}

resource "akamai_appsec_rate_policy" "low_rate" {
  config_id = akamai_appsec_configuration.config.config_id
  rate_policy = jsonencode(
    {
      "additionalMatchOptions" : null,
      "averageThreshold" : 75,
      "burstThreshold" : 250,
      "clientIdentifier" : "ip",
      "matchType" : "path",
      "name" : "Low Rate",
      "pathMatchType" : "Custom",
      "pathUriPositiveMatch" : true,
      "requestType" : "ClientRequest",
      "sameActionOnIpv6" : true,
      "type" : "WAF",
      "useXForwardForHeaders" : false
    }
  )
}

