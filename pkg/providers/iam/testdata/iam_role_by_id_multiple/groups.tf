resource "akamai_iam_group" "group_id_56789" {
  parent_group_id = 12345
  name            = "Custom group 1"
}

resource "akamai_iam_group" "group_id_98765" {
  parent_group_id = 12345
  name            = "Custom group 2"
}

