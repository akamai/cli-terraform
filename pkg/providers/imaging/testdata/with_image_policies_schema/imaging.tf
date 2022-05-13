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

data "akamai_imaging_policy_image" "data_policy_test_policy_image" {
  policy {
    breakpoints {

      widths = [280, 1080]
    }
    hosts = ["host1", "host2"]
    output {

      adaptive_quality   = 50
      perceptual_quality = "mediumHigh"
    }
    post_breakpoint_transformations {
      background_color {
        color = "#ffffff"
      }
    }
    post_breakpoint_transformations {
      if_dimension {
        default {

          compound {
            background_color {
              color = "#ffffff"
            }
          }
        }
        dimension = "height"
        greater_than {

          compound {
            background_color {
              color = "#ffffff"
            }
          }
        }
        value_var = "MaxDimOld"
      }
    }
    rollout_duration = 3600
    transformations {
      region_of_interest_crop {
        gravity = "Center"
        height  = 8
        region_of_interest {
          rectangle_shape {
            anchor {
              point_shape {
                x = 4
                y = 5
              }
            }
            height = 9
            width  = 8
          }
        }
        style = "fill"
        width = 7
      }
    }
    transformations {
      append {
        gravity          = "Center"
        gravity_priority = "horizontal"
        image {
          text_image {
            fill        = "#000000"
            size        = 72
            stroke      = "#FFFFFF"
            stroke_size = 0
            text        = "test"
            transformation {

            }
          }
        }
        preserve_minor_dimension = true
      }
    }
    transformations {
      trim {
        fuzz    = 0.08
        padding = 0
      }
    }
    transformations {
      if_dimension {
        default {

          compound {
            if_dimension {
              default {

                compound {
                  if_dimension {
                    default {

                      compound {
                        if_dimension {
                          default {

                            compound {
                              resize {
                                aspect     = "fit"
                                height_var = "ResizeDim"
                                type       = "normal"
                                width_var  = "ResizeDim"
                              }
                            }
                            compound {
                              crop {
                                allow_expansion = true
                                gravity         = "Center"
                                height_var      = "ResizeDim"
                                width_var       = "ResizeDim"
                                x_position      = 0
                                y_position      = 0
                              }
                            }
                            compound {
                              background_color {
                                color = "#ffffff"
                              }
                            }
                          }
                          dimension = "height"
                          greater_than {

                            compound {
                              resize {
                                aspect     = "fit"
                                height_var = "ResizeDimWithBorder"
                                type       = "normal"
                                width_var  = "ResizeDimWithBorder"
                              }
                            }
                            compound {
                              crop {
                                allow_expansion = true
                                gravity         = "Center"
                                height_var      = "ResizeDim"
                                width_var       = "ResizeDim"
                                x_position      = 0
                                y_position      = 0
                              }
                            }
                            compound {
                              background_color {
                                color = "#ffffff"
                              }
                            }
                          }
                          value_var = "MaxDimOld"
                        }
                      }
                    }
                    dimension = "height"
                    less_than {

                      compound {
                        resize {
                          aspect     = "fit"
                          height_var = "ResizeDimWithBorder"
                          type       = "normal"
                          width_var  = "ResizeDimWithBorder"
                        }
                      }
                      compound {
                        crop {
                          allow_expansion = true
                          gravity         = "Center"
                          height_var      = "ResizeDim"
                          width_var       = "ResizeDim"
                          x_position      = 0
                          y_position      = 0
                        }
                      }
                      compound {
                        background_color {
                          color = "#ffffff"
                        }
                      }
                    }
                    value_var = "MinDim"
                  }
                }
              }
              dimension = "width"
              less_than {

                compound {
                  resize {
                    aspect     = "fit"
                    height_var = "ResizeDimWithBorder"
                    type       = "normal"
                    width_var  = "ResizeDimWithBorder"
                  }
                }
                compound {
                  crop {
                    allow_expansion = true
                    gravity         = "Center"
                    height_var      = "ResizeDim"
                    width_var       = "ResizeDim"
                    x_position      = 0
                    y_position      = 0
                  }
                }
                compound {
                  background_color {
                    color = "#ffffff"
                  }
                }
              }
              value_var = "MinDim"
            }
          }
        }
        dimension = "width"
        value_var = "MaxDimOld"
      }
    }
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

      default_value = ""
      name          = "VariableWithoutDefaultValue"
      type          = "string"
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

resource "akamai_imaging_policy_image" "policy_test_policy_image" {
  policy_id              = "test_policy_image"
  contract_id            = "ctr_123"
  policyset_id           = "test_policyset_id"
  activate_on_production = true
  json                   = data.akamai_imaging_policy_image.data_policy_test_policy_image.json
}
