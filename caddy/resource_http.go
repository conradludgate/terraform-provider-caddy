package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceHTTP() *schema.Resource {
	return &schema.Resource{
		Create: resourceHTTPCreate,
		Read:   resourceHTTPRead,
		Update: resourceHTTPUpdate,
		Delete: resourceHTTPDelete,

		Schema: map[string]*schema.Schema{
			"http_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     80,
				Description: "specifies the port to use for HTTP (as opposed to HTTPS), which is used when setting up HTTP->HTTPS redirects or ACME HTTP challenge solvers. Default: 80.",
			},
			"https_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     443,
				Description: "specifies the port to use for HTTPS, which is used when solving the ACME TLS-ALPN challenges, or whenever HTTPS is needed but no specific port number is given. Default: 443.",
			},
			"grace_period": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0s",
				Description: `
					How long to wait for active connections when shutting down the server. Once the grace period is over, connections will be forcefully closed.
					Duration can be an integer or a string. An integer is interpreted as nanoseconds. If a string, it is a Go time.Duration value such as 300ms, 1.5h, or 2h45m; valid units are ns, us/µs, ms, s, m, h, and d.
				`,
			},
		},
	}
}

func resourceHTTPRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	server, err := c.GetHTTP()
	if err != nil {
		return err
	}

	d.Set("http_port", server.HTTPPort)
	d.Set("https_port", server.HTTPSPort)
	d.Set("grace_period", server.GracePeriod.String())

	return nil
}

func resourceHTTPCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	http := caddyapi.HTTP{
		HTTPPort:  d.Get("http_port").(int),
		HTTPSPort: d.Get("https_port").(int),
	}
	if err := http.GracePeriod.UnmarshalText(d.Get("grace_period").(string)); err != nil {
		return err
	}
	if err := c.CreateHTTP(http); err != nil {
		return err
	}
	d.SetId("@config/apps/http")

	return nil
}

func resourceHTTPUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	if d.HasChange("http_port") {
		if err := c.UpdateHTTPPort(d.Get("http_port").(int)); err != nil {
			return err
		}
	}
	if d.HasChange("https_port") {
		if err := c.UpdateHTTPSPort(d.Get("https_port").(int)); err != nil {
			return err
		}
	}
	if d.HasChange("grace_period") {
		var dur caddyapi.Duration
		if err := dur.UnmarshalText(d.Get("grace_period").(string)); err != nil {
			return err
		}
		if err := c.UpdateHTTPGracePeriod(dur); err != nil {
			return err
		}
	}

	return nil
}

func resourceHTTPDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)
	return c.DeleteHTTP()
}
