resource "akamai_appsec_waf_mode" "policy2" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  mode               = "KRS"
}

// WAF Rule Actions
// Akamai-X debug Pragma header detected and removed
resource "akamai_appsec_rule" "policy2_akamaipragma_deflection_699989" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "699989"
  rule_action        = "alert"
  condition_exception = jsonencode(
    { "exception" : {
      "specificHeaderCookieParamXmlOrJsonNames" : [
        {
          "names" : [
            "value",
            "body%5B0%5D%5Bvalue%5D"
          ],
          "selector" : "ARGS"
        }
      ]
    } },
    { "conditions" : [
      {
        "type" : "requestHeaderMatch",
        "name" : "User-Agent",
        "positiveMatch" : true,
        "value" : "*Slackbot-LinkExpanding*",
        "valueWildcard" : true
      }
    ] }
  )
}

// Request Indicates an automated program explored the site
resource "akamai_appsec_rule" "policy2_akamaibot_detect_3_v4_699996" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "699996"
  rule_action        = "alert"
}

// Session Fixation
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksession_fixation_950000" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950000"
  rule_action        = "alert"
}

// SQL Injection Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_950001" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950001"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (OS Commands 4)
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_950002" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950002"
  rule_action        = "alert"
}

// Session Fixation
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksession_fixation_950003" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950003"
  rule_action        = "alert"
}

// Remote File Access Attempt
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackfile_injection_950005" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950005"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (OS Commands 5)
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_950006" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950006"
  rule_action        = "alert"
}

// SQL Injection Attack (Blind Testing)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_950007" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950007"
  rule_action        = "alert"
}

// Injection of Undocumented ColdFusion Tags
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackcf_injection_950008" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950008"
  rule_action        = "alert"
}

// Session Fixation
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksession_fixation_950009" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950009"
  rule_action        = "alert"
}

// LDAP Injection Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackldap_injection_950010" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950010"
  rule_action        = "alert"
}

// Server-Side Include (SSI) Attack
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_950011" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950011"
  rule_action        = "alert"
}

// UPDF/XSS injection Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_950018" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950018"
  rule_action        = "alert"
}

// Email Injection Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackemail_injection_950019" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950019"
  rule_action        = "alert"
}

// Path Traversal Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackdir_traversal_950103" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950103"
  rule_action        = "alert"
}

// URL Encoding Abuse Attack Attempt
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationevasion_950107" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950107"
  rule_action        = "alert"
}

// URL Encoding Abuse Attack Attempt
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationevasion_950108" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950108"
  rule_action        = "alert"
}

// Multiple URL Encoding Detected
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationevasion_950109" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950109"
  rule_action        = "alert"
}

// Backdoor access
resource "akamai_appsec_rule" "policy2_owasp_crsmalicious_softwaretrojan_950110" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950110"
  rule_action        = "alert"
}

// Unicode Full/Half Width Abuse Attack Attempt
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationevasion_950116" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950116"
  rule_action        = "alert"
}

// Remote File Inclusion Attack (Remote URL with IP address)
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackrfi_950117" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950117"
  rule_action        = "alert"
}

// Remote File Inclusion Attack (Common PHP RFI Attacks)
resource "akamai_appsec_rule" "policy2_aseweb_attackrfi_950118" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950118"
  rule_action        = "alert"
}

// Remote File Inclusion Attack (Remote URL Ending with '?')
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackrfi_950119" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950119"
  rule_action        = "alert"
}

// Remote File Inclusion Attack (Remote URL Detected)
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackrfi_950120" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950120"
  rule_action        = "alert"
}

// SQL Injection Attack (Tautology Probes 1)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_950901" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950901"
  rule_action        = "alert"
}

// SQL Injection Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_950908" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950908"
  rule_action        = "alert"
}

// HTTP Response Splitting Attack (Header Injection)
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackhttp_response_splitting_950910" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950910"
  rule_action        = "alert"
}

// HTTP Response Splitting Attack (Response Injection)
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackhttp_response_splitting_950911" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950911"
  rule_action        = "alert"
}

// Backdoor access
resource "akamai_appsec_rule" "policy2_owasp_crsmalicious_softwaretrojan_950921" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "950921"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958000" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958000"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958001" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958001"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958002" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958002"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Fromcharcode Detected)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_958003" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958003"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958004" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958004"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958005" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958005"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958006" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958006"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958007" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958007"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (HTML INPUT IMAGE Tag)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_958008" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958008"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958009" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958009"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958010" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958010"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958011" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958011"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958012" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958012"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958013" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958013"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958016" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958016"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958017" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958017"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958018" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958018"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958019" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958019"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958020" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958020"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958022" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958022"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Javascript URL Protocol Handler with "lowsrc" Attribute)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_958023" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958023"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958024" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958024"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958025" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958025"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958026" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958026"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958027" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958027"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958028" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958028"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958030" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958030"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958031" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958031"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958032" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958032"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958033" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958033"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Style Attribute with 'expression' Keyword)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_958034" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958034"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958036" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958036"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958037" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958037"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958038" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958038"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958039" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958039"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958040" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958040"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958041" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958041"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958045" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958045"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958046" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958046"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958047" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958047"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958049" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958049"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Script Tag)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_958051" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958051"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Common PoC DOM Event Triggers)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_958052" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958052"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958054" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958054"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958056" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958056"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958057" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958057"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958059" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958059"
  rule_action        = "alert"
}

// Range: Invalid Last Byte Value
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_hreq_958230" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958230"
  rule_action        = "alert"
}

// Range: Too Many Fields
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_hreq_958231" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958231"
  rule_action        = "alert"
}

// Range: Field Exists and Begins With 0
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_hreq_958291" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958291"
  rule_action        = "alert"
}

// Multiple/Conflicting Connection Header Data Found
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_hreq_958295" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958295"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958404" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958404"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958405" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958405"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958406" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958406"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958407" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958407"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958408" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958408"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958409" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958409"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958410" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958410"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958411" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958411"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958412" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958412"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958413" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958413"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958414" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958414"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958415" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958415"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958416" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958416"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958417" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958417"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958418" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958418"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958419" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958419"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958420" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958420"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958421" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958421"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958422" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958422"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_958423" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958423"
  rule_action        = "alert"
}

// PHP Injection Attack (Common Functions)
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackphp_injection_958976" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958976"
  rule_action        = "alert"
}

// PHP Injection Attack (Configuration Override)
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackphp_injection_958977" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "958977"
  rule_action        = "alert"
}

// SQL Injection Attack (Merge, Execute, Having Probes)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_959070" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "959070"
  rule_action        = "alert"
}

// SQL Injection Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_959071" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "959071"
  rule_action        = "alert"
}

// SQL Injection Attack
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_959072" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "959072"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 1)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_959073" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "959073"
  rule_action        = "alert"
}

// PHP Injection Attack (Opening Tag)
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackphp_injection_959151" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "959151"
  rule_action        = "alert"
}

// Request content type is not allowed by policy
resource "akamai_appsec_rule" "policy2_owasp_crspolicyencoding_not_allowed_960010" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960010"
  rule_action        = "alert"
}

// GET or HEAD Request with Body Content
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_hreq_960011" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960011"
  rule_action        = "alert"
}

// POST Request Missing Content-Length Header
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_hreq_960012" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960012"
  rule_action        = "alert"
}

// Content-Length HTTP Header is Not Numeric
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_hreq_960016" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960016"
  rule_action        = "alert"
}

// Expect Header Not Allowed For HTTP 1.0
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_hreq_960022" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960022"
  rule_action        = "alert"
}

// HTTP Protocol Version is Not Allowed By Policy
resource "akamai_appsec_rule" "policy2_owasp_crspolicyprotocol_not_allowed_960034" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960034"
  rule_action        = "alert"
}

// URL File Extension is Restricted By Policy
resource "akamai_appsec_rule" "policy2_owasp_crspolicyext_restricted_960035" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960035"
  rule_action        = "alert"
}

// Argument value too long
resource "akamai_appsec_rule" "policy2_owasp_crspolicysize_limit_960208" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960208"
  rule_action        = "alert"
}

// Argument name too long
resource "akamai_appsec_rule" "policy2_owasp_crspolicysize_limit_960209" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960209"
  rule_action        = "alert"
}

// Too many arguments in request
resource "akamai_appsec_rule" "policy2_owasp_crspolicysize_limit_960335" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960335"
  rule_action        = "alert"
}

// Total arguments size exceeded
resource "akamai_appsec_rule" "policy2_owasp_crspolicysize_limit_960341" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960341"
  rule_action        = "alert"
}

// Invalid character in request
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationevasion_960901" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960901"
  rule_action        = "alert"
}

// Invalid Use of Identity Encoding
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_hreq_960902" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960902"
  rule_action        = "alert"
}

// Request Containing Content, but Missing Content-Type header
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationmissing_header_960904" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960904"
  rule_action        = "alert"
}

// Failed to Parse Request Body
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_req_960912" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960912"
  rule_action        = "alert"
}

// Multipart Request Body Failed Strict Validation
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_req_960913" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960913"
  rule_action        = "alert"
}

