package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/tfutils"
)

type ServerRouteHandler struct{}

func (ServerRouteHandler) Schema() tfutils.SchemaMap {
	return tfutils.SchemaMap{
		"static_response": tfutils.SchemaMap{
			"status_code": tfutils.String().Optional(),
			"headers":     tfutils.String().Map().Optional(),
			"body":        tfutils.String().Optional(),
			"close":       tfutils.Bool().Optional(),
		}.IntoSet().
			Optional(),

		"reverse_proxy": tfutils.SchemaMap{
			"upstream": tfutils.SchemaMap{
				"dial": tfutils.String().Optional(),
			}.IntoList().Optional(),
		}.IntoSet().
			Optional(),
	}
}

func ServerRouteHandlerFrom(d *MapData) map[string]interface{} {
	if d := GetObjectSet(d, "static_response"); len(d) == 1 {
		return StaticResponseFrom(&d[0])
	}

	if d := GetObjectSet(d, "reverse_proxy"); len(d) == 1 {
		return ReverseProxyFrom(&d[0])
	}

	return nil
}

func ServerRouteHandlersFrom(d []MapData) []map[string]interface{} {
	handlers := make([]map[string]interface{}, 0, len(d))
	for _, d := range d {
		handlers = append(handlers, ServerRouteHandlerFrom(&d))
	}
	return handlers
}

func StaticResponseFrom(d *MapData) map[string]interface{} {
	return map[string]interface{}{
		"handler":     "static_response",
		"status_code": GetString(d, "status_code"),
		"headers":     GetStringMap(d, "headers"),
		"body":        GetString(d, "body"),
		"close":       GetBool(d, "close"),
	}
}

func ReverseProxyFrom(d *MapData) map[string]interface{} {
	var upstreams []map[string]interface{}
	upstreamList := GetObjectList(d, "upstream")
	for _, d := range upstreamList {
		upstreams = append(upstreams, map[string]interface{}{
			"dial": GetString(&d, "dial"),
		})
	}

	return map[string]interface{}{
		"handler":  "reverse_proxy",
		"upstream": upstreams,
	}
}
