provider "privx" {
  privx_api_base_url        = var.PRIVX_API_BASE_URL
  privx_oauth_client_id     = var.PRIVX_OAUTH_CLIENT_ID
  privx_oauth_client_secret = var.PRIVX_OAUTH_CLIENT_SECRET
  privx_api_client_id       = var.PRIVX_API_CLIENT_ID
  privx_api_client_secret   = var.PRIVX_API_CLIENT_SECRET
  privx_debug               = var.PRIVX_DEBUG
}

resource "extender" "extender-test" {
  provider = privx
  name = "extendertestprovidertoto"
  type = "EXTENDER"
  web_proxy_port = 0
  web_proxy_address = ""
  routing_prefix = ""
  group_id = ""
}