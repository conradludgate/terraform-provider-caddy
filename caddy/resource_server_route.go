package caddy

import (
	"encoding/json"

	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/conradludgate/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServerRoute struct{}

func (sr ServerRoute) Schema() tfutils.SchemaMap {
	return tfutils.SchemaMap{
		"group":    tfutils.String().Optional(true),
		"terminal": tfutils.Bool().Optional(true),
		"match":    tfutils.ListOf(ServerRouteMatcher{}).Optional(true),
		"handle":   tfutils.ListOf(ServerRouteHandler{}).Optional(true),
	}
}

func (ServerRoute) Read(d *schema.ResourceData, m interface{}) error {
	r := caddyapi.Route{
		Group:    GetString(d, "group"),
		Terminal: GetBool(d, "terminal"),
		Matchers: ServerRouteMatchersFrom(GetObjectList(d, "match")),
		Handlers: ServerRouteHandlersFrom(GetObjectList(d, "handle")),
	}

	b, err := json.Marshal(r)
	if err != nil {
		return err
	}

	d.SetId(string(b))

	return nil
}

func ServerRouteFrom(d *MapData) caddyapi.Route {
	return caddyapi.Route{
		Group:    GetString(d, "group"),
		Terminal: GetBool(d, "terminal"),
		Matchers: ServerRouteMatchersFrom(GetObjectList(d, "match")),
		Handlers: ServerRouteHandlersFrom(GetObjectList(d, "handle")),
	}
}

func ServerRoutesFrom(d []MapData) []caddyapi.Route {
	routes := make([]caddyapi.Route, 0, len(d))
	for _, d := range d {
		routes = append(routes, ServerRouteFrom(&d))
	}
	return routes
}

func ServerRouteInto(route caddyapi.Route) map[string]interface{} {
	return map[string]interface{}{
		"group":    route.Group,
		"terminal": route.Terminal,
		"match":    ServerRouteMatchersInto(route.Matchers),
		"handle":   ServerRouteHandlersInto(route.Handlers),
	}
}

func ServerRoutesInto(routes []caddyapi.Route) []map[string]interface{} {
	d := make([]map[string]interface{}, 0, len(routes))
	for _, route := range routes {
		d = append(d, ServerRouteInto(route))
	}
	return d
}
