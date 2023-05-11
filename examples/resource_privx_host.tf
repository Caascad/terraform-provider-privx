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

resource "host" "myfirsthost" {
  provider=privx
  access_group_id = ""
  external_id = ""
  instance_id = ""
  source_id = ""
  name = ""
  contact_adress = ""
  cloud_provider = ""
  cloud_provider_region = ""
  tags = [""]
  addresses = [""]
}