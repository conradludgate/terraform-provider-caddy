package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/conradludgate/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HTTP Represents the caddy HTTP Application
type HTTP struct{}

func (HTTP) Schema() tfutils.SchemaMap {
	return tfutils.SchemaMap{
		"http_port":    tfutils.Int().Default(80),
		"https_port":   tfutils.Int().Default(443),
		"grace_period": tfutils.String().Default("0s"),
	}
}

func (HTTP) Read(d *schema.ResourceData, m interface{}) error {
	c := m.(Client)

	server, err := c.GetHTTP()
	if err != nil {
		return err
	}

	var gracePeriod caddyapi.Duration
	if err := gracePeriod.UnmarshalText(GetString(d, "grace_period")); err != nil {
		return err
	}
	if gracePeriod != server.GracePeriod {
		d.Set("grace_period", server.GracePeriod.String())
	}

	d.Set("http_port", server.HTTPPort)
	d.Set("https_port", server.HTTPSPort)

	return nil
}

func (HTTP) Create(d *schema.ResourceData, m interface{}) error {
	c := m.(Client)

	http := caddyapi.HTTP{
		HTTPPort:  GetInt(d, "http_port"),
		HTTPSPort: GetInt(d, "https_port"),
	}
	if err := http.GracePeriod.UnmarshalText(GetString(d, "grace_period")); err != nil {
		return err
	}
	if err := c.CreateHTTP(http); err != nil {
		return err
	}
	d.SetId("@config/apps/http")

	return nil
}

func (HTTP) Update(d *schema.ResourceData, m interface{}) error {
	c := m.(Client)

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

func (HTTP) Delete(d *schema.ResourceData, m interface{}) error {
	c := m.(Client)
	return c.DeleteHTTP()
}
