terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 2.0.0"
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

data "akamai_imaging_policy_video" "data_policy_test_policy_video" {
  policy {
    breakpoints {

      widths = [280, 1080]
    }
    hosts = ["host1", "host2"]
    output {

      perceptual_quality    = "mediumHigh"
      placeholder_video_url = "some"
    }
    rollout_duration = 3600
    variables {

      default_value = "280"
      name          = "ResizeDim"
      type          = "number"
    }
    variables {

      default_value = "260"
      name          = "ResizeDimWithBorder"
      type          = "number"
    }
    variables {

      default_value = "1000"
      enum_options {

        id    = "1"
        value = "value1"
      }
      enum_options {

        id    = "2"
        value = "value2"
      }
      name = "MinDim"
      type = "number"
    }
    variables {

      default_value = "1450"
      name          = "MinDimNew"
      type          = "number"
    }
    variables {

      default_value = "1500"
      name          = "MaxDimOld"
      type          = "number"
    }
  }
}

resource "akamai_imaging_policy_video" "policy_test_policy_video" {
  policy_id              = "test_policy_video"
  contract_id            = "ctr_123"
  policyset_id           = akamai_imaging_policy_set.policyset.id
  activate_on_production = true
  json                   = data.akamai_imaging_policy_video.data_policy_test_policy_video.json
}