// Multipart Parser Detected a Possible Unmatched Boundary
resource "akamai_appsec_rule" "policy2_owasp_crsprotocol_violationinvalid_req_960914" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "960914"
  rule_action        = "alert"
}

// Possible XSS Attack Detected - HTML Tag Handler
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973300" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973300"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973301" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973301"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973302" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973302"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973303" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973303"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973304" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973304"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (URL Protocols)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_973305" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973305"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973306" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973306"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Eval/Atob Functions)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_973307" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973307"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973308" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973308"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973309" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973309"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973310" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973310"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (XSS Unicode PoC String)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_973311" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973311"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Common PoC Payload)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_973312" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973312"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973313" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973313"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973314" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973314"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973315" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973315"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973316" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973316"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973317" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973317"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973318" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973318"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973319" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973319"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973320" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973320"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973321" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973321"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973322" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973322"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973323" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973323"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973324" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973324"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973325" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973325"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973326" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973326"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973327" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973327"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973328" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973328"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973329" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973329"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973330" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973330"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973331" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973331"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973332" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973332"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973333" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973333"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973334" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973334"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (IE XSS Filter Evasion Attempt)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_973335" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973335"
  rule_action        = "alert"
}

// XSS Filter - Category 1: Script Tag Vector
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973336" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973336"
  rule_action        = "alert"
}

// XSS Filter - Category 2: Event Handler Vector
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_973337" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "973337"
  rule_action        = "alert"
}

// Restricted SQL Character Anomaly Detection Alert - Total # of special characters exceeded
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackspecial_chars_981173" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981173"
  rule_action        = "alert"
}

// Conditional SQL Injection Attempts
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_981241" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981241"
  rule_action        = "alert"
}

// SQL Injection Attack (SQL Operator and Expression Probes 1)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981242" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981242"
  rule_action        = "alert"
}

// SQL Injection Attack (SQL Operator and Expression Probes 2)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981243" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981243"
  rule_action        = "alert"
}

// SQL Injection Attack (Tautology Probes 2)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981244" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981244"
  rule_action        = "alert"
}

// Basic SQL Authentication Bypass Attempts 2/3
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_981245" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981245"
  rule_action        = "alert"
}

// Basic SQL Authentication Bypass Attempts 3/3
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_981246" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981246"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 3)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981247" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981247"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 2)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981248" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981248"
  rule_action        = "alert"
}

// Chained SQL Injection Attempts 2/2
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_981249" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981249"
  rule_action        = "alert"
}

// SQL Benchmark And sleep() Injection Attempts Including Conditional Queries
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_981250" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981250"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 3)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981251" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981251"
  rule_action        = "alert"
}

// SQL Injection Attack (Charset manipulation)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981252" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981252"
  rule_action        = "alert"
}

// SQL Injection Attack (Stored Procedure Detected)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981253" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981253"
  rule_action        = "alert"
}

// SQL Injection Attack (Time-based Blind Probe)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981254" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981254"
  rule_action        = "alert"
}

// SQL Injection Attack (Sysadmin access functions)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981255" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981255"
  rule_action        = "alert"
}

// SQL Injection Attack (Merge, Execute, Match Probes)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981256" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981256"
  rule_action        = "alert"
}

// SQL Injection Attack (Hex Encoding Detected)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981260" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981260"
  rule_action        = "alert"
}

// SQL Injection Attack (NoSQL MongoDB Probes)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981270" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981270"
  rule_action        = "alert"
}

// Blind SQLi Tests Using sleep() or benchmark()
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_981272" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981272"
  rule_action        = "alert"
}

// SQL Injection Attack (UNION Attempt)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981276" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981276"
  rule_action        = "alert"
}

// Integer Overflow Attacks (Taken From Skipfish)
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_981277" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981277"
  rule_action        = "alert"
}

// SQL Injection Attack (SELECT Statement Anomaly Detected)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981300" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981300"
  rule_action        = "alert"
}

// SQL Injection Attack: Common Injection Testing Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_981318" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981318"
  rule_action        = "alert"
}

// SQL Injection Attack: SQL Operator Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_981319" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981319"
  rule_action        = "alert"
}

// SQL Injection Attack (Known/Default DB Resources Probe)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_981320" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "981320"
  rule_action        = "alert"
}

// Request Indicates a Security Scanner Scanned the Site
resource "akamai_appsec_rule" "policy2_owasp_crsautomationsecurity_scanner_990002" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "990002"
  rule_action        = "alert"
}

// Rogue Web Site Crawler
resource "akamai_appsec_rule" "policy2_owasp_crsautomationmalicious_990012" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "990012"
  rule_action        = "alert"
}

// Request Indicates a Security Scanner Scanned the Site
resource "akamai_appsec_rule" "policy2_owasp_crsautomationsecurity_scanner_990901" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "990901"
  rule_action        = "alert"
}

// Request Indicates a Security Scanner Scanned the Site
resource "akamai_appsec_rule" "policy2_owasp_crsautomationsecurity_scanner_990902" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "990902"
  rule_action        = "alert"
}

// SQL Injection Attack (GROUP BY/ORDER BY)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_3000000" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000000"
  rule_action        = "alert"
}

// HTTP Response Splitting (Header Injection Attempt)
resource "akamai_appsec_rule" "policy2_akamaiweb_attackhttp_response_splitting_3000001" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000001"
  rule_action        = "alert"
}

// Local System File Access Attempt
resource "akamai_appsec_rule" "policy2_akamaiweb_attackfile_injection_3000002" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000002"
  rule_action        = "alert"
}

// PHP Code Injection
resource "akamai_appsec_rule" "policy2_akamaiweb_attackphp_injection_3000003" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000003"
  rule_action        = "alert"
}

// Potential Remote File Inclusion (RFI) Attack: Suspicious Off-Domain URL Reference
resource "akamai_appsec_rule" "policy2_aseweb_attackrfi_3000004" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000004"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (OS commands with full path)
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000005" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000005"
  rule_action        = "alert"
}

// SQL Injection Attack (Comment String Termination)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_3000006" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000006"
  rule_action        = "alert"
}

// Command Injection (Unix File Leakage)
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000007" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000007"
  rule_action        = "alert"
}

// Pandora / DirtJumper DDoS Detection - HTTP GET Attacks
resource "akamai_appsec_rule" "policy2_akamaiautomationmalicious_3000008" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000008"
  rule_action        = "alert"
}

// Ruby on Rails YAML Injection Attack
resource "akamai_appsec_rule" "policy2_akamaiweb_attackruby_injection_3000009" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000009"
  rule_action        = "alert"
}

// LOIC 1.1 DoS Detection
resource "akamai_appsec_rule" "policy2_akamaiweb_attackloic11_3000010" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000010"
  rule_action        = "alert"
}

// HULK DoS Attack Tool Detection
resource "akamai_appsec_rule" "policy2_akamaiweb_attackhulk_3000011" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000011"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (OGNL Injection)
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000012" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000012"
  rule_action        = "alert"
}

// System Command Injection (Attacker Toolset Download)
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000013" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000013"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (OGNL Injection)
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000014" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000014"
  rule_action        = "alert"
}

// SQL Injection Attack (Database Timing Query)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_3000015" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000015"
  rule_action        = "alert"
}

// PHP Code Injection Using Data Stream Wrapper
resource "akamai_appsec_rule" "policy2_akamaiweb_attackphp_injection_3000016" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000016"
  rule_action        = "alert"
}

// MySQL Keywords Anomaly Detection Alert
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_3000017" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000017"
  rule_action        = "alert"
}

// DirtJumper DDoS Detection - HTTP POST Attacks
resource "akamai_appsec_rule" "policy2_akamaiautomationmalicious_3000018" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000018"
  rule_action        = "alert"
}

// Pandora DDoS Detection - HTTP POST Attacks
resource "akamai_appsec_rule" "policy2_akamaiautomationmalicious_3000019" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000019"
  rule_action        = "alert"
}

// Local File Inclusion (and Command Injection) Using '/proc/self/environ'
resource "akamai_appsec_rule" "policy2_akamaiweb_attacklfi_3000020" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000020"
  rule_action        = "alert"
}

// Detect Attempts to Access the Wordpress Pingback API
resource "akamai_appsec_rule" "policy2_akamaiweb_attackwordpress_pingback_3000021" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000021"
  rule_action        = "alert"
}

// SQL Injection (Built-in Functions, Objects and Keyword Probes 4)
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_3000022" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000022"
  rule_action        = "alert"
}

// Apache Struts ClassLoader Manipulation Remote Code Execution
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000023" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000023"
  rule_action        = "alert"
}

// Apache Commons FileUpload and Apache Tomcat DoS
resource "akamai_appsec_rule" "policy2_akamaiweb_attackapache_commons_dos_3000024" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000024"
  rule_action        = "alert"
}

// CVE-2014-6271 Bash Command Injection Attack
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000025" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000025"
  rule_action        = "alert"
}

// XXE External Entity
resource "akamai_appsec_rule" "policy2_akamaiweb_attackxxe_3000027" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000027"
  rule_action        = "alert"
}

// SQL Injection Attack: MySQL comments, conditions and ch(a)r injections
resource "akamai_appsec_rule" "policy2_aseweb_attacksql_injection_3000029" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000029"
  rule_action        = "alert"
}

