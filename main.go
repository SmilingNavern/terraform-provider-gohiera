package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/SmilingNavern/terraform-provider-gohiera/hiera"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hiera.Provider})
}
