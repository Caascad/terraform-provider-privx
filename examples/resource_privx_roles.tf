provider "privx" {
  privx_api_base_url        = var.PRIVX_API_BASE_URL
  privx_oauth_client_id     = var.PRIVX_OAUTH_CLIENT_ID
  privx_oauth_client_secret = var.PRIVX_OAUTH_CLIENT_SECRET
  privx_api_client_id       = var.PRIVX_API_CLIENT_ID
  privx_api_client_secret   = var.PRIVX_API_CLIENT_SECRET
  privx_debug               = var.PRIVX_DEBUG
}

resource "role" "role-test" {
  provider = privx
  name = "test-role"
  comment = "role test"
  access_group_id = ""
  permissions = ["connections-view"]
}

/* List of available permissions.
vault-add
vault-manage
connections-view
connections-manage
connections-playback
connections-trail
connections-terminate
connections-manual
connections-authorize
logs-view
logs-manage
roles-view
roles-manage
users-view
users-manage
hosts-view
hosts-manage
network-targets-view
network-targets-manage
sources-view
sources-manage
sources-data-push
access-groups-manage
workflows-view
workflows-manage
workflows-requests
workflows-requests-on-behalf
requests-view
settings-view
settings-manage
role-target-resources-view
role-target-resources-manage
api-clients-manage
licenses-manage
authorized-keys-manage
*/