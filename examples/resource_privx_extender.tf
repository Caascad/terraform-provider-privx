provider "privx" {
  privx_api_base_url        = var.PRIVX_API_BASE_URL
  /*Oauth auth can be replaced by token*/
  privx_oauth_client_id     = var.PRIVX_OAUTH_CLIENT_ID
  privx_oauth_client_secret = var.PRIVX_OAUTH_CLIENT_SECRET
  privx_api_client_id       = var.PRIVX_API_CLIENT_ID
  privx_api_client_secret   = var.PRIVX_API_CLIENT_SECRET
  /*here is the token way
  privx_api_bearer_token = var.PRIVX_API_BEARER_TOKEN
  */
  privx_debug               = var.PRIVX_DEBUG
}

resource "privx_extender" "myfirstextender" {
  name = ""
  type = "EXTENDER" /*Constant*/
  web_proxy_port = 0
  web_proxy_address = ""
  routing_prefix = ""
  group_id = ""
  extender_address = [""]
  subnets = [""]
}