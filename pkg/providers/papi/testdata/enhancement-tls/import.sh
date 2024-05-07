terraform init
terraform import akamai_edge_hostname.test-edgekey-net ehn_2867480,test_contract,grp_12345
terraform import akamai_property.test-edgekey-net prp_12345,test_contract,grp_12345,LATEST
terraform import akamai_property_activation.test-edgekey-net-staging prp_12345:STAGING
