terraform init
terraform import akamai_gtm_domain.test_name "test.name.akadns.net"
terraform import akamai_gtm_datacenter.TEST1 "test.name.akadns.net:123"
terraform import akamai_gtm_datacenter.TEST2 "test.name.akadns.net:124"
terraform import akamai_gtm_cidrmap.test_cidrmap "test.name.akadns.net:test_cidrmap"
terraform import akamai_gtm_geomap.test_geomap "test.name.akadns.net:test_geomap"
terraform import akamai_gtm_asmap.test_asmap "test.name.akadns.net:test_asmap"