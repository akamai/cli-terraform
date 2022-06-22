resource "akamai_appsec_selected_hostnames" "hostnames" {
  config_id = akamai_appsec_configuration.config.config_id
  hostnames = ["www.easyakamai.com", "konaneweahost9001.edgekey.net", "konaneweahost8012.edgekey.net", "konaneweahost9002.edgekey.net", "konaneweahost8013.edgekey.net", "www.vbhat.com", "www.andrew89.com", "konaneweahost8016.edgekey.net", "konaneweahost8014.edgekey.net", "aetsaitcwest.edgekey.net", "konaneweahost9000.edgekey.net"]
  mode      = "REPLACE"
}
