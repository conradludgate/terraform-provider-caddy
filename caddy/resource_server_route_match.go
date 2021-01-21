package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/conradludgate/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MapListString is a schema that represents map[string][]string
var MapListString = tfutils.SchemaMap{
	"name":   tfutils.String().Required(true),
	"values": tfutils.String().List().Required(true),
}.IntoSet().SetFunc(func(v interface{}) int {
	m := v.(map[string]interface{})
	return schema.HashString(m["name"])
})

// ParseMapListString converts the data from a MapListString schema to a map[string][]string type
func ParseMapListString(d *MapData, key string) map[string][]string {
	sets := GetObjectSet(d, key)
	var values map[string][]string
	if len(sets) > 0 {
		values = make(map[string][]string, len(sets))
		for _, d := range sets {
			values[GetString(&d, "name")] = GetStringList(&d, "values")
		}
	}
	return values
}

type ServerRouteMatcher struct {
	not bool
}

func (s ServerRouteMatcher) Schema() tfutils.SchemaMap {
	sm := tfutils.SchemaMap{
		"host":   tfutils.String().List().Optional(true),
		"path":   tfutils.String().List().Optional(true),
		"method": tfutils.String().List().Optional(true),
		"header": MapListString.Optional(true),
		"query":  MapListString.Optional(true),
	}
	if !s.not {
		sm["not"] = tfutils.ListOf(ServerRouteMatcher{true}).Optional(true)
	}
	return sm
}

func ServerRouteMatcherFrom(d *MapData) caddyapi.Match {
	match := caddyapi.Match{
		Host:   GetStringList(d, "host"),
		Path:   GetStringList(d, "path"),
		Method: GetStringList(d, "method"),
		Header: ParseMapListString(d, "header"),
		Query:  ParseMapListString(d, "query"),
	}

	if nots := GetObjectList(d, "not"); len(nots) > 0 {
		match.Not = ServerRouteMatchersFrom(nots)
	}

	return match
}

func ServerRouteMatchersFrom(d []MapData) []caddyapi.Match {
	matchers := make([]caddyapi.Match, 0, len(d))
	for _, d := range d {
		matchers = append(matchers, ServerRouteMatcherFrom(&d))
	}
	return matchers
}
