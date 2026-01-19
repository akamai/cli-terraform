terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 9.3.0"
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
    }
  ]
}

resource "akamai_property_domainownership_validation" "example-com" {
  domains = [
    {
      domain_name       = "example.com"
      validation_scope  = "DOMAIN"
      validation_method = "DNS_CNAME"
    }
  ]
}
