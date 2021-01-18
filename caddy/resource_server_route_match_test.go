package caddy

import (
	"testing"

	"github.com/conradludgate/terraform-provider-caddy/caddy/mocks"
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestServerRouteMatchNot(t *testing.T) {
	UnitTest(t, func(caddyMock *mocks.Client) {
		caddyMock.On("CreateServer", "Foo", serverMatchNot).Return("@config/apps/http/servers/Foo", nil)
		caddyMock.On("GetServer", "@config/apps/http/servers/Foo").Return(serverMatchNot, nil)
		caddyMock.On("DeleteServer", "@config/apps/http/servers/Foo").Return(nil)
	},
		resource.TestStep{
			Config: serverConfigMatchNot,
		},
	)
}

const serverConfigMatchNot = `
data "caddy_server_route" "server_test_route" {
	match {
		host = ["foo.example.com"]
		not {
			path = ["/foo/*"]
		}
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

var serverMatchNot = caddyapi.Server{
	Listen: []string{":443"},
	Routes: []caddyapi.Route{
		{
			Matchers: []caddyapi.Match{
				{
					Host: []string{"foo.example.com"},
					Not: []caddyapi.Match{
						{
							Path: []string{"/foo/*"},
						},
					},
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

func TestServerRouteMatchHeader(t *testing.T) {
	UnitTest(t, func(caddyMock *mocks.Client) {
		caddyMock.On("CreateServer", "Foo", serverMatchHeader).Return("@config/apps/http/servers/Foo", nil)
		caddyMock.On("GetServer", "@config/apps/http/servers/Foo").Return(serverMatchHeader, nil)
		caddyMock.On("DeleteServer", "@config/apps/http/servers/Foo").Return(nil)
	},
		resource.TestStep{
			Config: serverConfigMatchHeader,
		},
	)
}

const serverConfigMatchHeader = `
data "caddy_server_route" "server_test_route" {
	match {
		header {
			name = "Authorization"
			values = ["Bearer *"]
		}
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

var serverMatchHeader = caddyapi.Server{
	Listen: []string{":443"},
	Routes: []caddyapi.Route{
		{
			Matchers: []caddyapi.Match{
				{
					Header: map[string][]string{
						"Authorization": {"Bearer *"},
					},
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
