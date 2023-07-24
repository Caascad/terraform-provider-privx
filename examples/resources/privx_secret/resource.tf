resource "privx_secret" "foo" {
  name = "bar"
  data = jsonencode({})
  write_roles = [
    {
      id   = "dc012ffa-540c-563a-5293-51a644ce273b"
      name = "role_[Default]_[ADMIN]"
    }
  ]
  read_roles = [
    {
      id   = "dc012ffa-540c-563a-5293-51a644ce273b"
      name = "role_[Default]_[ADMIN]"
    }
  ]
}
