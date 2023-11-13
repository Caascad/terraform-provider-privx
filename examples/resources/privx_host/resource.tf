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

  #  certificate_template  = "" # FIXME: not implemented in privx-sdk-go v1.29
  #  host_certificate_raw  = "" # FIXME: not implemented in privx-sdk-go v1.29
  #  use_for_password_rotation = false # FIXME: not implemented in privx-sdk-go v1.29

  ssh_host_public_keys = [
    {
      key = ""
    }
  ]

  services = [
    {
      service = "SSH" # SSH | RDP | VNC | HTTP | HTTPS
      address = ""
      port    = ""
      source  = ""
      #     use_for_password_rotation = false # FIXME: not implemented in privx-sdk-go v1.29
    }
  ]

  principals = [
    {
      principal = ""
      # rotate                    = false # FIXME: not implemented in privx-sdk-go v1.29
      # use_for_password_rotation = false # FIXME: not implemented in privx-sdk-go v1.29
      use_user_account = false
      passphrase       = ""
      source           = ""
      roles = [
        {
          id = ""
        }
      ]
      applications = [
        {
          name = ""
          # application       = "" # FIXME: not implemented in privx-sdk-go v1.29
          # arguments         = "" # FIXME: not implemented in privx-sdk-go v1.29
          # working_directory = "" # FIXME: not implemented in privx-sdk-go v1.29
        }
      ]
      # service_options = { # FIXME: not implemented in privx-sdk-go v1.29
      #   ssh = {
      #     shell         = false
      #     file_transfer = false
      #     exec          = false
      #     tunnels       = false
      #     x11           = false
      #     other         = false
      #   }
      #   rdp = {
      #     file_transfer = false
      #     audio         = false
      #     clipboard     = false
      #     web           = false
      #   }
      #   web = {
      #     file_transfer = false
      #     audio         = false
      #     clipboard     = false
      #   }
      # }
      # command_restrictions = {
      #   enabled = false
      #   default_whitelist = {
      #     id      = ""
      #     name    = ""
      #     deleted = ""
      #   }
      #   rshell_variant = "" # bash | posix
      #   banner         = ""
      #   allow_no_match = false
      #   audit_match    = false
      #   audit_no_match = false
      #   whitelists = [
      #     {
      #       whitelist = {
      #         id   = ""
      #         name = ""
      #       }
      #       roles = [
      #         {
      #           id   = ""
      #           name = ""
      #         }
      #       ]
      #     }
      #   ]
      # }
    }
  ]
  # password_rotation = { # FIXME: not implemented in privx-sdk-go v1.29
  #   use_main_account   = true
  #   operating_system   = "LINUX" # LINUX | WINDOWS
  #   winrm_address      = ""
  #   winrm_port         = ""
  #   protocol           = "SSH" # SSH | WINRM
  #   password_policy_id = ""
  #   script_template_id = ""
  # }
}
