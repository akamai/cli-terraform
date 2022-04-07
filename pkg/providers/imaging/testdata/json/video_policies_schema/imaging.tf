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
  type        = "VIDEO"
  contract_id = "ctr_123"
}

resource "akamai_imaging_policy_video" "policy__auto" {
  policy_id              = ".auto"
  contract_id            = "ctr_123"
  policyset_id           = "test_policyset_id"
  activate_on_production = true
  policy {
    breakpoints {

      widths = [320, 640, 1024, 2048, 5000]
    }
    output {

      perceptual_quality = "mediumHigh"
    }
  }
}

resource "akamai_imaging_policy_video" "policy_test_policy_video" {
  policy_id              = "test_policy_video"
  contract_id            = "ctr_123"
  policyset_id           = "test_policyset_id"
  activate_on_production = true
  policy {
    breakpoints {

      widths = [420, 640, 1024, 2048, 5000]
    }
    output {

      perceptual_quality = "mediumHigh"
    }
  }
}
