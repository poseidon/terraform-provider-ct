package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/coreos/terraform-provider-ct/ct"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ct.Provider,
	})
}
