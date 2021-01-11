package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServerRouteHandle() *schema.Resource {
	return &schema.Resource{
		Read:   resourceServerRouteHandleRead,
		Create: resourceServerRouteHandleCreate,
		Update: resourceServerRouteHandleUpdate,
		Delete: resourceServerRouteHandleDelete,

		Schema: map[string]*schema.Schema{
			"route": {
				Type:     schema.TypeString,
				Required: true,
			},
			"handler": {
				Type:     schema.TypeString,
				Required: true,
			},

			"body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"headers": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"upstream": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dial": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceServerRouteHandleRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	if d.Id() == "" {
		return nil
	}

	handle, err := c.GetServerRouteHandle(d.Id())
	if err != nil {
		return err
	}

	d.Set("handler", handle.Handler)
	d.Set("body", handle.Body)

	return nil
}

func resourceServerRouteHandleCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	id, err := c.CreateServerRouteHandle(d.Get("route").(string), caddyapi.Handle{
		ID:      d.Id(),
		Handler: d.Get("handler").(string),
		Body:    GetOkS(d, "body"),
		// Headers: ,
		Upstreams: setIntoUpstreams(d.Get("upstream").(*schema.Set)),
	})
	if err != nil {
		return err
	}
	d.SetId(id)
	return nil
}

func resourceServerRouteHandleUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	return c.UpdateServerRouteHandle(d.Id(), caddyapi.Handle{
		ID:      d.Id(),
		Handler: d.Get("handler").(string),
		Body:    GetOkS(d, "body"),
		// Headers: ,
		Upstreams: setIntoUpstreams(d.Get("upstream").(*schema.Set)),
	})
}

func resourceServerRouteHandleDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	return c.DeleteServerRouteHandle(d.Id())
}

func setIntoUpstreams(set *schema.Set) []caddyapi.Upstream {
	var upstreams []caddyapi.Upstream
	for _, v := range set.List() {
		data := v.(map[string]interface{})
		upstreams = append(upstreams, caddyapi.Upstream{
			Dial: data["dial"].(string),
		})
	}
	return upstreams
}
