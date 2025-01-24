resource "akamai_appsec_rate_policy" "page_view_requests" {
  config_id = local.config_id
  rate_policy = jsonencode(
    {
      "additionalMatchOptions" : [
        {
          "positiveMatch" : true,
          "type" : "RequestMethodCondition",
          "values" : [
            "GET",
            "HEAD"
          ]
        }
      ],
      "averageThreshold" : 12,
      "burstThreshold" : 18,
      "clientIdentifiers" : [
        "ip"
      ],
      "description" : "A popular brute force attack that consists of sending a large number of requests for base page, HTML page or XHR requests (usually non-cacheable). This could destabilize the origin.",
      "fileExtensions" : {
        "positiveMatch" : false,
        "values" : [
          "aif",
          "aiff",
          "au",
          "avi",
          "bin",
          "bmp",
          "cab",
          "carb",
          "cct",
          "cdf",
          "class",
          "css",
          "doc",
          "dcr",
          "dtd",
          "exe",
          "flv",
          "gcf",
          "gff",
          "gif",
          "grv",
          "hdml",
          "hqx",
          "ico",
          "ini",
          "jpeg",
          "jpg",
          "js",
          "mov",
          "mp3",
          "nc",
          "pct",
          "pdf",
          "png",
          "ppc",
          "pws",
          "svg",
          "swa",
          "swf",
          "txt",
          "vbs",
          "w32",
          "wav",
          "wbmp",
          "wml",
          "wmlc",
          "wmls",
          "wmlsc",
          "xsd",
          "zip",
          "webp",
          "jxr",
          "hdp",
          "wdp",
          "webm",
          "ogv",
          "mp4",
          "ttf",
          "woff",
          "eot",
          "woff2"
        ]
      },
      "matchType" : "path",
      "name" : "Page View Requests",
      "pathMatchType" : "Custom",
      "pathUriPositiveMatch" : true,
      "requestType" : "ClientRequest",
      "sameActionOnIpv6" : true,
      "type" : "WAF",
      "useXForwardForHeaders" : false
    }
  )
}

resource "akamai_appsec_rate_policy" "origin_error" {
  config_id = local.config_id
  rate_policy = jsonencode(
    {
      "additionalMatchOptions" : [
        {
          "positiveMatch" : true,
          "type" : "ResponseStatusCondition",
          "values" : [
            "400",
            "401",
            "402",
            "403",
            "404",
            "405",
            "406",
            "407",
            "408",
            "409",
            "410",
            "500",
            "501",
            "502",
            "503",
            "504"
          ]
        }
      ],
      "averageThreshold" : 5,
      "burstThreshold" : 8,
      "clientIdentifiers" : [
        "ip"
      ],
      "description" : "An excessive error rate from the origin could indicate malicious activity by a bot scanning the site or a publishing error. In both cases, this would increase the origin traffic and could potentially destabilize it.",
      "matchType" : "path",
      "name" : "Origin Error",
      "pathMatchType" : "Custom",
      "pathUriPositiveMatch" : true,
      "requestType" : "ForwardResponse",
      "sameActionOnIpv6" : true,
      "type" : "WAF",
      "useXForwardForHeaders" : false
    }
  )
}

resource "akamai_appsec_rate_policy" "post_page_requests" {
  config_id = local.config_id
  rate_policy = jsonencode(
    {
      "additionalMatchOptions" : [
        {
          "positiveMatch" : true,
          "type" : "RequestMethodCondition",
          "values" : [
            "POST"
          ]
        }
      ],
      "averageThreshold" : 3,
      "burstThreshold" : 5,
      "clientIdentifiers" : [
        "ip"
      ],
      "description" : "Mitigating HTTP flood attacks using POST requests",
      "matchType" : "path",
      "name" : "POST Page Requests",
      "pathMatchType" : "Custom",
      "pathUriPositiveMatch" : true,
      "requestType" : "ClientRequest",
      "sameActionOnIpv6" : true,
      "type" : "WAF",
      "useXForwardForHeaders" : false
    }
  )
}

