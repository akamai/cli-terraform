resource "akamai_iam_role" "role_id_12345" {
  name          = "Custom role 12345"
  description   = "Custom role description"
  granted_roles = [992, 707, 452, 677, 726, 296, 457, 987]
}

resource "akamai_iam_role" "role_id_54321" {
  name          = "Custom role 54321"
  description   = "Custom role description"
  granted_roles = [992, 707, 452, 677, 726, 296, 457, 987]
}

