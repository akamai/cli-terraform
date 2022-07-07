resource "akamai_iam_group" "group_id_101" {
  parent_group_id = 111
  name            = "grp_101"
}

resource "akamai_iam_group" "group_id_102" {
  parent_group_id = 111
  name            = "grp_102"
}

resource "akamai_iam_group" "group_id_112" {
  parent_group_id = 102
  name            = "grp_112"
}

resource "akamai_iam_group" "group_id_123" {
  parent_group_id = 112
  name            = "grp_123"
}

