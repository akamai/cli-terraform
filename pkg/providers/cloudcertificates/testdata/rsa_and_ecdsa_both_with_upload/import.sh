terraform init
terraform import akamai_cloudcertificates_lineage.rsa-and-ecdsa-both-with-upload 123456
terraform import akamai_cloudcertificates_upload.rsa-and-ecdsa-both-with-upload 123456,567890
#terraform import akamai_cloudcertificates_activation.rsa-and-ecdsa-both-with-upload 123456,<activation_id>[,<activation_id>]
