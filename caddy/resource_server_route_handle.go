package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/conradludgate/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServerRouteHandler struct{}

func (ServerRouteHandler) Schema() tfutils.SchemaMap {
	return tfutils.SchemaMap{
		"static_response": tfutils.SchemaMap{
			"status_code": tfutils.String().Optional(true),
			"header":      MapListString.Optional(true),
			"body":        tfutils.String().Optional(true),
			"close":       tfutils.Bool().Optional(true),
		}.IntoSet().
			Optional(true).
			MaxItems(1),

		"reverse_proxy": tfutils.SchemaMap{
			"upstream": tfutils.SchemaMap{
				"dial": tfutils.String().Optional(true),
			}.IntoList().Required(true),
		}.IntoSet().
			Optional(true).
			MaxItems(1),

		"request_body": tfutils.SchemaMap{
			"max_size": tfutils.Int().Required(true),
		}.IntoSet().
			Optional(true).
			MaxItems(1),

		"file_server": tfutils.SchemaMap{
			"root":           tfutils.String().Optional(true),
			"hide":           tfutils.String().List().Optional(true),
			"index_names":    tfutils.String().List().Optional(true),
			"canonical_uris": tfutils.Bool().Optional(true),
			"pass_thru":      tfutils.Bool().Optional(true),
		}.IntoSet().
			Optional(true).
			MaxItems(1),

		"templates": tfutils.SchemaMap{
			"file_root":  tfutils.String().Optional(true),
			"mime_types": tfutils.String().List().Optional(true),
			"delimiters": tfutils.String().List().Optional(true),
		}.IntoSet().
			Optional(true).
			MaxItems(1),
	}
}

func ServerRouteHandlersFrom(d []MapData) []caddyapi.HandleMarshal {
	handlers := make([]caddyapi.HandleMarshal, 0, len(d))
	for _, d := range d {
		handlers = append(handlers, ServerRouteHandlerFrom(&d))
	}
	return handlers
}

func ServerRouteHandlerFrom(d *MapData) caddyapi.HandleMarshal {
	for key, v := range *d {
		s := AsSet(v)
		if len(s) != 1 {
			continue
		}
		d := &s[0]

		switch key {
		case "static_response":
			return caddyapi.HandleMarshal{Handle: IntoStaticResponse(d)}
		case "reverse_proxy":
			return caddyapi.HandleMarshal{Handle: IntoReverseProxy(d)}
		case "request_body":
			return caddyapi.HandleMarshal{Handle: IntoRequestBody(d)}
		case "file_server":
			return caddyapi.HandleMarshal{Handle: IntoFileServer(d)}
		case "templates":
			return caddyapi.HandleMarshal{Handle: IntoTemplates(d)}
		}
	}

	panic("no handler")
}

func IntoStaticResponse(d *MapData) caddyapi.StaticResponse {
	return caddyapi.StaticResponse{
		StatusCode: GetString(d, "status_code"),
		Headers:    ParseMapListString(d, "header"),
		Body:       GetString(d, "body"),
		Close:      GetBool(d, "close"),
	}
}

func IntoReverseProxy(d *MapData) caddyapi.ReverseProxy {
	var upstreams []caddyapi.Upstream
	upstreamList := GetObjectList(d, "upstream")
	for _, d := range upstreamList {
		upstreams = append(upstreams, caddyapi.Upstream{
			Dial: GetString(&d, "dial"),
		})
	}

	return caddyapi.ReverseProxy{
		Upstreams: upstreams,
	}
}

func IntoRequestBody(d *MapData) caddyapi.RequestBody {
	return caddyapi.RequestBody{
		MaxSize: GetInt(d, "max_size"),
	}
}

func IntoFileServer(d *MapData) caddyapi.FileServer {
	return caddyapi.FileServer{
		Root:          GetString(d, "root"),
		Hide:          GetStringList(d, "hide"),
		IndexNames:    GetStringList(d, "index_names"),
		CanonicalURIs: GetBool(d, "canonical_uris"),
		PassThru:      GetBool(d, "pass_thru"),
	}
}

func IntoTemplates(d *MapData) caddyapi.Templates {
	return caddyapi.Templates{
		FileRoot:   GetString(d, "file_root"),
		MimeTypes:  GetStringList(d, "mime_types"),
		Delimiters: GetStringList(d, "delimiters"),
	}
}

func ServerRouteHandlersInto(handlers []caddyapi.HandleMarshal) []map[string]interface{} {
	d := make([]map[string]interface{}, 0, len(handlers))
	for _, handle := range handlers {
		d = append(d, ServerRouteHandlerInto(handle))
	}
	return d
}

func ServerRouteHandlerInto(handle caddyapi.HandleMarshal) map[string]interface{} {
	h := handle.Handle

	var key string
	var val interface{}

	switch h.(type) {
	case caddyapi.StaticResponse:
		key = "static_response"
		val = FromStaticResponse(h.(caddyapi.StaticResponse))
	case caddyapi.ReverseProxy:
		key = "reverse_proxy"
		val = FromReverseProxy(h.(caddyapi.ReverseProxy))
	case caddyapi.RequestBody:
		key = "request_body"
		val = FromRequestBody(h.(caddyapi.RequestBody))
	case caddyapi.FileServer:
		key = "file_server"
		val = FromFileServer(h.(caddyapi.FileServer))
	case caddyapi.Templates:
		key = "templates"
		val = FromTemplates(h.(caddyapi.Templates))
	default:
		panic("unexpected handler type")
	}

	m := map[string]interface{}{}
	m[key] = schema.NewSet(schema.HashResource(
		ServerRouteHandler{}.Schema().BuildResource().Schema[key].Elem.(*schema.Resource),
	), []interface{}{val})
	return m
}

func FromStaticResponse(r caddyapi.StaticResponse) map[string]interface{} {
	return map[string]interface{}{
		"status_code": r.StatusCode,
		"headers":     IntoMapListString(r.Headers),
		"body":        r.Body,
		"close":       r.Close,
	}
}

func FromReverseProxy(r caddyapi.ReverseProxy) map[string]interface{} {
	upstreams := make([]map[string]interface{}, len(r.Upstreams))
	for i, upstream := range r.Upstreams {
		upstreams[i] = map[string]interface{}{
			"dial": upstream.Dial,
		}
	}

	return map[string]interface{}{
		"upstreams": upstreams,
	}
}

func FromRequestBody(r caddyapi.RequestBody) map[string]interface{} {
	return map[string]interface{}{
		"max_size": r.MaxSize,
	}
}

func FromFileServer(f caddyapi.FileServer) map[string]interface{} {
	return map[string]interface{}{
		"root":           f.Root,
		"hide":           f.Hide,
		"index_names":    f.IndexNames,
		"canonical_uris": f.CanonicalURIs,
		"pass_thru":      f.PassThru,
	}
}

func FromTemplates(t caddyapi.Templates) map[string]interface{} {
	return map[string]interface{}{
		"file_root":  t.FileRoot,
		"mime_types": t.MimeTypes,
		"delimiters": t.Delimiters,
	}
}
