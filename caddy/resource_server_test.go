package caddy

import (
	"testing"

	"github.com/conradludgate/terraform-provider-caddy/caddy/mocks"
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestServer(t *testing.T) {
	UnitTest(t, func(caddyMock *mocks.Client) {
		caddyMock.On("CreateServer", "Foo", server).Return("@config/apps/http/servers/Foo", nil)
		caddyMock.On("GetServer", "@config/apps/http/servers/Foo").Return(server, nil)
		caddyMock.On("DeleteServer", "@config/apps/http/servers/Foo").Return(nil)
	},
		resource.TestStep{
			Config: serverConfig,
		},
	)
}

const serverConfig = `
resource "caddy_server" "server_test" {
	name = "Foo"
	listen = [":443"]

	route {
		match {
			host = ["foo.example.com"]
		}

		handle {
			static_response {
				body = "hello world"
			}
		}
	}
}
`

var server = caddyapi.Server{
	Listen: []string{":443"},
	Routes: []caddyapi.Route{
		{
			Matchers: []caddyapi.Match{
				{
					Host: []string{"foo.example.com"},
				},
			},
			Handlers: []map[string]interface{}{
				{
					"status_code": "",
					"handler":     "static_response",
					"body":        "hello world",
					"headers":     map[string]string(nil),
					"close":       false,
				},
			},
		},
	},
}

func TestServerUpdate(t *testing.T) {
	UnitTest(t, func(caddyMock *mocks.Client) {
		caddyMock.On("CreateServer", "Foo", server).Return("@config/apps/http/servers/Foo", nil)
		caddyMock.On("GetServer", "@config/apps/http/servers/Foo").Return(server, nil)
		caddyMock.On("UpdateServerRoutes", "@config/apps/http/servers/Foo", serverUpdateRouteMatch.Routes).Return(nil)
		caddyMock.On("DeleteServer", "@config/apps/http/servers/Foo").Return(nil)
	},
		resource.TestStep{
			Config: serverConfig,
		},
		resource.TestStep{
			Config: serverConfigUpdateRouteMatch,
		},
	)
}

const serverConfigUpdateRouteMatch = `
resource "caddy_server" "server_test" {
	name = "Foo"
	listen = [":443"]

	route {
		match {
			host = ["bar.example.com"]
		}

		handle {
			reverse_proxy {
				upstream {
					dial = "localhost:2020"
				}
			}
		}
	}
}
`

var serverUpdateRouteMatch = caddyapi.Server{
	Listen: []string{":443"},
	Routes: []caddyapi.Route{
		{
			Matchers: []caddyapi.Match{
				{
					Host: []string{"bar.example.com"},
				},
			},
			Handlers: []map[string]interface{}{
				{
					"handler": "reverse_proxy",
					"upstreams": []map[string]interface{}{
						{
							"dial": "localhost:2020",
						},
					},
				},
			},
		},
	},
}

func TestServerSeparated(t *testing.T) {
	UnitTest(t, func(caddyMock *mocks.Client) {
		caddyMock.On("CreateServer", "Foo", serverSeparated).Return("@config/apps/http/servers/Foo", nil)
		caddyMock.On("GetServer", "@config/apps/http/servers/Foo").Return(serverSeparated, nil)
		caddyMock.On("DeleteServer", "@config/apps/http/servers/Foo").Return(nil)
	},
		resource.TestStep{
			Config: serverConfigSeperated,
		},
	)
}

const serverConfigSeperated = `
data "caddy_server_route" "server_test_route" {
	match {
		host = ["foo.example.com"]
	}

	handle {
		static_response {
			body = "hello world"
		}
	}
}

resource "caddy_server" "server_test" {
	name = "Foo"
	listen = [":443"]

	routes = [data.caddy_server_route.server_test_route.id]
}
`

var serverSeparated = caddyapi.Server{
	Listen: []string{":443"},
	Routes: []caddyapi.Route{
		{
			Matchers: []caddyapi.Match{
				{
					Host: []string{"foo.example.com"},
				},
			},
			Handlers: []map[string]interface{}{
				{
					"status_code": "",
					"handler":     "static_response",
					"body":        "hello world",
					"headers":     nil,
					"close":       false,
				},
			},
		},
	},
}
