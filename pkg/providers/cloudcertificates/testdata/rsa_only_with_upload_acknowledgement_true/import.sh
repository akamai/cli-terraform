terraform init
terraform import akamai_cloudcertificates_lineage.rsa-only-with-upload-acknowledgement-true 123456
terraform import akamai_cloudcertificates_upload.rsa-only-with-upload-acknowledgement-true 123456,567890,true
#terraform import akamai_cloudcertificates_activation.rsa-only-with-upload-acknowledgement-true 123456,<activation_id>[,<activation_id>]
