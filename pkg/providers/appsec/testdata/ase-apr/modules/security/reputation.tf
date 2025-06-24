// Client Reputation Actions
resource "akamai_appsec_reputation_profile_action" "default_policy_3017089" {
  config_id             = local.config_id
  security_policy_id    = akamai_appsec_reputation_protection.default_policy.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.dos_attackers_high_threat.reputation_profile_id
  action                = "deny"
}
