terraform init
terraform import akamai_edge_hostname.test-edgesuite-net ehn_2867480,test_contract,grp_12345
terraform import akamai_property.test-edgesuite-net prp_12345,test_contract,grp_12345,3
terraform import akamai_property_activation.test-edgesuite-net-staging prp_12345:STAGING
