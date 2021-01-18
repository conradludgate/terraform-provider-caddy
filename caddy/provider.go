package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider for caddy
func Provider() *schema.Provider {
	return tfutils.ProviderBuilder{
		Schema: tfutils.SchemaMap{
			"host": tfutils.String().Default("http://localhost:2019"),
			"ssh": tfutils.SchemaMap{
				"host":     tfutils.String().Required(),
				"key_file": tfutils.String().Required(),
				"host_key": tfutils.String().Required(),
			}.IntoSet().Optional().MaxItems(1),
		},
		Resources: tfutils.ResourceMap{
			"caddy_http":   HTTP{},
			"caddy_server": Server{},
		},
		DataSources: tfutils.DataSourceMap{
			"caddy_server_route": ServerRoute{},
		},
		ConfigureFunc: providerConfigurer,
	}.Build()
}
