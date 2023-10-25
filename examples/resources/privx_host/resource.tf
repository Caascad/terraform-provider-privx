provider "privx" {
}

resource "privx_host" "foo" {
  access_group_id       = "af1ac498-c3e0-4037-6dd5-456d4b2e17a2"
  external_id           = ""
  instance_id           = ""
  source_id             = ""
  common_name           = "test-dev-provider"
  contact_address       = ""
  cloud_provider        = ""
  cloud_provider_region = ""
  distinguished_name    = ""
  organization          = ""
  organizational_unit   = ""
  zone                  = ""
  host_type             = ""
  host_classification   = ""
  comment               = ""
  disabled              = "false" # BY_ADMIN | BY_LISCENCE | false
  deployable            = false
  tofu                  = false
  stand_alone_host      = false
  audit_enabled         = false
  scope                 = [""]
  tags                  = [""]
  addresses             = ["10.10.10.10"]
  certificate_template  = ""

  services = [
    {
      service                   = "SSH" # SSH | RDP | VNC | HTTP | HTTPS
      address                   = ""
      port                      = ""
      source                    = ""
      use_for_password_rotation = false
    }
  ]

  principals = [
    {
      principal                 = ""
      rotate                    = false
      use_for_password_rotation = false
      use_user_account          = false
      passphrase                = ""
      source                    = ""
      roles = [
        {
          id = ""
        }
      ]
      applications = [
        {
          name              = ""
          application       = ""
          arguments         = ""
          working_directory = ""
        }
      ]
      service_options = {
        ssh = {
          shell         = false
          file_transfer = false
          exec          = false
          tunnels       = false
          x11           = false
          other         = false
        }
        rdp = {
          file_transfer = false
          audio         = false
          clipboard     = false
          web           = false
        }
        web = {
          file_transfer = false
          audio         = false
          clipboard     = false
        }
      }
    }
  ]
}
