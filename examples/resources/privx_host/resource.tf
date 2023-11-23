resource "privx_host" "foo" {
  access_group_id       = "565381ce-0911-4ba8-8606-8eecd8074556"
  external_id           = ""
  instance_id           = ""
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
  tofu                  = false
  stand_alone_host      = false
  audit_enabled         = false
  scope                 = ["foo"]
  tags                  = ["foo"]
  addresses             = ["0.0.0.0"]
  ssh_host_public_keys = [
    {
      key = "" # Must be a valid format
    }
  ]

  services = [
    {
      service = "SSH" # SSH | RDP | VNC | HTTP | HTTPS
      address = "0.0.0.0"
      port    = 22
    }
  ]

  principals = [
    {
      principal        = "foo"
      passphrase       = "bar"
      use_user_account = false
      roles = [
        {
          id = "1fb15cfa-6137-4821-b60c-ffc0ba11bb86"
        }
      ]
    }
  ]
}

#resource "privx_host" "foo" {
#  access_group_id       = "565381ce-0911-4ba8-8606-8eecd8074556"
#  external_id           = ""
#  instance_id           = ""
#  common_name           = "test-dev-provider"
#  contact_address       = ""
#  cloud_provider        = ""
#  cloud_provider_region = ""
#  distinguished_name    = ""
#  organization          = ""
#  organizational_unit   = ""
#  zone                  = ""
#  host_type             = ""
#  host_classification   = ""
#  comment               = ""
#  tofu                  = false
#  stand_alone_host      = false
#  audit_enabled         = false
#  scope                 = []
#  tags                  = []
#  addresses             = []
#
#  #  certificate_template  = "" # FIXME: not implemented in privx-sdk-go v1.29
#  #  host_certificate_raw  = "" # FIXME: not implemented in privx-sdk-go v1.29
#  #  use_for_password_rotation = false # FIXME: not implemented in privx-sdk-go v1.29
#
#  ssh_host_public_keys = [
#    {
#      key = "" # Must be a valid format
#    }
#  ]
#
#  services = [
#    {
#      service = "SSH" # SSH | RDP | VNC | HTTP | HTTPS
#      address = "0.0.0.0"
#      port    = 22
#      #     use_for_password_rotation = false # FIXME: not implemented in privx-sdk-go v1.29
#    }
#  ]
#
#  principals = [
#    {
#      principal = "examplename"
#      # rotate                    = false # FIXME: not implemented in privx-sdk-go v1.29
#      # use_for_password_rotation = false # FIXME: not implemented in privx-sdk-go v1.29
#      use_user_account = false
#      passphrase       = "toto"
#      roles = [
#        {
#          id = "1fb15cfa-6137-4821-b60c-ffc0ba11bb86"
#        }
#      ]
#      #applications = [ # FIXME: not implemented in privx-sdk-go v1.29
#      #  {
#      #    name = "example-application"
#      #    # application       = "" # FIXME: not implemented in privx-sdk-go v1.29
#      #    # arguments         = "" # FIXME: not implemented in privx-sdk-go v1.29
#      #    # working_directory = "" # FIXME: not implemented in privx-sdk-go v1.29
#      #  }
#      #]
#      #service_options = { # FIXME: not implemented in privx-sdk-go v1.29
#      #  ssh = {
#      #    shell         = false
#      #    file_transfer = false
#      #    exec          = false
#      #    tunnels       = false
#      #    x11           = false
#      #    other         = false
#      #  }
#      #  rdp = {
#      #    file_transfer = false
#      #    audio         = false
#      #    clipboard     = false
#      #    web           = false
#      #  }
#      #  web = {
#      #    file_transfer = false
#      #    audio         = false
#      #    clipboard     = false
#      #  }
#      #}
#      #command_restrictions = { # FIXME: not implemented in privx-sdk-go v1.29
#      #  enabled = false
#      #  default_whitelist = {
#      #    id      = ""
#      #    name    = ""
#      #    deleted = ""
#      #  }
#      #  rshell_variant = "" # bash | posix
#      #  banner         = ""
#      #  allow_no_match = false
#      #  audit_match    = false
#      #  audit_no_match = false
#      #  whitelists = [
#      #    {
#      #      whitelist = {
#      #        id   = ""
#      #        name = ""
#      #      }
#      #      roles = [
#      #        {
#      #          id   = ""
#      #          name = ""
#      #        }
#      #      ]
#      #    }
#      #  ]
#      #}
#    }
#  ]
#  # password_rotation = { # FIXME: not implemented in privx-sdk-go v1.29
#  #   use_main_account   = true
#  #   operating_system   = "LINUX" # LINUX | WINDOWS
#  #   winrm_address      = ""
#  #   winrm_port         = 0
#  #   protocol           = "SSH" # SSH | WINRM
#  #   password_policy_id = ""
#  #   script_template_id = ""
#  # }
#}
