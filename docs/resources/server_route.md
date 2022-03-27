---
page_title: "caddy_server_route Resource - terraform-provider-caddy"
subcategory: ""
description: |-
  
---

# Resource `caddy_server_route`

https://caddyserver.com/docs/json/apps/http/servers/routes/

## Examples

```tf
resource "caddy_server_route" "foo" {
  server_name = "foo_server"
  route_id    = "foo"
  match {
    host = "foo.example.com"
  }

  handler {
    static_response {
      body = "Hello World!"
    }
  }
}
```

## Schema

### Optional
- **route_id** (String) Unique route identifier to find route in server route (populates @id)
- **server_name** (String) Name of the existing server this route will be added to
- **group** (String) group is an optional name for a group to which this route belongs. Grouping a route makes it mutually exclusive with others in its group; if a route belongs to a group, only the first matching route in that group will be executed.
- **terminal** (Boolean) If true, no more routes will be executed after this one
- **handle** (Block List) (see [below for nested schema](#nestedblock--handle))
- **match** (Block List) (see [below for nested schema](#nestedblock--match))

<a id="nestedblock--handle"></a>
### Nested Schema for `handle`

Optional:

Only one of these sets can be specified in a single handle. However, a route can have multiple handles

- **reverse_proxy** (Block Set) (see [below for nested schema](#nestedblock--handle--reverse_proxy))
- **static_response** (Block Set) (see [below for nested schema](#nestedblock--handle--static_response))

<a id="nestedblock--handle--reverse_proxy"></a>
### Nested Schema for `handle.reverse_proxy`

Reverse proxy takes the request and forwards it to one of the given upstreams.
https://caddyserver.com/docs/json/apps/http/servers/routes/handle/reverse_proxy/

Optional:

- **upstream** (Block List) (see [below for nested schema](#nestedblock--handle--reverse_proxy--upstream))

<a id="nestedblock--handle--reverse_proxy--upstream"></a>
### Nested Schema for `handle.reverse_proxy.upstream`

Optional:

- **dial** (String) The network address to dial to connect to the upstream. Must represent precisely one socket (i.e. no port ranges). A valid network address either has a host and port or is a unix socket address.

<a id="nestedblock--handle--static_response"></a>
### Nested Schema for `handle.static_response`

Static Response returns the static body and headers supplied.

Optional:

- **body** (String)
- **close** (Boolean) If true, the server will close the client's connection after writing the response.
- **headers** (Map of String)
- **status_code** (String)



<a id="nestedblock--match"></a>
### Nested Schema for `match`

A single match expression can have multiple values, for example

```
match {
  host = ["foo.example.com"]
  path = ["/images/*"]
}
```

will only match requests with host with "foo.example.com" AND path "/images/*".
However, multiple match expressions can be chained (acting as an OR statement).
If any matcher matches the request, then the handlers are triggered.

Optional:

- **host** (List of String)


