package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"http": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listen": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	if d.Id() == "" {
		return nil
	}

	server, err := c.GetServer(d.Id())
	if err != nil {
		return err
	}

	d.Set("listen", server.Listen)

	return nil
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	server := caddyapi.Server{
		Listen: sitos(d.Get("listen").([]interface{})),
		ID:     d.Id(),
	}

	id, err := c.CreateServer(d.Get("name").(string), server)
	if err != nil {
		return err
	}
	d.SetId(id)

	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)

	if d.HasChange("listen") {
		c.UpdateServerListen(d.Id(), sitos(d.Get("listen").([]interface{})))
	}

	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*caddyapi.Client)
	return c.DeleteServer(d.Id())
}
