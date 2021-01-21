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
		"group":    tfutils.String().Optional(),
		"terminal": tfutils.Bool().Optional(),
		"match":    tfutils.ListOf(ServerRouteMatcher{}).Optional(),
		"handle":   tfutils.ListOf(ServerRouteHandler{}).Optional(),
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

// func (ServerRoute) Read(d *schema.Resource, m interface{}) error {
// 	return nil
// }
