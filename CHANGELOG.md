# RELEASE NOTES

## 2.0.0 (Jan 5, 2025)

### BREAKING CHANGES:

* Removed deprecated commands:
  * `create-domain`
  * `create-zone`
  * `create-appsec`
  * `create-property`
  * `create-cloudlets-policy`
  * `create-edgekv`
  * `create-edgeworker`
  * `create-iam`
  * `create-imaging`
  * `create-cps`

* Removed the `include` subcommand from the `export-property` command.
* Removed the `with-includes` flag from the `export-property` command.
* Removed the `schema` flag from the `export-property`, `export-property-include` and `export-imaging` commands.

### FEATURES/ENHANCEMENTS:

* General
  * Migrated to Go `1.22`.
  * Improved code by resolving issues reported by linter.
  * Updated vulnerable dependencies.
  * Logic responsible for excluding endpoints from retries is now configurable with the `AKAMAI_RETRY_EXCLUDED_ENDPOINTS` environment variable.

* IAM
  * Added the `--only` option to the `group`, `role`, and `user` subcommands of the `export-iam` command. This option allows exporting only specific information.

* Logging
  * Changed logger from `apex` to `slog`.
    * Log output has not been changed.

* PAPI
  * Added support for the new rule format `v2025-01-13`.
  * Adjusted exported fields for current schema definition for the rule format `v2024-10-21` inside the `gov_cloud` behaviour.
  * Added a property name in the error message.
  * Introduced the `split-depth` flag for the `export-property` and `export-property-include` command. When used, each rule up to a specified nesting level will be generated in its own `.tf` file. Rules with higher nesting levels will be placed in a file of their closest ancestor. All rules will be generated in a dedicated `rules` module.

### BUG FIXES:

* General
  * Fixed a problem with invisible output in the light background by converting all colors to the monochromatic representation.

* AppSec
  *  Fixed issues with AAP/WAP terraform export
    * Removed `hostnames` block from `variables.tf` file
    * Modified `imports.tmpl` to import correct resources for AAP/WAP accounts
    * Renamed all references from `akamai_appsec_configuration.config` to `data.akamai_appsec_configuration.config` for WAP/AAP accounts

##  1.19.0 (Nov 21, 2024)

### FEATURES/ENHANCEMENTS:

* General
  * Added retryable logic for all GET requests to the API.
    This behavior can be disabled using the `AKAMAI_RETRY_DISABLED` environment variable.
    It can be fine-tuned using these environment variables:
    * `AKAMAI_RETRY_MAX` - The maximum retires number of API requests, default is 10.
    * `AKAMAI_RETRY_WAIT_MIN` - The minimum wait time in seconds between API requests retries, default is 1 sec.
    * `AKAMAI_RETRY_WAIT_MAX` - The maximum wait time in seconds between API requests retries, default is 30 sec.

* AppSec
  *  Added support to export the `akamai_botman_content_protection_rule` resource for the specified policy.
  *  Added support to export the `akamai_botman_content_protection_rule_sequence` resource for the specified policy.
  *  Added support to export the `akamai_botman_content_protection_javascript_injection_rule` resource for the specified policy.

* Cloud Access Manager
  * Added the `group_id` and `contract_id` flags for `export-cloudaccess` which allows exporting the `akamai_cloudaccess_key` resource with the specified group and contract IDs.

* DNS
  * Added the new `outbound_zone_transfer` field to the `akamai_dns_zone` resource.

* PAPI
  * Added support for the new rule format `v2024-10-21`.
  * Added the `product_id` attribute to the exported `akamai_property_include` resource.

### BUG FIXES:

* AppSec
  * Fixed improper resource generation for AAP/AAP accounts in the `export-appsec` output:
    * Removed the `akamai_appsec_configuration` resource and added the `akamai_appsec_configuration` data source.
    * Removed the `hostnames` variable from the `appsec-variables.tf` file.
    * Removed the `hostnames` field from the `security` module.


