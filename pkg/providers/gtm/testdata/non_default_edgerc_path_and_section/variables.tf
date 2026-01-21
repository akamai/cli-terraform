variable "edgerc_path" {
  type    = string
  default = "/non/default/path/to/edgerc"
}

variable "config_section" {
  type    = string
  default = "non_default_section"
}

variable "contractid" {
  type        = string
  default     = ""
  description = "Value unknown at the time of import. Please update."
}

variable "groupid" {
  type        = string
  default     = ""
  description = "Value unknown at the time of import. Please update."
}
