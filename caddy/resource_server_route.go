package caddy

import (
	"fmt"
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/conradludgate/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServerRouteResource struct {
	Nested int
}

func (sr ServerRouteResource) Schema() tfutils.SchemaMap {
	return tfutils.SchemaMap{
		"server_name": tfutils.String().Required(true),
		"route_id":    tfutils.String().Required(true),
		"group":       tfutils.String().Optional(true),
		"terminal":    tfutils.Bool().Optional(true),
		"match":       tfutils.ListOf(ServerRouteMatcher{}).Optional(true),
		"handle":      tfutils.ListOf(ServerRouteHandler{sr.Nested}).Optional(true),
	}
}

func (ServerRouteResource) Read(d *schema.ResourceData, m interface{}) error {
	c := m.(Client)

	serverName := GetString(d, "server_name")
	serverId := "@config/apps/http/servers/" + serverName
	server, err := c.GetServer(serverId)
	if err != nil {
		return err
	}

	for _, route := range server.Routes {
		if route.ID == d.Id() {
			d.Set("route_id", route.ID)
			d.Set("group", route.Group)
			d.Set("terminal", route.Terminal)
			d.Set("match", ServerRouteMatchersInto(route.Matchers))
			d.Set("handle", ServerRouteHandlersInto(route.Handlers))

			return nil
		}
	}

	return nil
}

func (ServerRouteResource) Create(d *schema.ResourceData, m interface{}) error {
	c := m.(Client)

	serverName := GetString(d, "server_name")
	serverId := "@config/apps/http/servers/" + serverName
	server, err := c.GetServer(serverId)
	if err != nil {
		return err
	}

	routeId := GetString(d, "route_id")
	route := caddyapi.Route{
		ID:       routeId,
		Group:    GetString(d, "group"),
		Terminal: GetBool(d, "terminal"),
		Matchers: ServerRouteMatchersFrom(GetObjectList(d, "match")),
		Handlers: ServerRouteHandlersFrom(GetObjectList(d, "handle")),
	}

	server.Routes = append(server.Routes, route)
	if err := c.UpdateServerRoutes(serverId, server.Routes); err != nil {
		return err
	}

	d.SetId(routeId)
	return nil
}

func (ServerRouteResource) Update(d *schema.ResourceData, m interface{}) error {
	if !d.HasChangesExcept() {
		return nil
	}

	c := m.(Client)
	serverName := GetString(d, "server_name")

	serverId := "@config/apps/http/servers/" + serverName
	server, err := c.GetServer(serverId)
	if err != nil {
		return err
	}

	routeId := GetString(d, "route_id")

	newRoute := caddyapi.Route{
		ID:       routeId,
		Group:    GetString(d, "group"),
		Terminal: GetBool(d, "terminal"),
		Matchers: ServerRouteMatchersFrom(GetObjectList(d, "match")),
		Handlers: ServerRouteHandlersFrom(GetObjectList(d, "handle")),
	}

	for i, route := range server.Routes {
		if route.ID == d.Id() {
			server.Routes[i] = newRoute
			if err := c.UpdateServerRoutes(serverId, server.Routes); err != nil {
				return err
			}

			d.SetId(routeId)
			return nil
		}
	}

	return fmt.Errorf("route %s not found", d.Id())
}

func (ServerRouteResource) Delete(d *schema.ResourceData, m interface{}) error {
	c := m.(Client)
	serverName := GetString(d, "server_name")
	serverId := "@config/apps/http/servers/" + serverName
	server, err := c.GetServer(serverId)
	if err != nil {
		return err
	}

	for i, route := range server.Routes {
		if route.ID == d.Id() {
			server.Routes = append(server.Routes[:i], server.Routes[i+1:]...)
			err := c.UpdateServerRoutes(serverId, server.Routes)
			if err != nil {
				return err
			}

			break
		}
	}

	return nil
}
