resource "akamai_appsec_selected_hostnames" "hostnames" {
  config_id = akamai_appsec_configuration.config.config_id
  hostnames = ["test.akamai.com"]
  mode      = "REPLACE"
}
