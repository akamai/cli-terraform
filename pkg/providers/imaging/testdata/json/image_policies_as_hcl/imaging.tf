terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
    }
  }
  required_version = ">= 1.0"
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

data "akamai_imaging_policy_image" "data_policy__auto" {
  policy {
    breakpoints {

      widths = [320, 640, 1024, 2048, 5000]
    }
    output {

      perceptual_quality = "mediumHigh"
    }
    transformations {
      max_colors {
        colors = 2
      }
    }
  }
}

data "akamai_imaging_policy_image" "data_policy_test_policy_image" {
  policy {
    breakpoints {

      widths = [420, 640, 1024, 2048, 5000]
    }
    output {

      perceptual_quality = "mediumHigh"
    }
    transformations {
      max_colors {
        colors = 2
      }
    }
  }
}

resource "akamai_imaging_policy_image" "policy__auto" {
  policy_id              = ".auto"
  contract_id            = "ctr_123"
  policyset_id           = akamai_imaging_policy_set.policyset.id
  activate_on_production = true
  json                   = data.akamai_imaging_policy_image.data_policy__auto.json
}

resource "akamai_imaging_policy_image" "policy_test_policy_image" {
  policy_id              = "test_policy_image"
  contract_id            = "ctr_123"
  policyset_id           = akamai_imaging_policy_set.policyset.id
  activate_on_production = true
  json                   = data.akamai_imaging_policy_image.data_policy_test_policy_image.json
}
