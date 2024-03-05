terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 5.2.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_cloudwrapper_configuration" "test_configuration" {
  config_name  = "test_configuration"
  contract_id  = "1234"
  property_ids = ["123", "456"]
  comments = trimsuffix(<<EOT
first
second

last
EOT
  , "\n")
  retain_idle_objects = false
  location {
    traffic_type_id = 1
    comments        = <<EOT
first
second
EOT
    capacity {
      value = 1
      unit  = "GB"
    }
  }
  location {
    traffic_type_id = 2
    comments        = "TestComments"
    capacity {
      value = 2
      unit  = "TB"
    }
  }
}

#resource "akamai_cloudwrapper_activation" "test_configuration_activation" {
#  config_id = akamai_cloudwrapper_configuration.test_configuration.id
#  revision  = akamai_cloudwrapper_configuration.test_configuration.revision
#}