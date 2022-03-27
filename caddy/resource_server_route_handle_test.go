package caddy

import (
	"testing"

	"github.com/conradludgate/terraform-provider-caddy/caddy/mocks"
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestServerRouteHandleSubroute(t *testing.T) {
	UnitTest(t, func(caddyMock *mocks.Client) {
		caddyMock.On("CreateServer", "Foo", serverHandleSubroute).Return("@config/apps/http/servers/Foo", nil)
		caddyMock.On("GetServer", "@config/apps/http/servers/Foo").Return(&serverHandleSubroute, nil)
		caddyMock.On("DeleteServer", "@config/apps/http/servers/Foo").Return(nil)
	},
		resource.TestStep{
			Config: serverConfigHandleSubroute,
		},
	)
}

const serverConfigHandleSubroute = `
data "caddy_server_route" "server_test_route" {
	match {
		host = ["example.com"]
	}

	handle {
		subroute {
			route {
				match {
					path = ["/foo"]
				}
				handle {
					static_response {
						body = "foo"
					}
				}
			}
			route {
				match {
					path = ["/bar"]
				}
				handle {
					static_response {
						body = "bar"
					}
				}
			}
		}
	}
}

resource "caddy_server" "server_test" {
	name = "Foo"
	listen = [":443"]

	routes = [data.caddy_server_route.server_test_route.id]
}
`

var serverHandleSubroute = caddyapi.Server{
	Listen: []string{":443"},
	Routes: []caddyapi.Route{
		{
			Matchers: []caddyapi.Match{
				{
					Host: []string{"example.com"},
				},
			},
			Handlers: []caddyapi.HandleMarshal{
				{
					Handle: caddyapi.Subroute{
						Routes: []caddyapi.Route{
							{
								Matchers: []caddyapi.Match{
									{
										Path: []string{"/foo"},
									},
								},
								Handlers: []caddyapi.HandleMarshal{
									{
										Handle: caddyapi.StaticResponse{
											Body: "foo",
										},
									},
								},
							},
							{
								Matchers: []caddyapi.Match{
									{
										Path: []string{"/bar"},
									},
								},
								Handlers: []caddyapi.HandleMarshal{
									{
										Handle: caddyapi.StaticResponse{
											Body: "bar",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}
