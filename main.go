package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/poseidon/terraform-provider-ct/ct"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ct.Provider,
	})
}
