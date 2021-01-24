package caddy

import (
	"github.com/conradludgate/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider for caddy
func Provider() *schema.Provider {
	return tfutils.Provider{
		Schema: tfutils.SchemaMap{
			"host": tfutils.String().Default("http://localhost:2019"),
			"ssh": tfutils.SchemaMap{
				"host":     tfutils.String().Required(true),
				"key_file": tfutils.String().Required(true),
				"host_key": tfutils.String().Required(true),
			}.IntoSet().Optional(true).MaxItems(1),
		},
		Resources: tfutils.ResourceMap{
			"caddy_http":   HTTP{},
			"caddy_server": Server{},
		},
		DataSources: tfutils.DataSourceMap{
			"caddy_server_route": ServerRoute{5},
		},
		ConfigureFunc: providerConfigurer,
	}.Build()
}
