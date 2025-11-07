terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 9.2.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_property_domainownership_domains" "example-com" {
  domains = [
    {
      domain_name      = "example.com"
      validation_scope = "DOMAIN"
    },
    {
      domain_name      = "sub.example.com"
      validation_scope = "HOST"
    }
  ]
}

resource "akamai_property_domainownership_validation" "example-com" {
  domains = [
    {
      domain_name      = "example.com"
      validation_scope = "DOMAIN"
    },
    {
      domain_name      = "sub.example.com"
      validation_scope = "HOST"
    }
  ]
}
