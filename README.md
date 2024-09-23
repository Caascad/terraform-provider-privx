# Terraform Provider PrivX (WIP)

This repository is an unofficial terraform provider for PrivX made by the Caascad team.
It uses [privx-sdk-go](https://github.com/SSHcom/privx-sdk-go).

## Using the provider

See examples (`examples/`) and generated documentation (`docs/`),

## Development

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

Optionnal
- [Nix](https://github.com/NixOS/nix) >= 2.15.0
- [Golangci-lint](https://github.com/golangci/golangci-lint) >= 1.53.3
- [Pre-commit](https://github.com/pre-commit/pre-commit) >= 3.3.3

### Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
go mod vendor
```

Then commit the changes to `vendor/`, `go.mod` and `go.sum`.

### Developing the Provider
_This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).
You can use [Nix](https://github.com/NixOS/nix) to easily get all dependencies:
```shell
nix-shell shell.nix
```

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.
In order to use the built provider for development, add the following `$HOME/.terraformrc`
```
provider_installation {

  dev_overrides {
      "privx" = "$GOPATH/bin/"
  }
}
```

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
