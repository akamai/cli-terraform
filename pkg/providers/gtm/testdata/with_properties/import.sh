terraform init
terraform import akamai_gtm_domain.test_name "test.name.akadns.net"
terraform import akamai_gtm_datacenter.TEST1 "test.name.akadns.net:123"
terraform import akamai_gtm_datacenter.TEST2 "test.name.akadns.net:124"
terraform import akamai_gtm_property.test_property1 "test.name.akadns.net:test property1"
terraform import akamai_gtm_property.test_property2 "test.name.akadns.net:test property2"