// Basic SQL Authentication Bypass Attempts 3/3
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attacksql_injection_3000030" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000030"
  rule_action        = "alert"
}

// HTTP.sys Remote Code Execution Vulnerability Attack Detected (CVE-2015-1635)
resource "akamai_appsec_rule" "policy2_akamaiweb_attackiis_range_3000031" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000031"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack Event Handler
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_3000032" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000032"
  rule_action        = "alert"
}

// PHP Wrapper Attack
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000033" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000033"
  rule_action        = "alert"
}

// Command Injection via the Java Runtime.getRuntime() Method
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000034" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000034"
  rule_action        = "alert"
}

// Potential Account Brute Force Guessing via Wordpress XML-RPC API authenticated methods
resource "akamai_appsec_rule" "policy2_akamaiweb_attackwordpress_bruteforce_3000035" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000035"
  rule_action        = "alert"
}

// Detected LOIC / HOIC client request based on query string
resource "akamai_appsec_rule" "policy2_akamaiddosloic_hoic_1_v1_3000036" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000036"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (JS On-Event Handler)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_3000037" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000037"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (DOM Window Properties)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_3000038" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000038"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (DOM Document Methods)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_3000039" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000039"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Document Methods
resource "akamai_appsec_rule" "policy2_akamaiweb_attackxss_3000040" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000040"
  rule_action        = "alert"
}

// Server Side Template Injection (SSTI)
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000041" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000041"
  rule_action        = "alert"
}

// Detected ARDT client request
resource "akamai_appsec_rule" "policy2_akamaiddosardt_3000042" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000042"
  rule_action        = "alert"
}

// Detect Attempts to Access the Wordpress system.multicall XML-RPC API
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksystem_multicall_3000043" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000043"
  rule_action        = "alert"
}

// SQL Injection Using Boolean Logic
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksql_injection_3000044" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000044"
  rule_action        = "alert"
}

// SQL Injection Using Boolean Logic 1
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksql_injection_3000045" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000045"
  rule_action        = "alert"
}

// SQL Injection Using Boolean Logic 2
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksql_injection_3000046" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000046"
  rule_action        = "alert"
}

// SQL Injection Using SQL Backup Command
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksql_injection_3000047" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000047"
  rule_action        = "alert"
}

// SQL Injection Using SQL Restore Command
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksql_injection_3000048" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000048"
  rule_action        = "alert"
}

// SQL Injection With Cursor Declaration
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksql_injection_3000049" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000049"
  rule_action        = "alert"
}

// SQL Injection Using Boolean Logic 3
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksql_injection_3000050" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000050"
  rule_action        = "alert"
}

// SQL Injection Using Boolean Logic 4
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksql_injection_3000051" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000051"
  rule_action        = "alert"
}

// SQL Injection Using EXISTS
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksql_injection_3000052" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000052"
  rule_action        = "alert"
}

// SQL Injection Using DELETE Statements
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksql_injection_3000053" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000053"
  rule_action        = "alert"
}

// SQL Injection Using UPDATE Statements
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksql_injection_3000054" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000054"
  rule_action        = "alert"
}

// Avzhan Bot DDOS Detection
resource "akamai_appsec_rule" "policy2_akamaiautomationmalicious_3000055" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000055"
  rule_action        = "alert"
}

// PHP Object Injection Attack Detected
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000056" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000056"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Common Attack Tool Keywords)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_3000057" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000057"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (OGNL Injection)
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000058" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000058"
  rule_action        = "alert"
}

// Request Headers indicate request came from Wordpress Pingback
resource "akamai_appsec_rule" "policy2_akamaiweb_attackwordpress_pingback_3000059" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000059"
  rule_action        = "alert"
}

// Mirai / Kaiten DDoS Detection - HTTP Attacks
resource "akamai_appsec_rule" "policy2_akamaiautomationmalicious_3000060" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000060"
  rule_action        = "alert"
}

// Cross-site Scripting Attack (Referer Header From OpenBugBounty Website)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_3000061" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000061"
  rule_action        = "alert"
}

// Mirai/Kaiten Bot DDOS Detection - Bogus Search Engine Referer
resource "akamai_appsec_rule" "policy2_akamaiautomationmalicious_3000062" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000062"
  rule_action        = "alert"
}

// Wordpress wp-json Attack Attempt - non-integer character(s) in ID parameter paylaod
resource "akamai_appsec_rule" "policy2_akamaiweb_attackwordpress_3000063" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000063"
  rule_action        = "alert"
}

// Application Layer Hash DoS Attack
resource "akamai_appsec_rule" "policy2_akamaiautomationmalicious_3000064" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000064"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (Deserialization Attack)
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000065" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000065"
  rule_action        = "alert"
}

// Potential WordPress Parameter Resource Consumption Remote DoS Attack (CVE-2018-6389)
resource "akamai_appsec_rule" "policy2_akamaiweb_attackwordpress_dos_3000066" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000066"
  rule_action        = "alert"
}

// Potential Drupal Attack (CVE-2018-7600)
resource "akamai_appsec_rule" "policy2_akamaiweb_attackdrupal_3000067" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000067"
  rule_action        = "alert"
}

// ESI injection Attack
resource "akamai_appsec_rule" "policy2_akamaiweb_attackesi_injection_3000068" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000068"
  rule_action        = "alert"
}

// Webshell/Backdoor File Upload Attempt
resource "akamai_appsec_rule" "policy2_owasp_crsmalicious_softwaretrojan_3000071" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000071"
  rule_action        = "alert"
}

// Deserialization Attack Detected
resource "akamai_appsec_rule" "policy2_aseweb_attackcmd_injection_3000072" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000072"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Invalid Transfer-Encoding Header Value
resource "akamai_appsec_rule" "policy2_akamaiprotocol_violationhttp_desync_3000073" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000073"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: HTTP Request Smuggling Detect in Request Body
resource "akamai_appsec_rule" "policy2_akamaiprotocol_violationhttp_desync_3000074" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000074"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Transfer-Encoding Header Name Obfuscation
resource "akamai_appsec_rule" "policy2_akamaiprotocol_violationhttp_desync_3000075" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000075"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Transfer-Encoding Header in Request Body
resource "akamai_appsec_rule" "policy2_akamaiprotocol_violationhttp_desync_3000076" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000076"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Chunked header value with invalid Header Name
resource "akamai_appsec_rule" "policy2_akamaiprotocol_violationhttp_desync_3000077" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000077"
  rule_action        = "alert"
}

// Microsoft Sharepoint Remote Command Execution (Deserialization Attack)
resource "akamai_appsec_rule" "policy2_akamaiweb_attacksharepoint_deserial_3000079" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000079"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Attribute Injection 1)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_3000080" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000080"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Attribute Injection 2)
resource "akamai_appsec_rule" "policy2_aseweb_attackxss_3000081" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000081"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackxss_3000082" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000082"
  rule_action        = "alert"
}

// Possible MS Exchange/OWA Attack Detected (CVE-2021-26855)
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackplatform_3000083" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000083"
  rule_action        = "alert"
}

// Possible MS Exchange/OWA Attack Detected (CVE-2021-27065)
resource "akamai_appsec_rule" "policy2_owasp_crsweb_attackplatform_3000084" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000084"
  rule_action        = "alert"
}

// PROXY Header Detected
resource "akamai_appsec_rule" "policy2_akamaiweb_attackproxy_header_detected_3000999" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy2.security_policy_id
  rule_id            = "3000999"
  rule_action        = "alert"
}


// WAF Attack Group Actions
resource "akamai_appsec_attack_group" "policy2_SQL" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy2.security_policy_id
  attack_group        = "SQL"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy2_XSS" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy2.security_policy_id
  attack_group        = "XSS"
  attack_group_action = "deny"
}

resource "akamai_appsec_attack_group" "policy2_CMD" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy2.security_policy_id
  attack_group        = "CMD"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy2_HTTP" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy2.security_policy_id
  attack_group        = "HTTP"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy2_RFI" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy2.security_policy_id
  attack_group        = "RFI"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy2_PHP" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy2.security_policy_id
  attack_group        = "PHP"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy2_TROJAN" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy2.security_policy_id
  attack_group        = "TROJAN"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy2_DDOS" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy2.security_policy_id
  attack_group        = "DDOS"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy2_IN" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy2.security_policy_id
  attack_group        = "IN"
  attack_group_action = "deny_custom_78842"
}

resource "akamai_appsec_attack_group" "policy2_OUT" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy2.security_policy_id
  attack_group        = "OUT"
  attack_group_action = "alert"
}

resource "akamai_appsec_waf_mode" "andrew" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  mode               = "KRS"
}

// WAF Rule Actions
// CMD Injection Attack Detected (OS Commands 4)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_950002" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "950002"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (OS Commands 5)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_950006" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "950006"
  rule_action        = "alert"
}

// SQL Injection Attack (Blind Testing)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_950007" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "950007"
  rule_action        = "alert"
}

// Server-Side Include (SSI) Attack
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_950011" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "950011"
  rule_action        = "alert"
}

// Remote File Inclusion Attack (Common PHP RFI Attacks)
resource "akamai_appsec_rule" "andrew_aseweb_attackrfi_950118" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "950118"
  rule_action        = "alert"
}

