terraform init
terraform import akamai_mtlstruststore_ca_set.test-ca-set-name "12345"
terraform import akamai_mtlstruststore_ca_set_activation.test-ca-set-name-staging "12345:STAGING"
#terraform import akamai_mtlstruststore_ca_set_activation.test-ca-set-name-production "12345:PRODUCTION"
