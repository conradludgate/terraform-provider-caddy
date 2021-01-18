package caddy

import (
	"encoding/json"

	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/conradludgate/terraform-provider-caddy/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Server struct{}

func (Server) Schema() tfutils.SchemaMap {
	return tfutils.SchemaMap{
		"name":   tfutils.String().Required(),
		"listen": tfutils.String().List().Required(),
		"route":  tfutils.ListOf(ServerRoute{}).Optional(),
		"routes": tfutils.String().List().Optional().ConflictsWith("route"),
	}
}

func (Server) Read(d *schema.ResourceData, m interface{}) error {
	c := m.(Client)

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

func (Server) Create(d *schema.ResourceData, m interface{}) error {
	c := m.(Client)

	server := caddyapi.Server{
		Listen: GetStringList(d, "listen"),
		ID:     d.Id(),
	}

	if routes := GetStringList(d, "routes"); len(routes) > 0 {
		server.Routes = make([]caddyapi.Route, len(routes))
		for i, route := range routes {
			if err := json.Unmarshal([]byte(route), &server.Routes[i]); err != nil {
				return err
			}
		}
	} else {
		server.Routes = ServerRoutesFrom(GetObjectList(d, "route"))
	}

	id, err := c.CreateServer(GetString(d, "name"), server)
	if err != nil {
		return err
	}
	d.SetId(id)

	return nil
}

func (Server) Update(d *schema.ResourceData, m interface{}) error {
	c := m.(Client)

	if d.HasChange("listen") {
		if err := c.UpdateServerListen(d.Id(), GetStringList(d, "listen")); err != nil {
			return err
		}
	}

	if d.HasChange("route") {
		if err := c.UpdateServerRoutes(d.Id(), ServerRoutesFrom(GetObjectList(d, "route"))); err != nil {
			return err
		}
	}

	if d.HasChange("routes") {
		jsonRoutes := GetStringList(d, "routes")
		routes := make([]caddyapi.Route, len(jsonRoutes))
		for i, route := range jsonRoutes {
			if err := json.Unmarshal([]byte(route), &routes[i]); err != nil {
				return err
			}
		}
		if err := c.UpdateServerRoutes(d.Id(), routes); err != nil {
			return err
		}
	}

	return nil
}

func (Server) Delete(d *schema.ResourceData, m interface{}) error {
	c := m.(Client)
	return c.DeleteServer(d.Id())
}