// Local File Inclusion (LFI) Attack (Directory Traversal and Obfuscation Attempts)
resource "akamai_appsec_rule" "andrew_aseweb_attacklfi_950203" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "950203"
  rule_action        = "alert"
}

// Unicode Full/Half Width Abuse Attack Attempt
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_950216" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "950216"
  rule_action        = "alert"
}

// Possible URL Redirector Abuse (Off-Domain URL)
resource "akamai_appsec_rule" "andrew_aseweb_attackpolicy_950220" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "950220"
  rule_action        = "alert"
}

// SQL Injection Attack (Tautology Probes 1)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_950901" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "950901"
  rule_action        = "alert"
}

// HTTP Response Splitting Attack (Header Injection)
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_951910" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "951910"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Fromcharcode Detected)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_958003" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "958003"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (HTML INPUT IMAGE Tag)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_958008" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "958008"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Javascript URL Protocol Handler with "lowsrc" Attribute)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_958023" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "958023"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Style Attribute with 'expression' Keyword)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_958034" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "958034"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Script Tag)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_958051" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "958051"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Common PoC DOM Event Triggers)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_958052" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "958052"
  rule_action        = "alert"
}

// SQL Injection Attack (Merge, Execute, Having Probes)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_959070" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "959070"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 1)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_959073" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "959073"
  rule_action        = "alert"
}

// PHP Injection Attack (Common Functions)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_959976" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "959976"
  rule_action        = "alert"
}

// PHP Injection Attack (Configuration Override)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_959977" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "959977"
  rule_action        = "alert"
}

// GET or HEAD Request with Body Content
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_961011" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "961011"
  rule_action        = "alert"
}

// POST Request Missing Content-Length Header
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_961012" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "961012"
  rule_action        = "alert"
}

// Invalid HTTP Protocol Version
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_961034" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "961034"
  rule_action        = "alert"
}

// Request Containing Content, but Missing Content-Type header
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_961904" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "961904"
  rule_action        = "alert"
}

// Failed to Parse Request Body for WAF Inspection
resource "akamai_appsec_rule" "andrew_aseweb_attackpolicy_961912" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "961912"
  rule_action        = "alert"
}

// HTTP Range Header: Invalid Last Byte Value
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_968230" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "968230"
  rule_action        = "alert"
}

// PHP Injection Attack (Opening Tag)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_969151" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "969151"
  rule_action        = "alert"
}

// SQL Information Leakage
resource "akamai_appsec_rule" "andrew_aseoutboundleakage_970003" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970003"
  rule_action        = "alert"
}

// IIS Information Leakage
resource "akamai_appsec_rule" "andrew_aseoutboundleakage_970004" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970004"
  rule_action        = "alert"
}

// PHP Information Leakage
resource "akamai_appsec_rule" "andrew_aseoutboundleakage_970009" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970009"
  rule_action        = "alert"
}

// File or Directory Names Leakage
resource "akamai_appsec_rule" "andrew_aseoutboundleakage_970011" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970011"
  rule_action        = "alert"
}

// Directory Listing
resource "akamai_appsec_rule" "andrew_aseoutboundleakage_970013" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970013"
  rule_action        = "alert"
}

// ASP/JSP Source Code Leakage
resource "akamai_appsec_rule" "andrew_aseoutboundleakage_970014" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970014"
  rule_action        = "alert"
}

// PHP Source Code Leakage
resource "akamai_appsec_rule" "andrew_aseoutboundleakage_970015" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970015"
  rule_action        = "alert"
}

// Application is not Available (Server-Side Exceptions)
resource "akamai_appsec_rule" "andrew_aseoutbounderror_970118" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970118"
  rule_action        = "alert"
}

// Application is not Available (HTTP 5XX)
resource "akamai_appsec_rule" "andrew_aseoutbounderror_970901" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970901"
  rule_action        = "alert"
}

// PHP Source Code Leakage
resource "akamai_appsec_rule" "andrew_aseoutboundleakage_970902" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970902"
  rule_action        = "alert"
}

// ASP/JSP Source Code Leakage
resource "akamai_appsec_rule" "andrew_aseoutboundleakage_970903" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970903"
  rule_action        = "alert"
}

// IIS Information Leakage
resource "akamai_appsec_rule" "andrew_aseoutboundleakage_970904" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "970904"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (URL Protocols)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_973305" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "973305"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Eval/Atob Functions)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_973307" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "973307"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (XSS Unicode PoC String)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_973311" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "973311"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Common PoC Payload)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_973312" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "973312"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (IE XSS Filter Evasion Attempt)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_973335" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "973335"
  rule_action        = "alert"
}

// SQL Injection Attack (SQL Conditional Probes)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981240" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981240"
  rule_action        = "alert"
}

// SQL Injection Attack (SQL Operator and Expression Probes 1)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981242" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981242"
  rule_action        = "alert"
}

// SQL Injection Attack (SQL Operator and Expression Probes 2)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981243" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981243"
  rule_action        = "alert"
}

// SQL Injection Attack (Tautology Probes 2)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981244" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981244"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 3)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981247" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981247"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 2)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981248" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981248"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 3)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981251" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981251"
  rule_action        = "alert"
}

// SQL Injection Attack (Charset manipulation)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981252" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981252"
  rule_action        = "alert"
}

// SQL Injection Attack (Stored Procedure Detected)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981253" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981253"
  rule_action        = "alert"
}

// SQL Injection Attack (Time-based Blind Probe)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981254" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981254"
  rule_action        = "alert"
}

// SQL Injection Attack (Sysadmin access functions)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981255" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981255"
  rule_action        = "alert"
}

// SQL Injection Attack (Merge, Execute, Match Probes)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981256" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981256"
  rule_action        = "alert"
}

// SQL Injection Attack (Hex Encoding Detected)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981260" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981260"
  rule_action        = "alert"
}

// SQL Injection Attack (NoSQL MongoDB Probes)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981270" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981270"
  rule_action        = "alert"
}

// SQL Injection Attack (UNION Attempt)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981276" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981276"
  rule_action        = "alert"
}

// SQL Injection Attack (SELECT Statement Anomaly Detected)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981300" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981300"
  rule_action        = "alert"
}

// SQL Injection Attack (Known/Default DB Resources Probe)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_981320" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "981320"
  rule_action        = "alert"
}

// Security Scanner/Web Attack Tool Detected (User-Agent)
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_999002" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "999002"
  rule_action        = "alert"
}

// Security Scanner/Web Attack Tool Detected (Request Header Names)
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_999901" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "999901"
  rule_action        = "alert"
}

// Security Scanner/Web Attack Tool Detected (Filename)
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_999902" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "999902"
  rule_action        = "alert"
}

// SQL Injection Attack (GROUP BY/ORDER BY)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_3000000" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000000"
  rule_action        = "alert"
}

// Potential Remote File Inclusion (RFI) Attack: Suspicious Off-Domain URL Reference
resource "akamai_appsec_rule" "andrew_aseweb_attackrfi_3000004" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000004"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (OS commands with full path)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000005" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000005"
  rule_action        = "alert"
}

// SQL Injection Attack (Comment String Termination)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_3000006" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000006"
  rule_action        = "alert"
}

// Command Injection (Unix File Leakage)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000007" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000007"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (OGNL Injection)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000012" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000012"
  rule_action        = "alert"
}

// System Command Injection (Attacker Toolset Download)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000013" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000013"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (OGNL Injection)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000014" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000014"
  rule_action        = "alert"
}

// SQL Injection Attack (Database Timing Query)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_3000015" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000015"
  rule_action        = "alert"
}

// MySQL Keywords Anomaly Detection Alert
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_3000017" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000017"
  rule_action        = "alert"
}

// SQL Injection (Built-in Functions, Objects and Keyword Probes 4)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_3000022" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000022"
  rule_action        = "alert"
}

// Apache Struts ClassLoader Manipulation Remote Code Execution
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000023" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000023"
  rule_action        = "alert"
}

// CVE-2014-6271 Bash Command Injection Attack
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000025" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000025"
  rule_action        = "alert"
}

// SQL Injection Attack: MySQL comments, conditions and ch(a)r injections
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_3000029" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000029"
  rule_action        = "alert"
}

// PHP Wrapper Attack
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000033" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000033"
  rule_action        = "alert"
}

// Command Injection via the Java Runtime.getRuntime() Method
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000034" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000034"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (JS On-Event Handler)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000037" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000037"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (DOM Window Properties)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000038" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000038"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (DOM Document Methods)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000039" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000039"
  rule_action        = "alert"
}

// Server Side Template Injection (SSTI)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000041" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000041"
  rule_action        = "alert"
}

// PHP Object Injection Attack Detected
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000056" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000056"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Common Attack Tool Keywords)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000057" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000057"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (OGNL Injection)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000058" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000058"
  rule_action        = "alert"
}

// Cross-site Scripting Attack (Referer Header From OpenBugBounty Website)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000061" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000061"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (Deserialization Attack)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000065" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000065"
  rule_action        = "alert"
}

// Deserialization Attack Detected
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000072" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000072"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Attribute Injection 1)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000080" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000080"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Attribute Injection 2)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000081" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000081"
  rule_action        = "alert"
}

