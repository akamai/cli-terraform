# Release Notes

## Version 1.11.0 (Dec 7, 2023)

### Features/Enhancements

* APPSEC
  * Added support for `asn_network_lists` to `akamai_appsec_ip_geo` resource for IP/Geo Firewall.
* Deprecated `--schema` flag and replaced with
  * `--policy-as-hcl` (for export-imaging command)
  * `--rules-as-hcl`  (for export-property command)
* Cloudlets
  * Add `origin_description` field export in `akamai_cloudlets_application_load_balancer` resource
* PAPI
  * `export-property` command with flag `--rules-as-hcl` now supports export of properties in frozen format `v2023-10-30`

### Bug fixes

* DNS (export-zone)
  * Fixed `target` field string escaping in `akamai_dns_record` resource
  * Changed provider version requirement from `~> 1.6.1` to `>= 1.6.1`
* PAPI
  * Fixed error with exporting schema containing `serialNumber`

## Version 1.10.0 (October 31, 2023)

### Features/Enhancements

* [IMPORTANT] Client Lists
  * Added command `export-clientlist` which allows export of `akamai_clientlist_list` and `akamai_clientlist_activation`
    resources
* Cloudlets
  * Added `matches_always` field to `akamai_cloudlets_edge_redirector_match_rule` export template
* PAPI
  * Added support for new rule format `v2023-09-20`

### Bug fixes

* Fixed generation of multiline text for:
  * `description` variable in AppSec configuration
  * `comments` and `location.comments` fields in `akamai_cloudwrapper_configuration`
  * `comment` field in `akamai_dns_zone`
  * `comment` field in `akamai_gtm_domain`
  * `comments` field in `akamai_gtm_property`
  * `description` field in `akamai_gtm_resource`
  * `note` field in `akamai_property_activation` and `akamai_property_include_activation`

## Version 1.9.1 (September 26, 2023)

### Bug fixes

* CPS
  * Fixed nil pointer evaluating *cps.DNSNameSettings.CloneDNSNames ([#52](https://github.com/akamai/cli-terraform/issues/52))

* Identity and Access Management (IAM)
  * Fixed newline escaping in `description` field after exporting a role

* PAPI
  * Add missing fields to `akamai_property_builder` for `origin` and `siteShield` behaviors

## Version 1.9.0 (August 29, 2023)


* [IMPORTANT] CloudWrapper
  * Added support for `export-cloudwrapper` command which allows export of `akamai_cloudwrapper_configuration` and `akamai_cloudwrapper_activation` resources

* APPSEC
  * Added import support for custom client sequence

### Bug fixes

* Image and Video Manager
  * Added description for `--schema` flag for `export-imaging` command in `README.md` ([#56](https://github.com/akamai/cli-terraform/issues/56))

* PAPI
  * Fixed `export-property` command to export `akamai_property_activation` resource attributes for latest active version.
  * Fixed `export-property` command to use `group_id` and `contract_id` as terraform variables, instead of data sources, which
  produced inconsistencies ([I#374](https://github.com/akamai/terraform-provider-akamai/issues/426))
  * `logStreamName` field from `datastream` behavior has changed from string to array of strings for rule
    format `v2023-05-30` ([#58](https://github.com/akamai/cli-terraform/issues/58))

## Version 1.8.0 (August 1, 2023)

### Features/Enhancements

* APPSEC
  * `export-appsec` command uses `challenge_injection_rules` resource instead of deprecated `challenge_interception_rules`
  * Added import support for `enable_pii_learning` in `akamai_appsec_advanced_settings_pii_learning` resource

### Bug fixes

* CPS
  * Fixed CN being added to an empty SANS list in `akamai_cps_third_party_enrollment`

* PAPI
  * Fixed exporting property and property include version comments in rule tree JSON
  * Fixed exporting rule tree variables with null value or description in `--schema` mode

## Version 1.7.0 (July 5, 2023)

### Features/Enhancements

* Migrated to Terraform 1.4.6 version

* PAPI
  * Added support for `export-property` command with flag `--schema` for properties in frozen formats `v2023-01-05` and `v2023-05-30`.
  * Added support for import of `akamai_property_activation` resource.
  * Added changes in `export-property` command:
    * Added support for `STAGING` and `PRODUCTION` network configurations for `akamai_property_activation` resource.
    * Removed support for `var.env` variable.
    * Added support for `auto_acknowledge_rule_warnings` default value in `akamai_property_activation` resource.

* APPSEC
  * Added support for `ukraine_geo_control_action` to `modules-security-firewall.tmpl` template for IP/Geo Firewall.

## Version 1.6.0 (May 31, 2023)

### Bug fixes

* Fix escaping of `akamai_gtm_property` `static_rr_set.rdata` field
* Export of some fields depends on the fact if the rule is default or not

### Features/Enhancements

* Migrate to Terraform 1.3.7 version
* Flag `schema` works now with `include` sub-command as well with `with-includes` flag

## Version 1.5.0 (Apr 27, 2023)

### Features/Enhancements

* APPSEC
  * Add import support for bot management resources

* EdgeKV
  * Export of `export-edgekv` uses `akamai_edgekv_group_items` resource instead of deprecated `initial_data` within `akamai_edgekv` resource

### Fixes

* GTM
  * Remove deprecated field `name` of `traffic_target` during export  ([I#374](https://github.com/akamai/terraform-provider-akamai/issues/374))

* PAPI
  * `is_secure` and `variable` fields can only be used in `default` datasource `akamai_property_rules_builder`
  * Support for `advanced_override` and `custom_override` fields in `default` datasource `akamai_property_rules_builder`
  * Fix ending newline character during export for heredoc in `akamai_property_rules_builder` datasource
  * Export `akamai_property.rule_format` as reference to `akamai_property_rules_builder`
  * Remove `certificate` and `product_id` from edgehostnames during export `akamai_edge_hostname` ([I#338](https://github.com/akamai/terraform-provider-akamai/issues/338))

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
