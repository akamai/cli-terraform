terraform init
terraform import akamai_property_include.test_include test_contract:test_group:inc_123456
terraform import akamai_property_include_activation.test_include_staging test_contract:test_group:inc_123456:STAGING
terraform import akamai_property_include_activation.test_include_production test_contract:test_group:inc_123456:PRODUCTION
