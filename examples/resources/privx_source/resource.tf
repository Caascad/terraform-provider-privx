resource "privx_source" "foo" {
  name    = "test-dev-provider"
  comment = ""
  tags    = ["tot"]
  enabled = true
  oidc_connection = {
    address             = "10.10.10.10"
    issuer              = "http://foo.com"
    button_title        = "http://foo.com/button"
    client_id           = "bar"
    client_secret       = "foobar"
    tags_attribute_name = "foo"
    enabled             = true
  }
}