* Cloud Access Manager
  *  Marked the `cloud_secret_access_key` field as sensitive in a template for the `akamai_cloudaccess_key` resource and moved its definition to the `variables.tf` file. ([I#580](https://github.com/akamai/terraform-provider-akamai/issues/580))

* PAPI
  * Fixed a missing child file when using uppercase letters for a second rule with lowercase letters ([#78](https://github.com/akamai/cli-terraform/issues/78)).

### DEPRECATIONS:

* Excluded the deprecated `akamai_appsec_wap_selected_hostnames` resource from the `export-appsec` command. Instead, use the `akamai_appsec_aap_selected_hostnames` resource to export.

## 1.18.0 (Oct 10, 2024)

### FEATURES/ENHANCEMENTS:

* AppSec
  * The `akamai_appsec_match_target` resource is created only if a target product is not an AAP account.
  * The `akamai_appsec_wap_selected_hostnames` resource can be exported for AAP accounts.

* Identity and Access Management (IAM)
  * Added support for generating the `enable_mfa` and `user_notifications` attributes when exporting the `akamai_iam_user` configuration.
  * Added the `allowlist` subcommand to the `export-iam` command that exports terraform configuration files for an account's IP allowlist and CIDR blocks.
  * Modified the `all` subcommand for the `export-iam` command to export an account's IP allowlist and CIDR blocks details.

### BUG FIXES:

* PAPI
  * Fixed an issue with property export where the hostname `cnameTo` begins with a number or contains a space.

## 1.17.0 (Sep 04, 2024)

### FEATURES/ENHANCEMENTS:

* AppSec
  * The `request_body_inspection_limit_override` field is added to the `akamai_appsec_advanced_settings_request_body` resource.

* PAPI
  * Added support for the new rule format `v2024-08-13`.

### BUG FIXES:

* Identity and Access Management (IAM)
  * Fixed handling of a new line character in fields from the exported `akamai_iam_user` resource.

## 1.16.0 (Jul 16, 2024)

### FEATURES/ENHANCEMENTS:

* Migrated go version to 1.21.12 for builds.

* [IMPORTANT] Cloud Access Manager
  * Added the `export-cloudaccess` command which allows exporting the `akamai_cloudaccess_key` resource.

* PAPI
  * If an edge hostname uses a custom TTL, it is exported in the `akamai_edge_hostname` resource.
  * Added support for the new rule format `v2024-05-31`.

## 1.15.0 (May 28, 2024)

### FEATURES/ENHANCEMENTS:

* General
  * Updated various dependencies.

* Cloudlets
  * Added import support for the `akamai_cloudlets_application_load_balancer_activation` resource.

* PAPI
  * Added export of the `certificate` for the `akamai_edge_hostname` resource ([I#338](https://github.com/akamai/terraform-provider-akamai/issues/338))

## 1.14.0 (Apr 23, 2024)

### FEATURES/ENHANCEMENTS:

* General
  * Updated various dependencies.

* Image and Video Manager
  * Added handling of `SmartCrop` transformation when exporting an image with the `policy-as-hcl` flag.

## 1.13.0 (Mar 26, 2024)

### FEATURES/ENHANCEMENTS:

* General
  * Updated minimal required terraform version to 1.0
  * Migrated to go 1.21
* AppSec
  *  Added support to export the `akamai_appsec_penalty_box_conditions` resource for the specified policy.
  *  Added support to export the `akamai_appsec_eval_penalty_box_conditions` resource for the specified policy.
* Cloudlets
  * Changed export for the `akamai_cloudlets_audience_segmentation_match_rule` resource to generate empty `forward_settings` when both `origin_id` and `path_and_qs` are empty and `use_incoming_query_string` is false.
* GTM
  * Added support for exporting fields:
    * `sign_and_serve`, `sign_and_serve_algorithm` for the `akamai_gtm_domain` resource.
    * `http_method`, `http_request_body`, `pre_2023_security_posture`, `alternate_ca_certificate` inside `liveness_test` in the `akamai_gtm_property` resource.
* PAPI
  * The `export-property` command with the `--rules-as-hcl` flag now supports exporting properties in frozen format `v2024-02-12`.
  * The `export-property-include` command with the `--rules-as-hcl` flag now supports exporting includes in frozen format `v2024-02-12`.

### BUG FIXES:

* PAPI
  * Fixed an issue that empty `custom_certificate_authorities` or `custom_certificates` where not generated during `export-property` with the `rules-as-hcl` flag.

### DEPRECATIONS:

* AppSec
  * Excluded the deprecated `akamai_appsec_selected_hostnames` resource from the `export-appsec` command. Instead, use the `akamai_appsec_configuration` resource to export.

## 1.12.0 (Feb 19, 2024)

### FEATURES/ENHANCEMENTS:

* Cloudlets
  * Modified export of active policies for cloudlets to generate the `akamai_cloudlets_policy_activation` resource entry in the `import.sh` script. When activating both networks, only production will be exported.
  * Added support to export shared (V3) policies.
  * Added support for exporting policies without any version.
* EdgeWorkers
  * Added support for generating the `note` field when exporting EdgeWorker configuration.
* IVM
  * Added support for generating `serve_stale_duration`, `allow_pristine_on_downsize` and `prefer_modern_formats` when exporting using the `--policy-as-hcl` flag.
* PAPI
  * Introduced two variables `activate_latest_on_staging` and `activate_latest_on_production` in exported configuration for property activation or include activation (when exporting includes alone), to drive which version to use for activation.
  * When there is no activation for given network, export activation is commented out.
  * Added a new export `export-property-include` command as a replacement for `export-property`.`include` subcommand. It'll generate `include` configuration without related properties
  * Deprecated the `include` subcommand available for `export-property`.
  * Deprecated the `--with-includes` flag available for `export-property`
  * The `export-property` and `export-property-include` commands with the `--rules-as-hcl` flag now support exporting properties in frozen format `v2024-01-09`.
  * Added support for the `export-property` command with the `--akamai-property-bootstrap` flag to export a property using the `akamai_property_bootstrap` resource. This option is false by default.

### BUG FIXES:

* AppSec
  * Fixed an issue where advanced exceptions were not generated for Rules and Risk Groups ([#61](https://github.com/akamai/cli-terraform/issues/61)).

## 1.11.0 (Dec 7, 2023)

### FEATURES/ENHANCEMENTS:

* AppSec
  * Added support for `asn_network_lists` to the `akamai_appsec_ip_geo` resource for IP/Geo Firewall.
* Deprecated the `--schema` flag and replaced with:
  * `--policy-as-hcl` (for export-imaging command)
  * `--rules-as-hcl`  (for export-property command)
* Cloudlets
  * Added the `origin_description` field export in the `akamai_cloudlets_application_load_balancer` resource.
* PAPI
  * The `export-property` command with the `--rules-as-hcl` flag now supports exporting properties in frozen format `v2023-10-30`.

### BUG FIXES:

* DNS (export-zone)
  * Fixed the `target` field string escaping in the `akamai_dns_record` resource.
  * Changed a provider version requirement from `~> 1.6.1` to `>= 1.6.1`.
* PAPI
  * Fixed an error with exporting a schema containing `serialNumber`.

## 1.10.0 (Oct 31, 2023)

### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Client Lists
  * Added the `export-clientlist` command which allows exporting the `akamai_clientlist_list` and `akamai_clientlist_activation` resources.
* Cloudlets
  * Added the `matches_always` field to the `akamai_cloudlets_edge_redirector_match_rule` export template.
* PAPI
  * Added support for the new rule format `v2023-09-20`.

### BUG FIXES:

* Fixed an issue with generating a multiline text for:
  * The `description` variable in AppSec configuration
  * The `comments` and `location.comments` fields in the `akamai_cloudwrapper_configuration` resource
  * The `comment` field in the `akamai_dns_zone`
  * The `comment` field in `akamai_gtm_domain` resource
  * The `comments` field in the `akamai_gtm_property` resource
  * The `description` field in the `akamai_gtm_resource` resource
  * The `note` field in the `akamai_property_activation` and `akamai_property_include_activation` resources

## 1.9.1 (Sep 26, 2023)

### BUG FIXES:

* CPS
  * Fixed s nil pointer evaluating *cps.DNSNameSettings.CloneDNSNames ([#52](https://github.com/akamai/cli-terraform/issues/52)).

* Identity and Access Management (IAM)
  * Fixed newline escaping in the `description` field after exporting a role.

* PAPI
  * Add missing fields to `akamai_property_builder` for the `origin` and `siteShield` behaviors.

## 1.9.0 (Aug 29, 2023)

### FEATURES/ENHANCEMENTS:

* [IMPORTANT] CloudWrapper
  * Added support for the `export-cloudwrapper` command which allows exporting the `akamai_cloudwrapper_configuration` and `akamai_cloudwrapper_activation` resources.

* AppSec
  * Added import support for a custom client sequence.

### BUG FIXES:

* Image and Video Manager
  * Added the description for the `--schema` flag for the `export-imaging` command in `README.md` ([#56](https://github.com/akamai/cli-terraform/issues/56)).

* PAPI
  * Fixed the `export-property` command to export the `akamai_property_activation` resource attributes for the latest active version.
  * Fixed the `export-property` command to use `group_id` and `contract_id` as terraform variables, instead of data sources, which
  produced inconsistencies ([I#374](https://github.com/akamai/terraform-provider-akamai/issues/426)).
  * The `logStreamName` field from the `datastream` behavior has changed from a string to an array of strings for the rule format `v2023-05-30` ([#58](https://github.com/akamai/cli-terraform/issues/58)).

## 1.8.0 (Aug 1, 2023)

### FEATURES/ENHANCEMENTS:

* AppSec
  * The `export-appsec` command uses now the `challenge_injection_rules` resource instead of the deprecated `challenge_interception_rules` resource.
  * Added import support for `enable_pii_learning` in the `akamai_appsec_advanced_settings_pii_learning` resource.

### BUG FIXES:

* CPS
  * Fixed CN being added to an empty SANS list in the `akamai_cps_third_party_enrollment` resource.

* PAPI
  * Fixed exporting a property and property include version comments in a rule tree JSON.
  * Fixed exporting rule tree variables with a null value or description in the `--schema` mode.

## 1.7.0 (Jul 5, 2023)

### FEATURES/ENHANCEMENTS:

* Migrated to Terraform 1.4.6 version.

* PAPI
  * Added support for the `export-property` command with the `--schema` flag for properties in frozen formats `v2023-01-05` and `v2023-05-30`.
  * Added support for importing the `akamai_property_activation` resource.
  * Added changes in the `export-property` command:
    * Added support for the `STAGING` and `PRODUCTION` network configurations for the `akamai_property_activation` resource.
    * Removed support for the `var.env` variable.
    * Added support for the `auto_acknowledge_rule_warnings` default value in the `akamai_property_activation` resource.

* AppSec
  * Added support for `ukraine_geo_control_action` to the `modules-security-firewall.tmpl` template for IP/Geo Firewall.

## 1.6.0 (May 31, 2023)

### FEATURES/ENHANCEMENTS:

* Migrated to Terraform 1.3.7 version.
* The `schema` flag works now with `include` subcommand as well as with the `with-includes` flag.

### BUG FIXES:

* Fixed escaping of the `akamai_gtm_property` `static_rr_set.rdata` field.
* Export of some fields depends on the fact if the rule is default or not.

## 1.5.0 (Apr 27, 2023)

### FEATURES/ENHANCEMENTS:

* AppSec
  * Added import support for bot management resources.

* EdgeKV
  * Export of `export-edgekv` uses the `akamai_edgekv_group_items` resource instead of the deprecated `initial_data` within the `akamai_edgekv` resource.

### BUG FIXES:

* GTM
  * Removed the deprecated `name` field of `traffic_target` during export  ([I#374](https://github.com/akamai/terraform-provider-akamai/issues/374)).

* PAPI
  * The `is_secure` and `variable` fields can only be used in `default` data source `akamai_property_rules_builder`.
  * Added support for the `advanced_override` and `custom_override` fields in `default` data source `akamai_property_rules_builder`.
  * Fixed ending a newline character during export for heredoc in the `akamai_property_rules_builder` data source.
  * Export `akamai_property.rule_format` as a reference to the `akamai_property_rules_builder` data source.
  * Removed `certificate` and `product_id` from edge hostnames when exporting `akamai_edge_hostname` ([I#338](https://github.com/akamai/terraform-provider-akamai/issues/338)).

## 1.4.0 (Mar 30, 2023)

### FEATURES/ENHANCEMENTS:

* AppSec
  * Added support for exporting the `akamai_appsec_advanced_settings_request_body` resource.

* PAPI
  * The new `--schema` flag available with the `export-property` command to export a rule tree as the `akamai_property_rules_builder` data source (Beta).

### BUG FIXES:

* PAPI
  * Fixed property export with an empty `EdgeHostnameID` ([I#41](https://github.com/akamai/cli-terraform/issues/41)).
  * Changed the exported attribute value in configuration of the `akamai_contract` data source from the deprecated `name`
    to `group_name` when exporting a property.
  * Commented out the `akamai_property_activation` resource if there is no currently active version.

## 1.3.1 (Feb 2, 2023)

### FEATURES/ENHANCEMENTS:

* General
  * Migrate to go 1.18
  * Added badges to readme and improved code quality based on golangci-lint.
* CPS
  * Added `preferred_trust_chain` in the `csr` set attribute.

## 1.3.0 (Dec 15, 2022)

### FEATURES/ENHANCEMENTS:

* PAPI
  * Added the new `export-property include` subcommand to export the `akamai_property_include` resource with the accompanying resources and data sources:
  `akamai_property_include_activation`, `akamai_property_include_parents`, and `akamai_property_rules_template`
  * Added the new `--with-includes` flag available with the `export-property` command to export resources and data sources for the Property Manager Includes that are referenced by the Property which is being exported.

## 1.2.0 (Dec 1, 2022)

### FEATURES/ENHANCEMENTS:

* CPS
  * Added the new `export-cps` command to export a DV enrollment (`akamai_cps_dv_enrollment`) or third-party enrollment with the accompanying resources and data source: (`akamai_cps_third_party_enrollment`,`akamai_cps_csr`, and `akamai_cps_upload_certificate`).

## 1.1.1 (Oct 27, 2022)

### BUG FIXES:

* GTM
  * Fixed exporting a GTM property with a default data center ([I#31](https://github.com/akamai/cli-terraform/issues/31)).

* PAPI
  * Fixed exporting the `cert_provisioning_type` field ([I#15](https://github.com/akamai/cli-terraform/issues/15)).

## 1.1.0 (Oct 10, 2022)

### FEATURES/ENHANCEMENTS:

* AppSec
  * Added import support for `malware_policy` and `malware_policy_action`.

### BUG FIXES:

* General
  * Resolved all the tflint warnings which were introduced with the tflint version `v0.40.0`.

* AppSec
  * Fixed an incorrect policy ID for malware protection.
  * Fix a drift on match targets and malware protection resources.

* PAPI
  * Fixed ignoring a property version for property-snippets when exporting a property.

## 1.0.0 (Aug 3, 2022)

### FEATURES/ENHANCEMENTS:

* GTM
  * Improved formatting of output configurations.
* DNS
  * Added support for additional default data centers.
  * Improved formatting of output configurations.

### BUG FIXES

* General
  * Fixed the default flag values in the help output.
* Identity and Access Management (IAM)
  * Fixed IAM role export failures with a broken user.

### DEPRECATIONS:

* [IMPORTANT] General
  * The `create-*` command names are now deprecated. Use the `export-*` command instead.

## 0.9.0 (Jul 07, 2022)

### FEATURES/ENHANCEMENTS:

* Identity and Access Management  (IAM)
  * Added the new `create-iam` command to export users, roles, and/or groups.
* AppSec
  * Added the new `create-appsec` command to create the Application Security resource.

## 0.8.0 (Jun 02, 2022)

### FEATURES/ENHANCEMENTS:

* General
  * Added `arm64` support (Apple M1).
* Image and Video Manager
  * Added the new `create-imaging` command to import image and video policies.
* PAPI
  * Added the optional `version` flag to the `create-property` command ([I#8](https://github.com/akamai/cli-terraform/issues/8)).

## 0.7.1 (May 11, 2022)

### FEATURES/ENHANCEMENTS:

* General
  * Added support for the `--acountkey` flag.
  * Added logging improvements.

### BUG FIXES:

* General
  * `./cli-terraform --help` should return a zero status.

* PAPI
  * Normalized rule names in the `create-property` command.

## 0.7.0 (Mar 31, 2022)

### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Image and Video Manager
  * Added support for importing the existing Video and Image Policies.
  * Added support for importing the existing Policy Set.

* Cloudlets
  * Added support for importing existing Cloudlets match rules for the Request Control with a related data source.

* General
  * Updated `urfave/cli` to v2.
    * **BREAKING CHANGE:** Now flags must come before args.
  * Updated the `README` file to reflect changes in command line flag ordering.

### BUG FIXES:

* DNS
  * Corrected handling of embedded double quotes in TXT record sets.

## 0.6.0 (Mar 3, 2022)

### FEATURES/ENHANCEMENTS:

* [IMPORTANT] EdgeWorkers and EdgeKV
  * EdgeWorkers
    * Added support for importing the existing EdgeWorker configuration with related resources.
  * EdgeKV
    * Added support for importing the existing EdgeKV configuration with a related resource.

* Cloudlets
  * Added support for importing the existing Cloudlets match rules for the Audience Segmentation with a related data source.

## 0.5.1 (Feb 10, 2022)

### FEATURES/ENHANCEMENTS:

* Cloudlets
  * Added a `default` value to the `config_section` variable in `create-cloudlets-policy` when none is specified.
  * Set the latest application load balancer version in the `akamai_cloudlets_application_load_balancer_activation` resource.

## 0.5.0 (Jan 27, 2022)

### FEATURES/ENHANCEMENTS:

* Cloudlets
  * Added support for importing the existing Cloudlets match rules for the Visitor Prioritization with a related data source.
  * Added support for importing the existing Cloudlets match rules for the Continuous Deployment/Phased Release with a related data source.
  * Added support for importing the existing Cloudlets match rules for the Forward Rewrite with a related data source.
  * Added support for importing the existing Cloudlets match rules for the API Prioritization with a related data source.

## 0.4.0 (Dec 6, 2021)

### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Cloudlets
  * Added support for importing the existing Cloudlets policy and policy versions configuration with a related resource.
  * Added support for importing the existing Cloudlets application load balancer configuration with a related resource.
  * Added support for importing the existing Cloudlets match rules for an application load balancer with a related data source.
  * Added support for importing the existing Cloudlets match rules for an edge redirector with a related data source.

## 0.3.1 (Nov 9, 2021)

### BUG FIXES:

* General
  * Updated binaries URL to fix a binary installation failure ([#6](https://github.com/akamai/cli-terraform/issues/6)).
  * Updated dependencies to fix an issue under MacOS Big Sur.

* PAPI
  * Add the `is_secure` field into a rule tree structure ([#9](https://github.com/akamai/cli-terraform/issues/9)).

## 0.3.0 (Aug 4, 2021)

### FEATURES/ENHANCEMENTS:

* General
  * Added `account-key` as an alias for the `accountkey` argument.

* DNS
  * Changed the edgerc section default to `default`.
  * Updated a generated Terraform config file header.
  * Populated the `contractId` variable when generating the `dnsvars.tf` file.

* GTM
  * Changed the edgerc section default to `default`.
  * Updated a generated Terraform config file header.

* PAPI
  * Removed deprecated CPCode support.

## 0.2.0 (May 12, 2021)

### FEATURES/ENHANCEMENTS:

* PAPI
  * Updated property resources to the latest syntax.

### BUG FIXES:

* DNS
  * Corrected AKAMAICDN record resource generation.
  * Corrected SVCB/HTTPS record resource generation.

* GTM
  * Do not include a Datacenter ID when generating the Datacenter resource.
  * Do not include a Traffic Target when generating the Static property resource.

## 0.1.0 (May 15, 2020)

### FEATURES/ENHANCEMENTS:

Initial release

* DNS
  * Added support for importing the existing DNS zones and related resources.

* GTM
  * Added support for importing the existing GTM domains and related resources.

* PAPI
  * Added support for importing the existing PAPI properties and related resources.
