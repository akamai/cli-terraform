resource "akamai_iam_group" "group_id_56789" {
  parent_group_id = 98765
  name            = "Custom group 56789"
}

resource "akamai_iam_group" "group_id_987" {
  parent_group_id = 98765
  name            = "Custom group 987"
}

