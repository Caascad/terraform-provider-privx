package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"privx": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PRIVX_API_BASE_URL"); v == "" {
		t.Fatal("PRIVX_API_BASE_URL must be set for acceptance tests")
	}
	if v := os.Getenv("PRIVX_OAUTH_CLIENT_ID"); v == "" {
		t.Fatal("PRIVX_OAUTH_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("PRIVX_OAUTH_CLIENT_SECRET"); v == "" {
		t.Fatal("PRIVX_OAUTH_CLIENT_SECRET must be set for acceptance tests")
	}
	if v := os.Getenv("PRIVX_API_CLIENT_ID"); v == "" {
		t.Fatal("PRIVX_API_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("PRIVX_API_CLIENT_SECRET"); v == "" {
		t.Fatal("PRIVX_API_CLIENT_SECRET must be set for acceptance tests")
	}
}