// SQL Injection Attack (SmartDetect)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_3000100" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000100"
  rule_action        = "alert"
}

// SQL Injection Attack (Common SQL Database Probes)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_3000101" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000101"
  rule_action        = "alert"
}

// SQL Injection Attack (Null Byte Detected)
resource "akamai_appsec_rule" "andrew_aseweb_attacksql_injection_3000102" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000102"
  rule_action        = "alert"
}

// Pandora / DirtJumper DDoS Detection - HTTP GET Attacks
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000108" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000108"
  rule_action        = "alert"
}

// Ruby on Rails YAML Injection Attack
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000109" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000109"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (SmartDetect)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000110" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000110"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Common PoC Probes 1)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000111" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000111"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Common PoC Probes 2)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000112" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000112"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Javascript Mixed Case Obfuscation)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000113" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000113"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Shell Script Execution)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000114" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000114"
  rule_action        = "alert"
}

// LOIC 1.1 DoS Detection
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000115" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000115"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (HTML Injection)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000116" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000116"
  rule_action        = "alert"
}

// HULK DoS Attack Tool Detected
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000117" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000117"
  rule_action        = "alert"
}

// DirtJumper DDoS Detection - HTTP POST Attacks
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000118" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000118"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (HTML Context Breaking)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000119" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000119"
  rule_action        = "alert"
}

// Local File Inclusion (LFI) Attack (Common OS Files 1)
resource "akamai_appsec_rule" "andrew_aseweb_attacklfi_3000120" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000120"
  rule_action        = "alert"
}

// Local File Inclusion (LFI) Attack (Common OS Files 2)
resource "akamai_appsec_rule" "andrew_aseweb_attacklfi_3000121" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000121"
  rule_action        = "alert"
}

// Local File Inclusion (LFI) Attack (Long Directory Traversal)
resource "akamai_appsec_rule" "andrew_aseweb_attacklfi_3000122" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000122"
  rule_action        = "alert"
}

// Local File Inclusion (LFI) Attack (Directory Traversal Obfuscation)
resource "akamai_appsec_rule" "andrew_aseweb_attacklfi_3000123" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000123"
  rule_action        = "alert"
}

// Local File Inclusion (LFI) Attack (Common OS Files 3)
resource "akamai_appsec_rule" "andrew_aseweb_attacklfi_3000124" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000124"
  rule_action        = "alert"
}

// Local File Inclusion (LFI) Attack (Common OS Files 4)
resource "akamai_appsec_rule" "andrew_aseweb_attacklfi_3000125" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000125"
  rule_action        = "alert"
}

// Local File Inclusion (LFI) Attack (Common OS Files 5)
resource "akamai_appsec_rule" "andrew_aseweb_attacklfi_3000126" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000126"
  rule_action        = "alert"
}

// Local File Inclusion (LFI) Attack (Nul Byte Detected)
resource "akamai_appsec_rule" "andrew_aseweb_attacklfi_3000127" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000127"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (HTML Entity Named Encoding Detected)
resource "akamai_appsec_rule" "andrew_aseweb_attackxss_3000128" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000128"
  rule_action        = "alert"
}

// Pandora DDoS Detection - HTTP POST Attacks
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000129" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000129"
  rule_action        = "alert"
}

// Remote File Inclusion Attack (Well-Known RFI Testing/Attack URL)
resource "akamai_appsec_rule" "andrew_aseweb_attackrfi_3000130" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000130"
  rule_action        = "alert"
}

// Remote File Inclusion Attack (Well-Known RFI Filename)
resource "akamai_appsec_rule" "andrew_aseweb_attackrfi_3000131" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000131"
  rule_action        = "alert"
}

// Detect Attempts to Access the Wordpress Pingback API
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000132" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000132"
  rule_action        = "alert"
}

// Apache Commons FileUpload and Apache Tomcat DoS
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000133" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000133"
  rule_action        = "alert"
}

// XML External Entity (XXE) Attack
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000134" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000134"
  rule_action        = "alert"
}

// HTTP.sys Remote Code Execution Vulnerability Attack Detected (CVE-2015-1635)
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000135" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000135"
  rule_action        = "alert"
}

// Potential Account Brute Force Guessing via Wordpress XML-RPC API authenticated methods
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000136" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000136"
  rule_action        = "alert"
}

// Detected LOIC / HOIC client request based on query string
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000137" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000137"
  rule_action        = "alert"
}

// Detected ARDT client request
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000138" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000138"
  rule_action        = "alert"
}

// Detect Attempts to Access the Wordpress system.multicall XML-RPC API
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000139" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000139"
  rule_action        = "alert"
}

// Avzhan Bot DDOS Detection
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000140" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000140"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (OS Commands 1)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000141" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000141"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (OS Commands 2)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000142" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000142"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (Bash with -c flag)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000143" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000143"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (Uname with -a flag)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000144" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000144"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (Cmd.exe with "dir" command)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000145" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000145"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (/bin/sh with pipe "|")
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000146" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000146"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (Shellshock Variant)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000147" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000147"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (Ping Beaconing)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000148" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000148"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (Common Uname PoC)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000149" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000149"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (Sleep with Bracketed IFS Obfuscation)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000150" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000150"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (Bracketed IFS Argument Separator Detected)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000151" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000151"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (IP Address Detected)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000152" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000152"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (OS Commands 3)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000153" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000153"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (Common PHP Function Detected)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000154" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000154"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (Php/Data Filter Detected)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000155" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000155"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (PHP High-Risk Functions)
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000156" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000156"
  rule_action        = "alert"
}

// Mirai / Kaiten DDoS Detection - HTTP Attacks
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000157" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000157"
  rule_action        = "alert"
}

// Security Scanner/Web Attack Tool Detected (PoC Testing Payload)
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000160" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000160"
  rule_action        = "alert"
}

// Mirai/Kaiten Bot DDOS Detection - Bogus Search Engine Referer
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000162" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000162"
  rule_action        = "alert"
}

// Application Layer Hash DoS Attack
resource "akamai_appsec_rule" "andrew_aseweb_attacktool_3000164" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000164"
  rule_action        = "alert"
}

// Potential Wordpress Javascript DoS Attack (CVE-2018-6389)
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000166" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000166"
  rule_action        = "alert"
}

// Potential Drupal Attack (CVE-2018-7600)
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000167" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000167"
  rule_action        = "alert"
}

// Edge Side Inclusion (ESI) injection Attack
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000168" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000168"
  rule_action        = "alert"
}

// Webshell/Backdoor File Upload Attempt
resource "akamai_appsec_rule" "andrew_aseweb_attackcmd_injection_3000171" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000171"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Invalid Transfer-Encoding Header Value
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_3000173" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000173"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: HTTP Request Smuggling Detect in Request Body
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_3000174" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000174"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Transfer-Encoding Header Name Obfuscation
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_3000175" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000175"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Transfer-Encoding Header in Request Body
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_3000176" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000176"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Chunked header value with invalid Header Name
resource "akamai_appsec_rule" "andrew_aseweb_attackprotocol_3000177" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000177"
  rule_action        = "alert"
}

// Microsoft Sharepoint Remote Command Execution (Deserialization Attack)
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000179" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000179"
  rule_action        = "alert"
}

// Partial Request Body Inspection Warning - Request Body is larger than the configured inspection limit
resource "akamai_appsec_rule" "andrew_aseweb_attackpolicy_3000180" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000180"
  rule_action        = "alert"
}

// Possible MS Exchange/OWA Attack Detected (CVE-2021-26855)
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000183" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000183"
  rule_action        = "alert"
}

// Possible MS Exchange/OWA Attack Detected (CVE-2021-27065)
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000184" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000184"
  rule_action        = "alert"
}

// Confluence/OGNLi Attack Detected (CVE-2021-26084)
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000185" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000185"
  rule_action        = "alert"
}

// PowerCMS Movable Type Attack Detected (CVE-2021-20837)
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000186" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000186"
  rule_action        = "alert"
}

// Magento vulnerability (Callback function) Attack Detected (CVE-2022-24086 CVE-2022-24087)
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000187" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000187"
  rule_action        = "alert"
}

// Magento vulnerability (validate_rules) Attack Detected (CVE-2022-24086 CVE-2022-24087)
resource "akamai_appsec_rule" "andrew_aseweb_attackplatform_3000188" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.andrew.security_policy_id
  rule_id            = "3000188"
  rule_action        = "alert"
}


// WAF Attack Group Actions
resource "akamai_appsec_attack_group" "andrew_POLICY" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.andrew.security_policy_id
  attack_group        = "POLICY"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "andrew_WAT" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.andrew.security_policy_id
  attack_group        = "WAT"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "andrew_PROTOCOL" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.andrew.security_policy_id
  attack_group        = "PROTOCOL"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "andrew_SQL" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.andrew.security_policy_id
  attack_group        = "SQL"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "andrew_XSS" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.andrew.security_policy_id
  attack_group        = "XSS"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "andrew_CMD" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.andrew.security_policy_id
  attack_group        = "CMD"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "andrew_LFI" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.andrew.security_policy_id
  attack_group        = "LFI"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "andrew_RFI" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.andrew.security_policy_id
  attack_group        = "RFI"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "andrew_PLATFORM" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.andrew.security_policy_id
  attack_group        = "PLATFORM"
  attack_group_action = "alert"
}

