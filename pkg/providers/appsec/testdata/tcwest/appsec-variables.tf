variable "group_name" {
  default = ""
}

variable "contract_id" {
  default = ""
}

variable "name" {
  default = "www.vbhat.com"
}

variable "description" {
  default = "A security config for demo"
}

variable "hostnames" {
  default = ["www.easyakamai.com", "konaneweahost9001.edgekey.net", "konaneweahost8012.edgekey.net", "konaneweahost9002.edgekey.net", "konaneweahost8013.edgekey.net", "www.vbhat.com", "www.andrew89.com", "konaneweahost8016.edgekey.net", "konaneweahost8014.edgekey.net", "aetsaitcwest.edgekey.net", "konaneweahost9000.edgekey.net"]
}

variable "emails" {
  default = ["noreply@example.org"]
}

variable "activation_note" {
  default = "Activated by Terraform"
}

variable "network" {
  default = "STAGING"
}
