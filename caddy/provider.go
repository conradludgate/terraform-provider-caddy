package caddy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider for caddy
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "http://localhost:2019", // ssh://terraform@ssh.conradludgate.com:22/tmp/caddy-admin.sock
				Description: "Caddy Admin API host. Must be a valid URL. If the scheme is `ssh`, then the path value is expected to be a unix socket",
			},
			"ssh_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSH Private key file. only needed if `host` points to an ssh server",
			},
			"host_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ignore_host_key"},
				Description:   "SSH Host key file. Only needed if `host` points to an ssh server",
			},
			"ignore_host_key": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"host_key"},
				Description:   "Ignore SSH Host Key",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"caddy_http":                resourceHTTP(),
			"caddy_server":              resourceServer(),
			"caddy_server_route":        resourceServerRoute(),
			"caddy_server_route_match":  resourceServerRouteMatch(),
			"caddy_server_route_handle": resourceServerRouteHandle(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  providerConfigurer,
	}
}
