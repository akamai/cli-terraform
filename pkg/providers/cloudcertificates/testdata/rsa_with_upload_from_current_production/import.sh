terraform init
terraform import akamai_cloudcertificates_lineage.rsa-with-upload-from-current-production 123456
terraform import akamai_cloudcertificates_upload.rsa-with-upload-from-current-production 123456,123401
terraform import akamai_cloudcertificates_activation.rsa-with-upload-from-current-production 123456,1001,1002
