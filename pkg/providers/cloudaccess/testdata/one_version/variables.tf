variable "edgerc_path" {
  type    = string
  default = "~/.edgerc"
}

variable "config_section" {
  type    = string
  default = "test_section"
}

variable "secret_access_key_a" {
  type      = string
  sensitive = true
}