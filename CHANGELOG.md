# Release Notes

## Version 0.7.0 (Mar ?, 2022)

### Features / Enhancements

* General
  * Update urfave/cli to v2
    * BREAKING CHANGE: now flags must come before args
  * Update the README to reflect changes in command line flag ordering

### Bug Fixes

* DNS
  * Corrected handling of embedded double quotes in TXT Recordsets

## Version 0.6.0 (Mar 3, 2022)

### Features / Enhancements

* [IMPORTANT] EdgeWorkers and EdgeKV
  * EDGEWORKERS
    * Support importing existing EdgeWorker configuration with related resources
  * EDGEKV
    * Support importing existing EdgeKV configuration with related resource

* CLOUDLETS
  * Support importing existing Cloudlets match rules for Audience Segmentation with related data source

## Version 0.5.1 (Feb 10, 2022)

### Features / Enhancements

* CLOUDLETS
  * Add a `default` value to the `config_section` variable in `create-cloudlets-policy` in case the user did not specify any
  * Set the latest application load balancer version in the resource `akamai_cloudlets_application_load_balancer_activation`

## Version 0.5.0 (Jan 27, 2022)

### Features / Enhancements

* CLOUDLETS
  * Support importing existing Cloudlets match rules for Visitor Prioritization with related data source
  * Support importing existing Cloudlets match rules for Continuous Deployment/Phased Release with related data source
  * Support importing existing Cloudlets match rules for Forward Rewrite with related data source
  * Support importing existing Cloudlets match rules for API Prioritization with related data source

## Version 0.4.0

### Features/Enhancements

* [IMPORTANT] Cloudlets
  * Support importing existing Cloudlets policy and policy versions configuration with related resource
  * Support importing existing Cloudlets application load balancer configuration with related resource
  * Support importing existing Cloudlets match rules for application load balancer with related data source
  * Support importing existing Cloudlets match rules for edge redirector with related data source

## Version 0.3.1

### Bug Fixes

* General
  * Update binaries URL to fix binary installation failure ([#6](https://github.com/akamai/cli-terraform/issues/6))
  * Update dependencies to fix issue under MacOS Big Sur

* PAPI
  * Add field `is_secure` into rule tree structure ([#9](https://github.com/akamai/cli-terraform/issues/9))

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
