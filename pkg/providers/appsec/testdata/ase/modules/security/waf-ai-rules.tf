// WAF AI Rules
// WAF AI Rules for policy default_policy
resource "akamai_appsec_waf_ai_rules" "default_policy_status" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.default_policy.security_policy_id
  ai_rule_status     = "ENABLED"
}

resource "akamai_appsec_waf_ai_rules" "default_policy_3001000" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.default_policy.security_policy_id
  rule_id            = 3001000
  action             = "alert"
}

resource "akamai_appsec_waf_ai_rules" "default_policy_3001001" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.default_policy.security_policy_id
  rule_id            = 3001001
  action             = "deny"
}

