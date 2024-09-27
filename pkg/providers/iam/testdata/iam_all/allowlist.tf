resource "akamai_iam_ip_allowlist" "allowlist" {
  enable = true
}

resource "akamai_iam_cidr_block" "cidr_1_1_1_1-1" {
  cidr_block = "1.1.1.1/1"
  enabled    = true
  comments   = "comment"
}

resource "akamai_iam_cidr_block" "cidr_2_2_2_2-2" {
  cidr_block = "2.2.2.2/2"
  enabled    = false
}

