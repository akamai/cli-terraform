terraform init
terraform import akamai_property.test-edgesuite-net prp_12345,test_contract,grp_12345,LATEST
terraform import akamai_property_activation.test-edgesuite-net-staging prp_12345:STAGING
