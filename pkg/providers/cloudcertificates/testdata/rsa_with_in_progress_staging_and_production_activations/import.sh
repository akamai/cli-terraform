terraform init
terraform import akamai_cloudcertificates_lineage.rsa-with-in-progress-staging-and-production-activations 123456
terraform import akamai_cloudcertificates_upload.rsa-with-in-progress-staging-and-production-activations 123456,567890
#terraform import akamai_cloudcertificates_activation.rsa-with-in-progress-staging-and-production-activations 123456,<activation_id>[,<activation_id>]
