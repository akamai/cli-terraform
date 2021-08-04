# Release Notes

## Version 0.3.0

### Features/Enhancements

* General
  * Add `account-key` as alias for `accountkey` argument

* DNS
  * Change edgerc section default to `default`
  * Update generated Terraform config file header
  * Populate contractId variable in dnsvars.tf generation

* GTM
  * Change edgerc section default to `default`
  * Update generated Terraform config file header

* PAPI
  * Remove deprecated CPCode support

## Version 0.2.0 

### Bug Fixes

* DNS
  * Correct AKAMAICDN record resource generation
  * Correct SVCB/HTTPS record resource generation

* GTM
  * Do not include Datacenter ID in Datacenter resource generation
  * Do not include Traffic Target in Static property resource generation

### Features/Enhancements

* PAPI
  * Update property resources to latest syntax

## Version 0.1.0

### Features/Enhancements

Initial release

* DNS
  * Support importing existing DNS zones and related resources

* GTM 
  * Support importing existing GTM domains and related resources

* PAPI
  * Support importing existing PAPI properties and related resources
