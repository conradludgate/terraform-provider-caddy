package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/conradludgate/terraform-provider-caddy/caddy"
)

//go:generate tfplugindocs

func main() {
	defer caddy.CloseConns()
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: caddy.Provider,
	})
}
