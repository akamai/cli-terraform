variable "name" {
  type = string
}

variable "config_id" {
  type = number
}

variable "note" {
  type = string
}

variable "network" {
  type = string
}

variable "notification_emails" {
  type = list(string)
}
