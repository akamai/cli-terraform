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

data "akamai_imaging_policy_image" "data_policy_test_policy_image" {
  policy {
    output {

      allow_pristine_on_downsize = true
      perceptual_quality         = "mediumHigh"
      prefer_modern_formats      = false
    }
    serve_stale_duration = 3600
    transformations {
      append {
        gravity          = "Center"
        gravity_priority = "horizontal"
        image {
          box_image {
            transformation {

              compound {
                append {
                  gravity          = "Center"
                  gravity_priority = "horizontal"
                  image {
                    box_image {
                      transformation {

                        compound {
                          append {
                            gravity          = "Center"
                            gravity_priority = "horizontal"
                            image {
                              box_image {
                                transformation {

                                }
                              }
                            }
                            preserve_minor_dimension = false
                          }
                        }
                      }
                    }
                  }
                  preserve_minor_dimension = false
                }
              }
            }
          }
        }
        preserve_minor_dimension = false
      }
    }
  }
}

resource "akamai_imaging_policy_image" "policy_test_policy_image" {
  policy_id              = "test_policy_image"
  contract_id            = "ctr_123"
  policyset_id           = akamai_imaging_policy_set.policyset.id
  activate_on_production = true
  json                   = data.akamai_imaging_policy_image.data_policy_test_policy_image.json
}
