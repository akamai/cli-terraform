resource "akamai_iam_ip_allowlist" "allowlist" {
  enable = true
}

resource "akamai_iam_cidr_blocks" "cidr_blocks" {
  cidr_blocks = [
    {
      cidr_block = "1.1.1.1/1"
      enabled    = true
      comments   = "comment"
    },
    {
      cidr_block = "2.2.2.2/2"
      enabled    = false
    },
  ]
}
