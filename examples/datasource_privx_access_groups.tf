provider "privx" {
  privx_api_base_url        = var.PRIVX_API_BASE_URL
  /*Oauth auth can be replaced by token*/
  privx_oauth_client_id     = var.PRIVX_OAUTH_CLIENT_ID
  privx_oauth_client_secret = var.PRIVX_OAUTH_CLIENT_SECRET
  privx_api_client_id       = var.PRIVX_API_CLIENT_ID
  privx_api_client_secret   = var.PRIVX_API_CLIENT_SECRET
  /*here is the token way

  */
  privx_debug               = var.PRIVX_DEBUG
}

data "access_groups" "access_groups_list" {
  provider = privx
}

output "access_groups_output" {
  value=data.access_groups.access_groups_list
}