resource "akamai_appsec_waf_mode" "policy1" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  mode               = "KRS"
}

// WAF Rule Actions
// Akamai-X debug Pragma header detected and removed
resource "akamai_appsec_rule" "policy1_akamaipragma_deflection_699989" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "699989"
  rule_action        = "deny"
}

// Request Indicates an automated program explored the site
resource "akamai_appsec_rule" "policy1_akamaibot_detect_3_v4_699996" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "699996"
  rule_action        = "alert"
}

// Session Fixation
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksession_fixation_950000" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950000"
  rule_action        = "alert"
}

// SQL Injection Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_950001" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950001"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (OS Commands 4)
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_950002" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950002"
  rule_action        = "alert"
}

// Session Fixation
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksession_fixation_950003" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950003"
  rule_action        = "alert"
}

// Remote File Access Attempt
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackfile_injection_950005" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950005"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (OS Commands 5)
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_950006" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950006"
  rule_action        = "alert"
}

// SQL Injection Attack (Blind Testing)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_950007" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950007"
  rule_action        = "alert"
}

// Injection of Undocumented ColdFusion Tags
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackcf_injection_950008" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950008"
  rule_action        = "alert"
}

// Session Fixation
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksession_fixation_950009" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950009"
  rule_action        = "alert"
}

// LDAP Injection Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackldap_injection_950010" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950010"
  rule_action        = "alert"
}

// Server-Side Include (SSI) Attack
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_950011" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950011"
  rule_action        = "alert"
}

// UPDF/XSS injection Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_950018" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950018"
  rule_action        = "alert"
}

// Email Injection Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackemail_injection_950019" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950019"
  rule_action        = "alert"
}

// Path Traversal Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackdir_traversal_950103" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950103"
  rule_action        = "alert"
}

// URL Encoding Abuse Attack Attempt
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationevasion_950107" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950107"
  rule_action        = "alert"
}

// URL Encoding Abuse Attack Attempt
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationevasion_950108" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950108"
  rule_action        = "alert"
}

// Multiple URL Encoding Detected
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationevasion_950109" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950109"
  rule_action        = "alert"
}

// Backdoor access
resource "akamai_appsec_rule" "policy1_owasp_crsmalicious_softwaretrojan_950110" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950110"
  rule_action        = "alert"
}

// Unicode Full/Half Width Abuse Attack Attempt
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationevasion_950116" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950116"
  rule_action        = "alert"
}

// Remote File Inclusion Attack (Remote URL with IP address)
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackrfi_950117" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950117"
  rule_action        = "alert"
}

// Remote File Inclusion Attack (Common PHP RFI Attacks)
resource "akamai_appsec_rule" "policy1_aseweb_attackrfi_950118" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950118"
  rule_action        = "alert"
}

// Remote File Inclusion Attack (Remote URL Ending with '?')
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackrfi_950119" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950119"
  rule_action        = "alert"
}

// Remote File Inclusion Attack (Remote URL Detected)
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackrfi_950120" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950120"
  rule_action        = "alert"
}

// SQL Injection Attack (Tautology Probes 1)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_950901" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950901"
  rule_action        = "alert"
}

// SQL Injection Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_950908" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950908"
  rule_action        = "alert"
}

// HTTP Response Splitting Attack (Header Injection)
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackhttp_response_splitting_950910" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950910"
  rule_action        = "alert"
}

// HTTP Response Splitting Attack (Response Injection)
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackhttp_response_splitting_950911" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950911"
  rule_action        = "alert"
}

// Backdoor access
resource "akamai_appsec_rule" "policy1_owasp_crsmalicious_softwaretrojan_950921" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "950921"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958000" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958000"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958001" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958001"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958002" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958002"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Fromcharcode Detected)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_958003" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958003"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958004" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958004"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958005" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958005"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958006" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958006"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958007" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958007"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (HTML INPUT IMAGE Tag)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_958008" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958008"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958009" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958009"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958010" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958010"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958011" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958011"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958012" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958012"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958013" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958013"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958016" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958016"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958017" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958017"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958018" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958018"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958019" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958019"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958020" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958020"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958022" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958022"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Javascript URL Protocol Handler with "lowsrc" Attribute)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_958023" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958023"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958024" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958024"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958025" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958025"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958026" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958026"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958027" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958027"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958028" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958028"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958030" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958030"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958031" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958031"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958032" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958032"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958033" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958033"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Style Attribute with 'expression' Keyword)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_958034" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958034"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958036" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958036"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958037" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958037"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958038" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958038"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958039" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958039"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958040" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958040"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958041" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958041"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958045" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958045"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958046" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958046"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958047" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958047"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958049" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958049"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Script Tag)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_958051" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958051"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Common PoC DOM Event Triggers)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_958052" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958052"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958054" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958054"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958056" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958056"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958057" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958057"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958059" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958059"
  rule_action        = "alert"
}

// Range: Invalid Last Byte Value
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_hreq_958230" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958230"
  rule_action        = "alert"
}

// Range: Too Many Fields
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_hreq_958231" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958231"
  rule_action        = "alert"
}

// Range: Field Exists and Begins With 0
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_hreq_958291" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958291"
  rule_action        = "alert"
}

// Multiple/Conflicting Connection Header Data Found
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_hreq_958295" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958295"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958404" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958404"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958405" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958405"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958406" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958406"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958407" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958407"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958408" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958408"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958409" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958409"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958410" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958410"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958411" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958411"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958412" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958412"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958413" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958413"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958414" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958414"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958415" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958415"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958416" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958416"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958417" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958417"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958418" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958418"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958419" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958419"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958420" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958420"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958421" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958421"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958422" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958422"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_958423" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958423"
  rule_action        = "alert"
}

// PHP Injection Attack (Common Functions)
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackphp_injection_958976" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958976"
  rule_action        = "alert"
}

// PHP Injection Attack (Configuration Override)
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackphp_injection_958977" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "958977"
  rule_action        = "alert"
}

// SQL Injection Attack (Merge, Execute, Having Probes)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_959070" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "959070"
  rule_action        = "alert"
}

// SQL Injection Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_959071" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "959071"
  rule_action        = "alert"
}

// SQL Injection Attack
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_959072" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "959072"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 1)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_959073" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "959073"
  rule_action        = "alert"
}

// PHP Injection Attack (Opening Tag)
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackphp_injection_959151" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "959151"
  rule_action        = "alert"
}

// Request content type is not allowed by policy
resource "akamai_appsec_rule" "policy1_owasp_crspolicyencoding_not_allowed_960010" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960010"
  rule_action        = "alert"
}

// GET or HEAD Request with Body Content
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_hreq_960011" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960011"
  rule_action        = "alert"
}

// POST Request Missing Content-Length Header
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_hreq_960012" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960012"
  rule_action        = "alert"
}

// Content-Length HTTP Header is Not Numeric
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_hreq_960016" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960016"
  rule_action        = "alert"
}

// Expect Header Not Allowed For HTTP 1.0
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_hreq_960022" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960022"
  rule_action        = "alert"
}

// HTTP Protocol Version is Not Allowed By Policy
resource "akamai_appsec_rule" "policy1_owasp_crspolicyprotocol_not_allowed_960034" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960034"
  rule_action        = "alert"
}

// URL File Extension is Restricted By Policy
resource "akamai_appsec_rule" "policy1_owasp_crspolicyext_restricted_960035" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960035"
  rule_action        = "alert"
}

// Argument value too long
resource "akamai_appsec_rule" "policy1_owasp_crspolicysize_limit_960208" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960208"
  rule_action        = "alert"
}

// Argument name too long
resource "akamai_appsec_rule" "policy1_owasp_crspolicysize_limit_960209" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960209"
  rule_action        = "alert"
}

// Too many arguments in request
resource "akamai_appsec_rule" "policy1_owasp_crspolicysize_limit_960335" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960335"
  rule_action        = "alert"
}

// Total arguments size exceeded
resource "akamai_appsec_rule" "policy1_owasp_crspolicysize_limit_960341" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960341"
  rule_action        = "alert"
}

// Invalid character in request
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationevasion_960901" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960901"
  rule_action        = "alert"
}

// Invalid Use of Identity Encoding
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_hreq_960902" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960902"
  rule_action        = "alert"
}

// Request Containing Content, but Missing Content-Type header
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationmissing_header_960904" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960904"
  rule_action        = "alert"
}

// Failed to Parse Request Body
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_req_960912" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960912"
  rule_action        = "alert"
}

// Multipart Request Body Failed Strict Validation
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_req_960913" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960913"
  rule_action        = "alert"
}

// Multipart Parser Detected a Possible Unmatched Boundary
resource "akamai_appsec_rule" "policy1_owasp_crsprotocol_violationinvalid_req_960914" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "960914"
  rule_action        = "alert"
}

