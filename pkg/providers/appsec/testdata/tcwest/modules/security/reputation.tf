// Client Reputation Actions
resource "akamai_appsec_reputation_profile_action" "andrew_2670508" {
  config_id             = akamai_appsec_configuration.config.config_id
  security_policy_id    = akamai_appsec_reputation_protection.andrew.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.web_attackers_high_threat.reputation_profile_id
  action                = "deny"
}
resource "akamai_appsec_reputation_profile_action" "andrew_2670509" {
  config_id             = akamai_appsec_configuration.config.config_id
  security_policy_id    = akamai_appsec_reputation_protection.andrew.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.dos_attackers_high_threat.reputation_profile_id
  action                = "alert"
}
// Client Reputation Actions
resource "akamai_appsec_reputation_profile_action" "policy1_2670508" {
  config_id             = akamai_appsec_configuration.config.config_id
  security_policy_id    = akamai_appsec_reputation_protection.policy1.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.web_attackers_high_threat.reputation_profile_id
  action                = "alert"
}
resource "akamai_appsec_reputation_profile_action" "policy1_2670509" {
  config_id             = akamai_appsec_configuration.config.config_id
  security_policy_id    = akamai_appsec_reputation_protection.policy1.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.dos_attackers_high_threat.reputation_profile_id
  action                = "deny"
}
resource "akamai_appsec_reputation_profile_action" "policy1_2670510" {
  config_id             = akamai_appsec_configuration.config.config_id
  security_policy_id    = akamai_appsec_reputation_protection.policy1.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.scanning_tools_high_threat.reputation_profile_id
  action                = "deny"
}
resource "akamai_appsec_reputation_profile_action" "policy1_2670511" {
  config_id             = akamai_appsec_configuration.config.config_id
  security_policy_id    = akamai_appsec_reputation_protection.policy1.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.web_attackers_low_threat.reputation_profile_id
  action                = "alert"
}
resource "akamai_appsec_reputation_profile_action" "policy1_2670512" {
  config_id             = akamai_appsec_configuration.config.config_id
  security_policy_id    = akamai_appsec_reputation_protection.policy1.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.dos_attackers_low_threat.reputation_profile_id
  action                = "alert"
}
resource "akamai_appsec_reputation_profile_action" "policy1_2670513" {
  config_id             = akamai_appsec_configuration.config.config_id
  security_policy_id    = akamai_appsec_reputation_protection.policy1.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.scanning_tools_low_threat.reputation_profile_id
  action                = "alert"
}
resource "akamai_appsec_reputation_profile_action" "policy1_2670514" {
  config_id             = akamai_appsec_configuration.config.config_id
  security_policy_id    = akamai_appsec_reputation_protection.policy1.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.web_scrapers_low_threat.reputation_profile_id
  action                = "alert"
}
resource "akamai_appsec_reputation_profile_action" "policy1_2670515" {
  config_id             = akamai_appsec_configuration.config.config_id
  security_policy_id    = akamai_appsec_reputation_protection.policy1.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.web_scrapers_high_threat.reputation_profile_id
  action                = "deny"
}
