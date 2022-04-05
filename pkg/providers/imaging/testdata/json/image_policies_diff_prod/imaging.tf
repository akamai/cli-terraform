terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_imaging_policy_set" "policyset" {
  name        = "some policy set"
  region      = "EMEA"
  type        = "IMAGE"
  contract_id = "ctr_123"
}

resource "akamai_imaging_policy_image" "policy__auto" {
  policy_id              = ".auto"
  contract_id            = "ctr_123"
  policyset_id           = "test_policyset_id"
  activate_on_production = true
  json                   = file("testdata/res/json/image_policies_diff_prod/_auto.json")
}

resource "akamai_imaging_policy_image" "policy_test_policy_image" {
  policy_id              = "test_policy_image"
  contract_id            = "ctr_123"
  policyset_id           = "test_policyset_id"
  activate_on_production = false
  json                   = file("testdata/res/json/image_policies_diff_prod/test_policy_image.json")
}
