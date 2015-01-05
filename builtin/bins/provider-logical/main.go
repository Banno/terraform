package main

import (
	"github.com/hashicorp/terraform/builtin/providers/logical"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: logical.Provider,
	})
}
