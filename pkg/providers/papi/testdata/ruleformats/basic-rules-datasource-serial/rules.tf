
data "akamai_property_rules_builder" "test-edgesuite-net_rule_default" {
  rules_v2023_01_05 {
    name      = "default"
    is_secure = false
    comments  = "The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings."
    uuid      = "default"
    behavior {
      origin {
        cache_key_hostname = "REQUEST_HOST_HEADER"
        compress           = true
        custom_certificate_authorities {
          can_be_ca   = true
          can_be_leaf = false
          issuer_rdns {
            c  = "US"
            cn = "GeoTrust Primary Certification Authority"
            o  = "GeoTrust Inc."
          }
          not_after                 = 1698710399000
          not_before                = 1383177600000
          pem_encoded_cert          = <<EOT
-----BEGIN CERTIFICATE-----
MIIEbjFAKEagAwIBAgIQboqQ68/wRIpyDQgF0IKlRDANBgkqhkiG9w0BAQsFADBY
MQswCQYDVQQGEwJVUzEWMBQGA1UEChMNR2VvVHJ1c3QgSW5jLjExMC8GA1UEAxMo
R2VvVHJ1c3QgUHJpbWFyeSBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eTAeFw0xMzEw
MzEwMDAwMDBaFw0yMzEwMzAyMzU5NTlaMEcxCzAJBgNVBAYTAlVTMRYwFAYDVQQK
Ew1HZW9UcnVzdCBJbmMuMSAwHgYDVQQDExdHZW9UcnVzdCBFViBTU0wgQ0EgLSBH
NDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANm0BfI4Zw8J53z1Yyrl
uV6oEa51cdlMhGetiV38KD0qsKXV1OYwCoTU5BjLhTfFRnHrHHtp22VpjDAFPgfh
bzzBC2HmOET8vIwvTnVX9ZaZfD6HHw+QS3DDPzlFOzpry7t7QFTRi0uhctIE6eBy
GpMRei/xq52cmFiuLOp3Xy8uh6+4a+Pi4j/WPeCWRN8RVWNSL/QmeMQPIE0KwGhw
FYY47rd2iKsYj081HtSMydt+PUTUNozBN7VZW4f56fHUxSi9HdzMlnLReqGnILW4
r/hupWB7K40f7vQr1mnNr8qAWCnoTAAgikkKbo6MqNEAEoS2xeKVosA7pGvwgtCW
XSUCAwEAAaOCAUMwggE/MBIGA1UdEwEB/wQIMAYBAf8CAQAwDgYDVR0PAQH/BAQD
AgEGMC8GCCsGAQUFBwEBBCMwITAfBggrBgEFBQcwAYYTaHR0cDovL2cyLnN5bWNi
LmNvbTBHBgNVHSAEQDA+MDwGBFUdIAAwNDAyBggrBgEFBQcCARYmaHR0cHM6Ly93
d3cuZ2VvdHJ1c3QuY29tL3Jlc291cmNlcy9jcHMwNAYDVR0fBC0wKzApoCegJYYj
aHR0cDovL2cxLnN5bWNiLmNvbS9HZW9UcnVzdFBDQS5jcmwwKQYDVR0RBCIwIKQe
MBwxGjAYBgNVBAMTEVN5bWFudGVjUEtJLTEtNTM4MB0GA1UdDgQWBBTez1xQt64C
HxUXqhboDbUonWpa8zAfBgNVHSMEGDAWgBQs1VBBlxWL8I82YVtK+2vZmckzkjAN
BgkqhkiG9w0BAQsFAAOCAQEAtI69B7mahew7Z70HYGHmhNHU7+sbuguCS5VktmZT
I723hN3ke40J2s+y9fHDv4eEvk6mqMLnEjkoNOCkVkRADJ+IoxXT6NNe4xwEYPtp
Nk9qfgwqKMHzqlgObM4dB8NKwJyNw3SxroLwGuH5Tim9Rt63Hfl929kPhMuSRcwc
sxj2oM9xbwwum9Its5mTg0SsFaqbLmfsT4hpBVZ7i7JDqTpsHBMzJRv9qMhXAvsc
4NG9O1ZEZcNj9Rvv7DDZ424uE+k5CCoMcvOazPYnKYTT70zHhBFlH8bjgQPbh8x4
97Wdlj5qf7wRhXp15kF9Dc/55YVpJY/HjQct+GkPy0FTAA==
-----END CERTIFICATE-----
EOT
          public_key                = "MIIBIjANBgkqhkiG9w0BFAKEAAOCAQ8AMIIBCgKCAQEA2bQF8jhnDwnnfPVjKuW5XqgRrnVx2UyEZ62JXfwoPSqwpdXU5jAKhNTkGMuFN8VGcesce2nbZWmMMAU+B+FvPMELYeY4RPy8jC9OdVf1lpl8PocfD5BLcMM/OUU7OmvLu3tAVNGLS6Fy0gTp4HIakxF6L/GrnZyYWK4s6ndfLy6Hr7hr4+LiP9Y94JZE3xFVY1Iv9CZ4xA8gTQrAaHAVhjjut3aIqxiPTzUe1IzJ2349RNQ2jME3tVlbh/np8dTFKL0d3MyWctF6oacgtbiv+G6lYHsrjR/u9CvWac2vyoBYKehMACCKSQpujoyo0QAShLbF4pWiwDuka/CC0JZdJQIDAQAB"
          public_key_algorithm      = "RSA"
          public_key_format         = "X.509"
          self_signed               = false
          serial_number             = "1.4693455585277353e+38"
          sha1_fingerprint          = "3056b343485b9d55f3e2b177a895bb0463ee3efd"
          sig_alg_name              = "SHA256withRSA"
          subject_alternative_names = []
          subject_cn                = "GeoTrust EV SSL CA - G4"
          subject_rdns {
            c  = "US"
            cn = "GeoTrust EV SSL CA - G4"
            o  = "GeoTrust Inc."
          }
          version = 3
        }
        custom_certificates {
          can_be_ca   = false
          can_be_leaf = true
          issuer_rdns {
            c  = "US"
            cn = "Amazon"
            o  = "Amazon"
            ou = "Server CA 1B"
          }
          not_after                 = 1661644799000
          not_before                = 1627516800000
          pem_encoded_cert          = <<EOT
-----BEGIN CERTIFICATE-----
MIIGCzCCBFAKEwIBAgIQA3fkBpb9nE+XoSK/oQBCozANBgkqhkiG9w0BAQsFADBG
MQswCQYDVQQGEwJVUzEPMA0GA1UEChMGQW1hem9uMRUwEwYDVQQLEwxTZXJ2ZXIg
Q0EgMUIxDzANBgNVBAMTBkFtYXpvbjAeFw0yMTA3MjkwMDAwMDBaFw0yMjA4Mjcy
MzU5NTlaMBoxGDAWBgNVBAMMDyoubXl0aGVyZXNhLmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBANYKWksTTlwgu4qtGianyiT74exN1LORJoAdZdRc
cMD1fZMaaTuh5N0fIw2OTiWPs9y2PH5hs7XlXrKm2xD7xp722n0n+J+/vAV4gPUL
wsqKth0tGDxi8az8G+pTDEf1EVEvVJSw3jpQFCO3vGMJeJqs02fNSHMevZjDS3JD
3qQa3ZL9wWAOELBbbzM57OyLBq2KlD4Uxdc19rMQzPRkzyAackcziVVNBCsXS6ph
mw9rx+IbopER2xnhcfkhcdjSNm65z3WQi367YdV/VcD7zUd9Ea64K8dqiN9QY0G5
jRJiEs6zfjvV803A6U5tY9TPY68aAHzs4lPUCEfvzz/yqqECAwEAAaOCAx8wggMb
MB8GA1UdIwQYMBaAFFmkZgZSoHuVkjyjlAcnlnRb+T3QMB0GA1UdDgQWBBT3mQuP
JXVtAKh7a86V6NZmh5YYPjBPBgNVHREESDBGgg8qLm15dGhlcmVzYS5jb22CDW15
dGhlcmVzYS5jb22CEioubXl0aGVyZXNhLmNvbS5jboIQbXl0aGVyZXNhLmNvbS5j
bjAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMC
MDsGA1UdHwQ0MDIwMKAuoCyGKmh0dHA6Ly9jcmwuc2NhMWIuYW1hem9udHJ1c3Qu
Y29tL3NjYTFiLmNybDATBgNVHSAEDDAKMAgGBmeBDAECATB1BggrBgEFBQcBAQRp
MGcwLQYIKwYBBQUHMAGGIWh0dHA6Ly9vY3NwLnNjYTFiLmFtYXpvbnRydXN0LmNv
bTA2BggrBgEFBQcwAoYqaHR0cDovL2NydC5zY2ExYi5hbWF6b250cnVzdC5jb20v
c2NhMWIuY3J0MAwGA1UdEwEB/wQCMAAwggGABgorBgEEAdZ5AgQCBIIBcASCAWwB
agB3ACl5vvCeOTkh8FZzn2Old+W+V32cYAr4+U1dJlwlXceEAAABevCFP78AAAQD
AEgwRgIhAI0dmeYj3HQrGrRjwzT1MaoIxVWWmx1pDRXrp0FnmWDuAiEAirOnsVWd
EtT1Y+lNaj5ZWYtO1bHoRG4XVnJHKaVwL0MAdwBRo7D1/QF5nFZtuDd4jwykeswb
J8v3nohCmg3+1IsF5QAAAXrwhT+1AAAEAwBIMEYCIQCt33r54fDpd/yaYN7cwazr
m1/RtIY+ysHFTGap8DCQewIhAN0rQ25+HU3ybP5n9LMNijvzr61w5Ip61d55gFE5
8PseAHYAQcjKsd8iRkoQxqE6CUKHXk4xixsD6+tLx2jwkGKWBvYAAAF68IU/bAAA
BAMARzBFAiBR/vMRwvCfKPZ7ewre7Nb8ZlXldO7tKwpqKtre9R4CFQIhAK44D9Oj
Kk/CMsKRpFFkdgDQpDWx1w71QwJTWFjRcFpCMA0GCSqGSIb3DQEBCwUAA4IBAQCq
VLjSaLp2lPm5HXDJBo6Tc/qcavM5lpeLsxN3Q7W9NOgn7NdT6gnGmWASlnMmtMzy
w98VsFazjs+lJeXVaRCD7NMTD5ijQ9wixYVGvdBkDVqGqImhW3xI0SuDsYlJSvFt
RcBhnop9HxADKIAP0kRMgLVSouQNyIYifnqCUPzw0Dz8zeN5fuFjVkSzZQe5CFek
YH6nNKTANBa6SU3DgoTIw5cebYdMsiVPZ1fwOoMxTGflz9WEO2kHb/+mYcnzuNuh
CNMsdiMUf48dIhpTs5Oe5VBuWfdeC6IVhh5cVwu5yGhG9zuPq9/IHW+eJBG/6Nrv
PlTQzLLqEPFBItdmQwe3
-----END CERTIFICATE-----
EOT
          public_key                = "MIIBIjANBgkqFAKE9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1gpaSxNOXCC7iq0aJqfKJPvh7E3Us5EmgB1l1FxwwPV9kxppO6Hk3R8jDY5OJY+z3LY8fmGzteVesqbbEPvGnvbafSf4n7+8BXiA9QvCyoq2HS0YPGLxrPwb6lMMR/URUS9UlLDeOlAUI7e8Ywl4mqzTZ81Icx69mMNLckPepBrdkv3BYA4QsFtvMzns7IsGrYqUPhTF1zX2sxDM9GTPIBpyRzOJVU0EKxdLqmGbD2vH4huikRHbGeFx+SFx2NI2brnPdZCLfrth1X9VwPvNR30Rrrgrx2qI31BjQbmNEmISzrN+O9XzTcDpTm1j1M9jrxoAfOziU9QIR+/PP/KqoQIDAQAB"
          public_key_algorithm      = "RSA"
          public_key_format         = "X.509"
          self_signed               = false
          serial_number             = "4610192225008347635903659133613654691"
          sha1_fingerprint          = "4d427593493f5c689a709d7efa5a429946a91c91"
          sig_alg_name              = "SHA256withRSA"
          subject_alternative_names = ["*.foo.com", "foo.com.cn", ]
          subject_cn                = "*.foo.com"
          subject_rdns {
            cn = "*.foo.com"
          }
          version = 3
        }
        custom_valid_cn_values           = ["{{Origin Hostname}}", "{{Forward Host Header}}", ]
        enable_true_client_ip            = true
        forward_host_header              = "REQUEST_HOST_HEADER"
        hostname                         = "foo.amazonaws.com"
        http_port                        = 80
        https_port                       = 443
        origin_certificate               = ""
        origin_certs_to_honor            = "COMBO"
        origin_sni                       = true
        origin_type                      = "CUSTOMER"
        ports                            = ""
        standard_certificate_authorities = ["akamai-permissive", ]
        true_client_ip_client_setting    = false
        true_client_ip_header            = "True-Client-IP"
        verification_mode                = "CUSTOM"
      }
    }
    behavior {
      cp_code {
        value {
          id = 452416
        }
      }
    }
  }
}