// Possible XSS Attack Detected - HTML Tag Handler
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973300" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973300"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973301" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973301"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973302" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973302"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973303" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973303"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973304" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973304"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (URL Protocols)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_973305" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973305"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973306" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973306"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Eval/Atob Functions)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_973307" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973307"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973308" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973308"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973309" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973309"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973310" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973310"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (XSS Unicode PoC String)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_973311" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973311"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Common PoC Payload)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_973312" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973312"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973313" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973313"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973314" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973314"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973315" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973315"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973316" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973316"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973317" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973317"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973318" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973318"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973319" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973319"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973320" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973320"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973321" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973321"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973322" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973322"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973323" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973323"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973324" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973324"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973325" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973325"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973326" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973326"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973327" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973327"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973328" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973328"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973329" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973329"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973330" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973330"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973331" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973331"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973332" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973332"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973333" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973333"
  rule_action        = "alert"
}

// IE XSS Filters - Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973334" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973334"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (IE XSS Filter Evasion Attempt)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_973335" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973335"
  rule_action        = "alert"
}

// XSS Filter - Category 1: Script Tag Vector
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973336" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973336"
  rule_action        = "alert"
}

// XSS Filter - Category 2: Event Handler Vector
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_973337" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "973337"
  rule_action        = "alert"
}

// Restricted SQL Character Anomaly Detection Alert - Total # of special characters exceeded
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackspecial_chars_981173" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981173"
  rule_action        = "alert"
}

// Conditional SQL Injection Attempts
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_981241" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981241"
  rule_action        = "alert"
}

// SQL Injection Attack (SQL Operator and Expression Probes 1)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981242" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981242"
  rule_action        = "alert"
}

// SQL Injection Attack (SQL Operator and Expression Probes 2)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981243" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981243"
  rule_action        = "alert"
}

// SQL Injection Attack (Tautology Probes 2)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981244" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981244"
  rule_action        = "alert"
}

// Basic SQL Authentication Bypass Attempts 2/3
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_981245" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981245"
  rule_action        = "alert"
}

// Basic SQL Authentication Bypass Attempts 3/3
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_981246" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981246"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 3)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981247" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981247"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 2)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981248" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981248"
  rule_action        = "alert"
}

// Chained SQL Injection Attempts 2/2
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_981249" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981249"
  rule_action        = "alert"
}

// SQL Benchmark And sleep() Injection Attempts Including Conditional Queries
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_981250" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981250"
  rule_action        = "alert"
}

// SQL Injection Attack (Built-in Functions, Objects and Keyword Probes 3)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981251" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981251"
  rule_action        = "alert"
}

// SQL Injection Attack (Charset manipulation)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981252" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981252"
  rule_action        = "alert"
}

// SQL Injection Attack (Stored Procedure Detected)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981253" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981253"
  rule_action        = "alert"
}

// SQL Injection Attack (Time-based Blind Probe)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981254" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981254"
  rule_action        = "alert"
}

// SQL Injection Attack (Sysadmin access functions)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981255" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981255"
  rule_action        = "alert"
}

// SQL Injection Attack (Merge, Execute, Match Probes)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981256" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981256"
  rule_action        = "alert"
}

// SQL Injection Attack (Hex Encoding Detected)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981260" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981260"
  rule_action        = "alert"
}

// SQL Injection Attack (NoSQL MongoDB Probes)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981270" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981270"
  rule_action        = "alert"
}

// Blind SQLi Tests Using sleep() or benchmark()
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_981272" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981272"
  rule_action        = "alert"
}

// SQL Injection Attack (UNION Attempt)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981276" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981276"
  rule_action        = "alert"
}

// Integer Overflow Attacks (Taken From Skipfish)
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_981277" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981277"
  rule_action        = "alert"
}

// SQL Injection Attack (SELECT Statement Anomaly Detected)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981300" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981300"
  rule_action        = "alert"
}

// SQL Injection Attack: Common Injection Testing Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_981318" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981318"
  rule_action        = "alert"
}

// SQL Injection Attack: SQL Operator Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_981319" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981319"
  rule_action        = "alert"
}

// SQL Injection Attack (Known/Default DB Resources Probe)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_981320" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "981320"
  rule_action        = "alert"
}

// Request Indicates a Security Scanner Scanned the Site
resource "akamai_appsec_rule" "policy1_owasp_crsautomationsecurity_scanner_990002" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "990002"
  rule_action        = "alert"
}

// Rogue Web Site Crawler
resource "akamai_appsec_rule" "policy1_owasp_crsautomationmalicious_990012" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "990012"
  rule_action        = "alert"
}

// Request Indicates a Security Scanner Scanned the Site
resource "akamai_appsec_rule" "policy1_owasp_crsautomationsecurity_scanner_990901" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "990901"
  rule_action        = "alert"
}

// Request Indicates a Security Scanner Scanned the Site
resource "akamai_appsec_rule" "policy1_owasp_crsautomationsecurity_scanner_990902" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "990902"
  rule_action        = "alert"
}

// SQL Injection Attack (GROUP BY/ORDER BY)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_3000000" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000000"
  rule_action        = "alert"
}

// HTTP Response Splitting (Header Injection Attempt)
resource "akamai_appsec_rule" "policy1_akamaiweb_attackhttp_response_splitting_3000001" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000001"
  rule_action        = "alert"
}

// Local System File Access Attempt
resource "akamai_appsec_rule" "policy1_akamaiweb_attackfile_injection_3000002" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000002"
  rule_action        = "alert"
}

// PHP Code Injection
resource "akamai_appsec_rule" "policy1_akamaiweb_attackphp_injection_3000003" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000003"
  rule_action        = "alert"
}

// Potential Remote File Inclusion (RFI) Attack: Suspicious Off-Domain URL Reference
resource "akamai_appsec_rule" "policy1_aseweb_attackrfi_3000004" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000004"
  rule_action        = "alert"
}

// CMD Injection Attack Detected (OS commands with full path)
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000005" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000005"
  rule_action        = "alert"
}

// SQL Injection Attack (Comment String Termination)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_3000006" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000006"
  rule_action        = "alert"
}

// Command Injection (Unix File Leakage)
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000007" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000007"
  rule_action        = "alert"
}

// Pandora / DirtJumper DDoS Detection - HTTP GET Attacks
resource "akamai_appsec_rule" "policy1_akamaiautomationmalicious_3000008" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000008"
  rule_action        = "alert"
}

// Ruby on Rails YAML Injection Attack
resource "akamai_appsec_rule" "policy1_akamaiweb_attackruby_injection_3000009" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000009"
  rule_action        = "alert"
}

// LOIC 1.1 DoS Detection
resource "akamai_appsec_rule" "policy1_akamaiweb_attackloic11_3000010" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000010"
  rule_action        = "alert"
}

// HULK DoS Attack Tool Detection
resource "akamai_appsec_rule" "policy1_akamaiweb_attackhulk_3000011" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000011"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (OGNL Injection)
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000012" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000012"
  rule_action        = "alert"
}

// System Command Injection (Attacker Toolset Download)
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000013" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000013"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (OGNL Injection)
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000014" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000014"
  rule_action        = "alert"
}

// SQL Injection Attack (Database Timing Query)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_3000015" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000015"
  rule_action        = "alert"
}

// PHP Code Injection Using Data Stream Wrapper
resource "akamai_appsec_rule" "policy1_akamaiweb_attackphp_injection_3000016" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000016"
  rule_action        = "alert"
}

// MySQL Keywords Anomaly Detection Alert
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_3000017" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000017"
  rule_action        = "alert"
}

// DirtJumper DDoS Detection - HTTP POST Attacks
resource "akamai_appsec_rule" "policy1_akamaiautomationmalicious_3000018" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000018"
  rule_action        = "alert"
}

// Pandora DDoS Detection - HTTP POST Attacks
resource "akamai_appsec_rule" "policy1_akamaiautomationmalicious_3000019" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000019"
  rule_action        = "alert"
}

// Local File Inclusion (and Command Injection) Using '/proc/self/environ'
resource "akamai_appsec_rule" "policy1_akamaiweb_attacklfi_3000020" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000020"
  rule_action        = "alert"
}

// Detect Attempts to Access the Wordpress Pingback API
resource "akamai_appsec_rule" "policy1_akamaiweb_attackwordpress_pingback_3000021" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000021"
  rule_action        = "alert"
}

// SQL Injection (Built-in Functions, Objects and Keyword Probes 4)
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_3000022" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000022"
  rule_action        = "alert"
}

// Apache Struts ClassLoader Manipulation Remote Code Execution
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000023" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000023"
  rule_action        = "alert"
}

// Apache Commons FileUpload and Apache Tomcat DoS
resource "akamai_appsec_rule" "policy1_akamaiweb_attackapache_commons_dos_3000024" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000024"
  rule_action        = "alert"
}

// CVE-2014-6271 Bash Command Injection Attack
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000025" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000025"
  rule_action        = "alert"
}

// XXE External Entity
resource "akamai_appsec_rule" "policy1_akamaiweb_attackxxe_3000027" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000027"
  rule_action        = "alert"
}

// SQL Injection Attack: MySQL comments, conditions and ch(a)r injections
resource "akamai_appsec_rule" "policy1_aseweb_attacksql_injection_3000029" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000029"
  rule_action        = "alert"
}

// Basic SQL Authentication Bypass Attempts 3/3
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attacksql_injection_3000030" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000030"
  rule_action        = "alert"
}

