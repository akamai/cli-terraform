variable "group_name" {
  type    = string
  default = ""
}

variable "contract_id" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = "www.vbhat.com"
}

variable "description" {
  type    = string
  default = "A security config for demo"
}

variable "hostnames" {
  type    = list(string)
  default = ["www.easyakamai.com", "konaneweahost9001.edgekey.net", "konaneweahost8012.edgekey.net", "konaneweahost9002.edgekey.net", "konaneweahost8013.edgekey.net", "www.vbhat.com", "www.andrew89.com", "konaneweahost8016.edgekey.net", "konaneweahost8014.edgekey.net", "aetsaitcwest.edgekey.net", "konaneweahost9000.edgekey.net"]
}

variable "emails" {
  type    = list(string)
  default = ["noreply@example.org"]
}

variable "activation_note" {
  type    = string
  default = "Activated by Terraform"
}

variable "network" {
  type    = string
  default = "STAGING"
}
