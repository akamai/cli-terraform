variable "edgerc_path" {
  type    = string
  default = "~/.edgerc"
}

variable "config_section" {
  type    = string
  default = "default"
}

variable "activate_latest_on_staging" {
  type    = bool
  default = false
}

#variable "activate_latest_on_production" {
#  type    = bool
#  default = true
#}
