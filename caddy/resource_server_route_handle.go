package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/tfutils"
)

type ServerRouteHandler struct{}

func (ServerRouteHandler) Schema() tfutils.SchemaMap {
	return tfutils.SchemaMap{
		"static_response": tfutils.SchemaMap{
			"status_code": tfutils.String().Optional(),
			"header":      MapListString.Optional(),
			"body":        tfutils.String().Optional(),
			"close":       tfutils.Bool().Optional(),
		}.IntoSet().
			Optional().
			MaxItems(1),

		"reverse_proxy": tfutils.SchemaMap{
			"upstream": tfutils.SchemaMap{
				"dial": tfutils.String().Optional(),
			}.IntoList().Required(),
		}.IntoSet().
			Optional().
			MaxItems(1),

		"request_body": tfutils.SchemaMap{
			"max_size": tfutils.Int().Required(),
		}.IntoSet().
			Optional().
			MaxItems(1),

		"file_server": tfutils.SchemaMap{
			"root":           tfutils.String().Optional(),
			"hide":           tfutils.String().List().Optional(),
			"index_names":    tfutils.String().List().Optional(),
			"canonical_uris": tfutils.Bool().Optional(),
			"pass_thru":      tfutils.Bool().Optional(),
		}.IntoSet().
			Optional().
			MaxItems(1),

		"templates": tfutils.SchemaMap{
			"file_root":  tfutils.String().Optional(),
			"mime_types": tfutils.String().List().Optional(),
			"delimiters": tfutils.String().List().Optional(),
		}.IntoSet().
			Optional().
			MaxItems(1),
	}
}

func ServerRouteHandlersFrom(d []MapData) []map[string]interface{} {
	handlers := make([]map[string]interface{}, 0, len(d))
	for _, d := range d {
		handlers = append(handlers, ServerRouteHandlerFrom(&d))
	}
	return handlers
}

func ServerRouteHandlerFrom(d *MapData) map[string]interface{} {
	handlerFuncs := map[string]func(d *MapData) map[string]interface{}{
		"static_response": StaticResponseFrom,
		"reverse_proxy":   ReverseProxyFrom,
		"request_body":    RequestBodyFrom,
		"file_server":     FileServerFrom,
		"templates":       TemplatesFrom,
	}

	for key, fn := range handlerFuncs {
		if d := GetObjectSet(d, key); len(d) == 1 {
			m := fn(&d[0])
			m["handler"] = key
			return m
		}
	}

	return nil
}

func StaticResponseFrom(d *MapData) map[string]interface{} {
	return map[string]interface{}{
		"status_code": GetString(d, "status_code"),
		"headers":     ParseMapListString(d, "header"),
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
		"upstreams": upstreams,
	}
}

func RequestBodyFrom(d *MapData) map[string]interface{} {
	return map[string]interface{}{
		"max_size": GetInt(d, "max_size"),
	}
}

func FileServerFrom(d *MapData) map[string]interface{} {
	return map[string]interface{}{
		"root":           GetString(d, "root"),
		"hide":           GetStringList(d, "hide"),
		"index_names":    GetStringList(d, "index_names"),
		"canonical_uris": GetBool(d, "canonical_uris"),
		"pass_thru":      GetBool(d, "pass_thru"),
	}
}

func TemplatesFrom(d *MapData) map[string]interface{} {
	return map[string]interface{}{
		"file_root":  GetString(d, "file_root"),
		"mime_types": GetStringList(d, "mime_types"),
		"delimiters": GetStringList(d, "delimiters"),
	}
}
