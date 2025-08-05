terraform init
terraform import akamai_mtlstruststore_ca_set.a_funny-set_name "12345"
terraform import akamai_mtlstruststore_ca_set_activation.a_funny-set_name-staging "12345:STAGING"
terraform import akamai_mtlstruststore_ca_set_activation.a_funny-set_name-production "12345:PRODUCTION"
