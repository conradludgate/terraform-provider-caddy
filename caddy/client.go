package caddy

import "github.com/conradludgate/terraform-provider-caddy/caddyapi"

//go:generate mockery --name=Client

type Client interface {
	GetHTTP() (caddyapi.HTTP, error)
	DeleteHTTP() error
	CreateHTTP(http caddyapi.HTTP) error
	UpdateHTTPPort(httpPort int) error
	UpdateHTTPSPort(httpsPort int) error
	UpdateHTTPGracePeriod(gracePeriod caddyapi.Duration) error

	CreateServer(name string, server caddyapi.Server) (string, error)
	UpdateServerListen(id string, listen []string) error
	UpdateServerRoutes(id string, routes []caddyapi.Route) error
	GetServer(id string) (caddyapi.Server, error)
	DeleteServer(id string) error
}
