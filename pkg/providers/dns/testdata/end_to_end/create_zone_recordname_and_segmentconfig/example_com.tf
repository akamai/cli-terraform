terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 9.0.0"
    }
  }
}

locals {
  zone = "example.com"
}

module "example_com" {
  source = "testdata/res/create_zone_recordname_and_segmentconfig/modules/example_com"

  contract = var.contractid
  group    = var.groupid
  name     = local.zone
}

module "example_com_abc_example_com_TXT" {
  source = "testdata/res/create_zone_recordname_and_segmentconfig/modules/example_com_abc_example_com_TXT"

  zonename = local.zone
}
