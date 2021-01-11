package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServerRouteMatch() *schema.Resource {
	return &schema.Resource{
		Read:   resourceServerRouteMatchRead,
		Create: resourceServerRouteMatchCreate,
		Update: resourceServerRouteMatchUpdate,
		Delete: resourceServerRouteMatchDelete,

		Schema: map[string]*schema.Schema{
			"route": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceServerRouteMatchRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	if d.Id() == "" {
		return nil
	}

	match, err := c.GetServerRouteMatch(d.Id())
	if err != nil {
		return err
	}

	d.Set("host", schema.NewSet(schema.HashString, sstoi(match.Host)))

	return nil
}

func resourceServerRouteMatchCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	id, err := c.CreateServerRouteMatch(d.Get("route").(string), caddyapi.Match{
		ID:   d.Id(),
		Host: sitos(d.Get("host").(*schema.Set).List()),
	})
	if err != nil {
		return err
	}
	d.SetId(id)
	return nil
}

func resourceServerRouteMatchUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	return c.UpdateServerRouteMatch(d.Id(), caddyapi.Match{
		Host: sitos(d.Get("host").(*schema.Set).List()),
	})
}

func resourceServerRouteMatchDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	return c.DeleteServerRouteMatch(d.Id())
}
