resource "akamai_appsec_match_target" "website_4092331" {
  config_id = akamai_appsec_configuration.config.config_id
  match_target = jsonencode(
    {
      "defaultFile" : "NO_MATCH",
      "filePaths" : [
        "/*"
      ],
      "hostnames" : [
        "www.easyakamai.com",
        "www.andrew89.com"
      ],
      "isNegativeFileExtensionMatch" : false,
      "isNegativePathMatch" : false,
      "securityPolicy" : {
        "policyId" : "${akamai_appsec_security_policy.andrew.security_policy_id}"
      },
      "sequence" : 0,
      "type" : "website"
    }
  )
}
resource "akamai_appsec_match_target" "website_2034325" {
  config_id = akamai_appsec_configuration.config.config_id
  match_target = jsonencode(
    {
      "defaultFile" : "NO_MATCH",
      "filePaths" : [
        "/andrewaaa"
      ],
      "hostnames" : [
        "www.vbhat.com"
      ],
      "isNegativeFileExtensionMatch" : false,
      "isNegativePathMatch" : false,
      "bypassNetworkLists" : [
        {
          "id" : "46506_AKAMAIPORTALTESTCENTERA",
          "name" : "Akamai Test Center (ATC) Agents"
        }
      ],
      "securityPolicy" : {
        "policyId" : "${akamai_appsec_security_policy.policy1.security_policy_id}"
      },
      "sequence" : 0,
      "type" : "website"
    }
  )
}
resource "akamai_appsec_match_target" "website_4092261" {
  config_id = akamai_appsec_configuration.config.config_id
  match_target = jsonencode(
    {
      "defaultFile" : "NO_MATCH",
      "filePaths" : [
        "/*"
      ],
      "hostnames" : [
        "konaneweahost8012.edgekey.net",
        "konaneweahost9001.edgekey.net",
        "konaneweahost8013.edgekey.net",
        "konaneweahost9002.edgekey.net",
        "konaneweahost8014.edgekey.net",
        "konaneweahost8016.edgekey.net",
        "aetsaitcwest.edgekey.net",
        "konaneweahost9000.edgekey.net"
      ],
      "isNegativeFileExtensionMatch" : false,
      "isNegativePathMatch" : false,
      "bypassNetworkLists" : [
        {
          "id" : "46506_AKAMAIPORTALTESTCENTERA",
          "name" : "Akamai Test Center (ATC) Agents"
        }
      ],
      "securityPolicy" : {
        "policyId" : "${akamai_appsec_security_policy.policy2.security_policy_id}"
      },
      "sequence" : 0,
      "type" : "website"
    }
  )
}

resource "akamai_appsec_match_target" "api_4124908" {
  config_id = akamai_appsec_configuration.config.config_id
  match_target = jsonencode(
    {
      "apis" : [
        {
          "id" : 767805,
          "name" : "andrew3"
        }
      ],
      "securityPolicy" : {
        "policyId" : "${akamai_appsec_security_policy.andrew.security_policy_id}"
      },
      "sequence" : 0,
      "type" : "api"
    }
  )
}
