terraform init
terraform import akamai_edge_hostname.test-edgesuite-net ehn_2867480,test_contract,grp_12345
terraform import akamai_property.test-edgesuite-net prp_12345,test_contract,grp_12345,LATEST
terraform import akamai_property_include.test_include test_contract:test_group:inc_123456
terraform import akamai_property_include_activation.test_include_staging test_contract:test_group:inc_123456:STAGING
terraform import akamai_property_include_activation.test_include_production test_contract:test_group:inc_123456:PRODUCTION
terraform import akamai_property_include.test_include_1 test_contract:test_group:inc_78910
terraform import akamai_property_include_activation.test_include_1_staging test_contract:test_group:inc_78910:STAGING