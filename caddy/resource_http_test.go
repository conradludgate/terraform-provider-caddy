package caddy

import (
	"testing"
	"time"

	"github.com/conradludgate/terraform-provider-caddy/caddy/mocks"
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestHTTP(t *testing.T) {
	UnitTest(t, func(caddyMock *mocks.Client) {
		caddyMock.On("CreateHTTP", http_).Return(nil)
		caddyMock.On("GetHTTP").Return(&http_, nil)
		caddyMock.On("DeleteHTTP").Return(nil)
	},
		resource.TestStep{
			Config: httpConfig,
		},
	)
}

const httpConfig = `
resource "caddy_http" "http" {
	http_port = 90
	https_port = 444
	grace_period = "10m"
}
`

var http_ = caddyapi.HTTP{
	HTTPPort:    90,
	HTTPSPort:   444,
	GracePeriod: caddyapi.Duration(10 * time.Minute),
}
