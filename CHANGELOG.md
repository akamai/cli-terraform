# Release Notes

## Version 1.4.x (Apr xx, 2023)

### Features/Enhancements

### Fixes

* GTM
  * Remove deprecated field `name` of `traffic_target` during export  ([I#374](https://github.com/akamai/terraform-provider-akamai/issues/374))

* PAPI
  * `is_secure` and `variable` fields can only be used in `default` datasource `akamai_property_rules_builder`

## Version 1.4.0 (Mar 30, 2023)

### Features/Enhancements

* APPSEC
  * Support for exporting `akamai_appsec_advanced_settings_request_body` resource

* PAPI
  * New `--schema` flag available with `export-property` command to export rule tree as `akamai_property_rules_builder` data source (Beta)

### Fixes

* PAPI
  * Fix property export with empty EdgeHostnameID ([I#41](https://github.com/akamai/cli-terraform/issues/41))
  * Change exported attribute value in configuration of `akamai_contract` data source from deprecated `name`
    to `group_name` when exporting property
  * Comment out `akamai_property_activation` resource if there is no currently active version

## Version 1.3.1 (Feb 2, 2023)

### Features/Enhancements

* General
  * Migrate to go 1.18
  * Add badges to readme and improve code quality based on golangci-lint
* CPS
  * Add `preferred_trust_chain` in `csr` set attribute

## Version 1.3.0 (Dec 15, 2022)

### Features/Enhancements

* PAPI
  * New `export-property include` subcommand to export Property Manager Include `akamai_property_include` with accompanying resources and data sources:
  `akamai_property_include_activation`, `akamai_property_include_parents` and `akamai_property_rules_template`
  * New `--with-includes` flag available with `export-property` command to export resources and data sources for the Property Manager Includes that are referenced by the Property which is being exported

## Version 1.2.0 (Dec 1, 2022)

### Features/Enhancements

* CPS
  * New `export-cps` command to export DV enrollment (`akamai_cps_dv_enrollment`) or third-party enrollment with accompanying resources and data source (`akamai_cps_third_party_enrollment`,`akamai_cps_csr` and `akamai_cps_upload_certificate`)
  
## Version 1.1.1 (Oct 27, 2022)

### Fixes

* GTM
  * Fix exporting GTM property with default datacenter ([I#31](https://github.com/akamai/cli-terraform/issues/31))

* PAPI
  * Fix `cert_provisioning_type` field exporting ([I#15](https://github.com/akamai/cli-terraform/issues/15))

## Version 1.1.0 (Oct 10, 2022)

### Features/Enhancements

* Application Security
  * Add import support for `malware_policy` and `malware_policy_action`

### Fixes

* General
  * Resolve all the tflint warnings which were introduced with the tflint version v0.40.0

* Application Security
  * Fix incorrect policy ID for malware protection
  * Fix drift on match targets & malware protection resources

* PAPI
  * Fix ignoring property version for property-snippets during exporting some property

## Version 1.0.0 (Aug 3, 2022)

### Deprecations

* [IMPORTANT] General
  * `create-*` command names are now deprecated, use `export-*` instead

### Features/Enhancements

* GTM
  * Improve formatting of output configurations
* DNS
  * Add support for additional default datacenters
  * Improve formatting of output configurations

### Fixes

* General
  * Fix default flag values in help output
* Identity and Access Management (IAM)
  * Fix IAM role export failures with broken user

## Version 0.9.0 (Jul 07, 2022)

### Features / Enhancements

* Identity and Access Management  (IAM)
  * New `create-iam` command to export users, roles and/or groups
* Application Security
  * New `create-appsec` command to create Application Security resource

## Version 0.8.0 (Jun 02, 2022)

### Features / Enhancements

* General
  * Add `arm64` support (Apple M1)
* Image and Video Manager
  * New `create-imaging` command to import image and video policies
* PAPI
  * Add optional `version` flag to `create-property` command ([I#8](https://github.com/akamai/cli-terraform/issues/8))

## Version 0.7.1 (May 11, 2022)

### Bug fixes

* General
  * `./cli-terraform --help` should return a zero status

* PAPI
  * Normalize rule names in `create-property` command

### Features / Enhancements

* General
  * Support `--acountkey` flag
  * Logging improvements

## Version 0.7.0 (Mar 31, 2022)

### Features / Enhancements

* [IMPORTANT] Image and Video Manager
  * Support importing existing Video and Image Policies
  * Support importing existing Policy Set

* CLOUDLETS
  * Support importing existing Cloudlets match rules for Request Control with related data source

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