// HTTP.sys Remote Code Execution Vulnerability Attack Detected (CVE-2015-1635)
resource "akamai_appsec_rule" "policy1_akamaiweb_attackiis_range_3000031" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000031"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack Event Handler
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_3000032" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000032"
  rule_action        = "alert"
}

// PHP Wrapper Attack
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000033" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000033"
  rule_action        = "alert"
}

// Command Injection via the Java Runtime.getRuntime() Method
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000034" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000034"
  rule_action        = "alert"
}

// Potential Account Brute Force Guessing via Wordpress XML-RPC API authenticated methods
resource "akamai_appsec_rule" "policy1_akamaiweb_attackwordpress_bruteforce_3000035" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000035"
  rule_action        = "alert"
}

// Detected LOIC / HOIC client request based on query string
resource "akamai_appsec_rule" "policy1_akamaiddosloic_hoic_1_v1_3000036" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000036"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (JS On-Event Handler)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_3000037" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000037"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (DOM Window Properties)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_3000038" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000038"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (DOM Document Methods)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_3000039" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000039"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Document Methods
resource "akamai_appsec_rule" "policy1_akamaiweb_attackxss_3000040" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000040"
  rule_action        = "alert"
}

// Server Side Template Injection (SSTI)
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000041" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000041"
  rule_action        = "alert"
}

// Detected ARDT client request
resource "akamai_appsec_rule" "policy1_akamaiddosardt_3000042" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000042"
  rule_action        = "alert"
}

// Detect Attempts to Access the Wordpress system.multicall XML-RPC API
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksystem_multicall_3000043" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000043"
  rule_action        = "alert"
}

// SQL Injection Using Boolean Logic
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksql_injection_3000044" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000044"
  rule_action        = "alert"
}

// SQL Injection Using Boolean Logic 1
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksql_injection_3000045" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000045"
  rule_action        = "alert"
}

// SQL Injection Using Boolean Logic 2
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksql_injection_3000046" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000046"
  rule_action        = "alert"
}

// SQL Injection Using SQL Backup Command
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksql_injection_3000047" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000047"
  rule_action        = "alert"
}

// SQL Injection Using SQL Restore Command
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksql_injection_3000048" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000048"
  rule_action        = "alert"
}

// SQL Injection With Cursor Declaration
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksql_injection_3000049" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000049"
  rule_action        = "alert"
}

// SQL Injection Using Boolean Logic 3
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksql_injection_3000050" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000050"
  rule_action        = "alert"
}

// SQL Injection Using Boolean Logic 4
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksql_injection_3000051" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000051"
  rule_action        = "alert"
}

// SQL Injection Using EXISTS
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksql_injection_3000052" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000052"
  rule_action        = "alert"
}

// SQL Injection Using DELETE Statements
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksql_injection_3000053" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000053"
  rule_action        = "alert"
}

// SQL Injection Using UPDATE Statements
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksql_injection_3000054" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000054"
  rule_action        = "alert"
}

// Avzhan Bot DDOS Detection
resource "akamai_appsec_rule" "policy1_akamaiautomationmalicious_3000055" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000055"
  rule_action        = "alert"
}

// PHP Object Injection Attack Detected
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000056" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000056"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Common Attack Tool Keywords)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_3000057" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000057"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (OGNL Injection)
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000058" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000058"
  rule_action        = "alert"
}

// Request Headers indicate request came from Wordpress Pingback
resource "akamai_appsec_rule" "policy1_akamaiweb_attackwordpress_pingback_3000059" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000059"
  rule_action        = "alert"
}

// Mirai / Kaiten DDoS Detection - HTTP Attacks
resource "akamai_appsec_rule" "policy1_akamaiautomationmalicious_3000060" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000060"
  rule_action        = "alert"
}

// Cross-site Scripting Attack (Referer Header From OpenBugBounty Website)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_3000061" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000061"
  rule_action        = "alert"
}

// Mirai/Kaiten Bot DDOS Detection - Bogus Search Engine Referer
resource "akamai_appsec_rule" "policy1_akamaiautomationmalicious_3000062" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000062"
  rule_action        = "alert"
}

// Wordpress wp-json Attack Attempt - non-integer character(s) in ID parameter paylaod
resource "akamai_appsec_rule" "policy1_akamaiweb_attackwordpress_3000063" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000063"
  rule_action        = "alert"
}

// Application Layer Hash DoS Attack
resource "akamai_appsec_rule" "policy1_akamaiautomationmalicious_3000064" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000064"
  rule_action        = "alert"
}

// Apache Struts Remote Command Execution (Deserialization Attack)
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000065" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000065"
  rule_action        = "alert"
}

// Potential WordPress Parameter Resource Consumption Remote DoS Attack (CVE-2018-6389)
resource "akamai_appsec_rule" "policy1_akamaiweb_attackwordpress_dos_3000066" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000066"
  rule_action        = "alert"
}

// Potential Drupal Attack (CVE-2018-7600)
resource "akamai_appsec_rule" "policy1_akamaiweb_attackdrupal_3000067" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000067"
  rule_action        = "alert"
}

// ESI injection Attack
resource "akamai_appsec_rule" "policy1_akamaiweb_attackesi_injection_3000068" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000068"
  rule_action        = "alert"
}

// Webshell/Backdoor File Upload Attempt
resource "akamai_appsec_rule" "policy1_owasp_crsmalicious_softwaretrojan_3000071" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000071"
  rule_action        = "alert"
}

// Deserialization Attack Detected
resource "akamai_appsec_rule" "policy1_aseweb_attackcmd_injection_3000072" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000072"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Invalid Transfer-Encoding Header Value
resource "akamai_appsec_rule" "policy1_akamaiprotocol_violationhttp_desync_3000073" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000073"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: HTTP Request Smuggling Detect in Request Body
resource "akamai_appsec_rule" "policy1_akamaiprotocol_violationhttp_desync_3000074" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000074"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Transfer-Encoding Header Name Obfuscation
resource "akamai_appsec_rule" "policy1_akamaiprotocol_violationhttp_desync_3000075" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000075"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Transfer-Encoding Header in Request Body
resource "akamai_appsec_rule" "policy1_akamaiprotocol_violationhttp_desync_3000076" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000076"
  rule_action        = "alert"
}

// Potential HTTP Desync Attack: Chunked header value with invalid Header Name
resource "akamai_appsec_rule" "policy1_akamaiprotocol_violationhttp_desync_3000077" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000077"
  rule_action        = "alert"
}

// Microsoft Sharepoint Remote Command Execution (Deserialization Attack)
resource "akamai_appsec_rule" "policy1_akamaiweb_attacksharepoint_deserial_3000079" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000079"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Attribute Injection 1)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_3000080" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000080"
  rule_action        = "alert"
}

// Cross-site Scripting (XSS) Attack (Attribute Injection 2)
resource "akamai_appsec_rule" "policy1_aseweb_attackxss_3000081" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000081"
  rule_action        = "alert"
}

// XSS Attack Detected
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackxss_3000082" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000082"
  rule_action        = "alert"
}

// Possible MS Exchange/OWA Attack Detected (CVE-2021-26855)
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackplatform_3000083" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000083"
  rule_action        = "alert"
}

// Possible MS Exchange/OWA Attack Detected (CVE-2021-27065)
resource "akamai_appsec_rule" "policy1_owasp_crsweb_attackplatform_3000084" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000084"
  rule_action        = "alert"
}

// PROXY Header Detected
resource "akamai_appsec_rule" "policy1_akamaiweb_attackproxy_header_detected_3000999" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  rule_id            = "3000999"
  rule_action        = "alert"
}


resource "akamai_appsec_custom_rule_action" "policy1_60088542" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_waf_protection.policy1.security_policy_id
  custom_rule_id     = akamai_appsec_custom_rule.custom_rule_1_60088542.custom_rule_id
  custom_rule_action = "deny"
}

// WAF Attack Group Actions
resource "akamai_appsec_attack_group" "policy1_SQL" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy1.security_policy_id
  attack_group        = "SQL"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy1_XSS" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy1.security_policy_id
  attack_group        = "XSS"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy1_CMD" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy1.security_policy_id
  attack_group        = "CMD"
  attack_group_action = "deny"
}

resource "akamai_appsec_attack_group" "policy1_HTTP" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy1.security_policy_id
  attack_group        = "HTTP"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy1_RFI" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy1.security_policy_id
  attack_group        = "RFI"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy1_PHP" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy1.security_policy_id
  attack_group        = "PHP"
  attack_group_action = "deny_custom_78842"
}

resource "akamai_appsec_attack_group" "policy1_TROJAN" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy1.security_policy_id
  attack_group        = "TROJAN"
  attack_group_action = "alert"
}

resource "akamai_appsec_attack_group" "policy1_IN" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy1.security_policy_id
  attack_group        = "IN"
  attack_group_action = "deny"
}

resource "akamai_appsec_attack_group" "policy1_OUT" {
  config_id           = akamai_appsec_configuration.config.config_id
  security_policy_id  = akamai_appsec_waf_protection.policy1.security_policy_id
  attack_group        = "OUT"
  attack_group_action = "alert"
}

