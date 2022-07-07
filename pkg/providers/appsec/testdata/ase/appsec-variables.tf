variable "group_name" {
  default = ""
}

variable "contract_id" {
  default = ""
}

variable "name" {
  default = "TFDEMO"
}

variable "description" {
  default = "A security config for demo"
}

variable "hostnames" {
  default = ["test.akamai.com"]
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
