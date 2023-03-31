data "roles" "list_roles" {
  provider = privx
}

output "roles" {
  value=data.roles.list_roles
}