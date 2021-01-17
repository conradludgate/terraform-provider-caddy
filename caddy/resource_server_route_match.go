package caddy

import (
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/conradludgate/terraform-provider-caddy/tfutils"
)

type ServerRouteMatcher struct{}

func (ServerRouteMatcher) Schema() tfutils.SchemaMap {
	return tfutils.SchemaMap{
		"host": tfutils.String().List().Optional(),
	}
}

func ServerRouteMatcherFrom(d *MapData) caddyapi.Match {
	return caddyapi.Match{
		Host: GetStringList(d, "host"),
	}
}

func ServerRouteMatchersFrom(d []MapData) []caddyapi.Match {
	matchers := make([]caddyapi.Match, 0, len(d))
	for _, d := range d {
		matchers = append(matchers, ServerRouteMatcherFrom(&d))
	}
	return matchers
}
