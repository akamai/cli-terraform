terraform init
terraform import akamai_cloudcertificates_lineage.ecdsa-only-with-upload 123456
terraform import akamai_cloudcertificates_upload.ecdsa-only-with-upload 123456,567890
#terraform import akamai_cloudcertificates_activation.ecdsa-only-with-upload 123456,<activation_id>[,<activation_id>]
