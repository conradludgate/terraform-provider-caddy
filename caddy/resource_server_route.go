package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServerRoute() *schema.Resource {
	return &schema.Resource{
		Read:   resourceServerRouteRead,
		Create: resourceServerRouteCreate,
		Update: resourceServerRouteUpdate,
		Delete: resourceServerRouteDelete,

		Schema: map[string]*schema.Schema{
			"server": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"terminal": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceServerRouteRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	if d.Id() == "" {
		return nil
	}

	route, err := c.GetServerRoute(d.Id())
	if err != nil {
		return err
	}

	if route.Group != nil {
		d.Set("group", route.Group)
	}
	if route.Terminal != nil {
		d.Set("terminal", route.Terminal)
	}

	return nil
}

func resourceServerRouteCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	id, err := c.CreateServerRoute(d.Get("server").(string), caddyapi.Route{
		ID:       d.Id(),
		Group:    GetOkS(d, "group"),
		Terminal: GetOkB(d, "terminal"),
	})
	if err != nil {
		return err
	}
	d.SetId(id)

	return nil
}

func resourceServerRouteUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	if d.HasChange("group") {
		c.UpdateServerRouteGroup(d.Id(), d.Get("group").(string))
	}

	if d.HasChange("terminal") {
		c.UpdateServerRouteTerminal(d.Id(), d.Get("terminal").(bool))
	}

	return nil
}

func resourceServerRouteDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)
	return c.DeleteServerRoute(d.Id())
}
