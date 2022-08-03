
resource "akamai_dns_record" "zoneName_someName_someType" {
  zone       = local.zone
  hardware   = "INTEL-386"
  software   = "Unix"
  name       = "someName"
  recordtype = "someType"
  ttl        = 1000
}
