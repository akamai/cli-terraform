resource "akamai_iam_role" "role_id_12345" {
  name          = "Custom role"
  description   = "Custom role description"
  granted_roles = [992, 707, 452, 677, 726, 296, 457, 987]
}

