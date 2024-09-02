resource "privx_carrier" "carrier-test_2" {
  name            = "my_carrier"
  access_group_id = "an_access_group_id"
  subnets = [
    "0.0.0.0/0"
  ]
  enabled = true
  extender_address = [
    "0.0.0.0/0",
  ]
  web_proxy_address                 = "0.0.0.0"
  routing_prefix                    = "routing_prefix"
  web_proxy_extender_route_patterns = ["route_pattern"]
}