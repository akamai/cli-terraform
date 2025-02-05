resource "akamai_appsec_selected_hostnames" "hostnames" {
  config_id = local.config_id
  hostnames = ["test.akamai.com"]
  mode      = "REPLACE"
}
