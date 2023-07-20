//provider "privx" {
//}

resource "privx_host" "foo" {
  access_group_id       = ""
  external_id           = ""
  instance_id           = ""
  source_id             = ""
  common_name           = "test-dev-provider"
  contact_address       = ""
  cloud_provider_id     = ""
  cloud_provider_region = ""
  tags                  = [""]
  addresses             = ["10.10.10.10"]
}
