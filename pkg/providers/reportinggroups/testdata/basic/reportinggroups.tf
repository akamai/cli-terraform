terraform {
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 10.1.0"
    }
  }
  required_version = ">= 1.0"
}

provider "akamai" {
  edgerc         = var.edgerc_path
  config_section = var.config_section
}

resource "akamai_reportinggroups_group" "test_reporting_group" {
  reporting_group_name = "test reporting group"

  access_group = {
    contract_id = "1-ACCGRP"
  }

  contract = {
    contract_id = "1-CNTR"
    cp_codes = [
      {
        cp_code_id = "12345"
      },
    ]
  }
}
