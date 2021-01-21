---
page_title: "caddy_server Resource - terraform-provider-caddy"
subcategory: ""
description: |-
  
---

# Resource `caddy_server`

Resource to manage a [caddy HTTP server](https://caddyserver.com/docs/json/apps/http/servers/).

## Example

```tf
resource "caddy_server" "https" {
  name   = "https"
  listen = [":443"]

  route {
    match {
      host = "foo.example.com"
    }

    handler {
      static_response {
        body = "Hello World!"
      }
    }
  }

  route {
    match {
      host = "bar.example.com"
    }

    handler {
      reverse_proxy {
        upstream {
          dial = "localhost:2020"
        }
      }
    }
  }
}
```

Or in separated form


```tf
resource "caddy_server" "https" {
  name   = "https"
  listen = [":443"]

  routes = [
    data.caddy_server_route.foo.id,
    data.caddy_server_route.bar.id,
  ]
}

data "caddy_server_route" "foo" {
  match {
    host = "foo.example.com"
  }

  handler {
    static_response {
      body = "Hello World!"
    }
  }
}

data "caddy_server_route" "bar" {
  match {
    host = "bar.example.com"
  }

  handler {
    reverse_proxy {
      upstream {
        dial = "localhost:2020"
      }
    }
  }
}
```

## Schema

### Required

- **listen** (List of String) List of ports to listen on
- **name** (String) Name of the server

### Optional

Either `route` or `routes` can be provided. `route` is for nested schema, whereas `routes` is for separated schema.

- **route** (Block List) (see [below for nested schema](#nestedblock--route))
- **routes** (List of String)

If an error occurred during the handling of a request, the can get processed through this list of error routes. Similar to `routes`,
either `error` or `errors` can be provided. `error` is for nested schema, whereas `errors` is for separated schema.

- **error** (Block List) (see [below for nested schema](#nestedblock--route))
- **errors** (List of String)

- **logs** (Block Set) (see [below for nested schema](#nestedblock--logs))

<a id="nestedblock--route"></a>
### Nested Schema for `route`

Required:

- **group** (String)
- **terminal** (Boolean)
- **handle** (Block List) (see [below for nested schema](#nestedblock--route--handle))
- **match** (Block List) (see [below for nested schema](#nestedblock--route--match))

<a id="nestedblock--route--handle"></a>
### Nested Schema for `route.handle`

Optional:

Only one of these sets can be specified in a single handle. However, a route can have multiple handles

- **reverse_proxy** (Block Set) (see [below for nested schema](#nestedblock--route--handle--reverse_proxy))
- **static_response** (Block Set) (see [below for nested schema](#nestedblock--route--handle--static_response))

<a id="nestedblock--route--handle--reverse_proxy"></a>
### Nested Schema for `route.handle.reverse_proxy`

Reverse proxy takes the request and forwards it to one of the given upstreams.
https://caddyserver.com/docs/json/apps/http/servers/routes/handle/reverse_proxy/

Optional:

- **upstream** (Block List) (see [below for nested schema](#nestedblock--route--handle--reverse_proxy--upstream))

<a id="nestedblock--route--handle--reverse_proxy--upstream"></a>
### Nested Schema for `route.handle.reverse_proxy.upstream`

Optional:

- **dial** (String) The network address to dial to connect to the upstream. Must represent precisely one socket (i.e. no port ranges). A valid network address either has a host and port or is a unix socket address.


<a id="nestedblock--route--route--handle--static_response"></a>
### Nested Schema for `route.handle.static_response`

Optional:

- **body** (String)
- **close** (Boolean) If true, the server will close the client's connection after writing the response.
- **headers** (Map of String)
- **status_code** (String)



<a id="nestedblock--route--match"></a>
### Nested Schema for `route.match`

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



<a id="nestedblock--logs"></a>
### Nested Schema for `logs`

The logs set is for configuring logs emitted by caddy. Include an empty `logs` set if you want a minimal logging configuration

Optional:

- **default_logger_name** (String) The default name for caddy to use when logging requests
- **logger_names** (Map of String) Map the request host to a new logger name
- **skip_hosts** (List of String) Skip logging for any requests using these hosts
- **skip_unmapped_hosts** (Boolean) Skip logging any request that does not appear in `logger_names`
