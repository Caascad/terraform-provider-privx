package main

import (
	"github.com/caascad/terraform-provider-privx/privx"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	opts := &plugin.ServeOpts{ProviderFunc: privx.Provider}
	plugin.Serve(opts)
}
