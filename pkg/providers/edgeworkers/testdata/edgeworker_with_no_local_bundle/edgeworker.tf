terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = "~> 2.0.0"
    }
  }
  required_version = ">= 0.13"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_edgeworker" "edgeworker" {
  name             = "test_edgeworker"
  group_id         = 1
  resource_tier_id = 2
  // Local Bundle will default to helloworld.tgz from https://github.com/akamai/edgeworkers-examples/tree/master/edgecompute/examples/getting-started/hello-world%20(EW)
}
