resource "akamai_appsec_wap_selected_hostnames" "policy2" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.policy2.security_policy_id
  protected_hosts    = ["konaneweahost8012.edgekey.net", "konaneweahost9001.edgekey.net", "konaneweahost8013.edgekey.net", "konaneweahost9002.edgekey.net", "konaneweahost8014.edgekey.net", "konaneweahost8016.edgekey.net", "aetsaitcwest.edgekey.net", "konaneweahost9000.edgekey.net"]
  evaluated_hosts    = []
}

resource "akamai_appsec_wap_selected_hostnames" "andrew" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.andrew.security_policy_id
  protected_hosts    = ["www.easyakamai.com", "www.andrew89.com"]
  evaluated_hosts    = []
}

resource "akamai_appsec_wap_selected_hostnames" "policy1" {
  config_id          = akamai_appsec_configuration.config.config_id
  security_policy_id = akamai_appsec_security_policy.policy1.security_policy_id
  protected_hosts    = ["www.vbhat.com"]
  evaluated_hosts    = []
}

