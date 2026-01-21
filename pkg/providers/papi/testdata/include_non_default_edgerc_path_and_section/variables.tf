variable "edgerc_path" {
  type    = string
  default = "/non/default/path/to/edgerc"
}

variable "config_section" {
  type    = string
  default = "non_default_section"
}

variable "activate_latest_on_staging" {
  type    = bool
  default = false
}

variable "activate_latest_on_production" {
  type    = bool
  default = false
}
