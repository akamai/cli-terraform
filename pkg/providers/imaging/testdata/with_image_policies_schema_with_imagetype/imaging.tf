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

resource "akamai_imaging_policy_image" "policy_test_policy_image" {
  policy_id              = "test_policy_image"
  contract_id            = "ctr_123"
  policyset_id           = "test_policyset_id"
  activate_on_production = true
  policy {
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
