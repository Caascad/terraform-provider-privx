resource "privx_api_client" "foo" {
  name = "test-provider-terraform"
  roles = [
    {
      id   = "dc012ffa-540c-563a-5293-51a644ce273b"
      name = "role_[Default]_[ADMIN]"
    }
  ]
}
