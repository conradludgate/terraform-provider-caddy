---
page_title: "caddy_http Resource - terraform-provider-caddy"
subcategory: ""
description: |-
  
---

# Resource `caddy_http`





## Schema

### Optional

- **grace_period** (String) How long to wait for active connections when shutting down the server. Once the grace period is over, connections will be forcefully closed.
					Duration can be an integer or a string. An integer is interpreted as nanoseconds. If a string, it is a Go time.Duration value such as 300ms, 1.5h, or 2h45m; valid units are ns, us/µs, ms, s, m, h, and d.
- **http_port** (Number) specifies the port to use for HTTP (as opposed to HTTPS), which is used when setting up HTTP->HTTPS redirects or ACME HTTP challenge solvers. Default: 80.
- **https_port** (Number) specifies the port to use for HTTPS, which is used when solving the ACME TLS-ALPN challenges, or whenever HTTPS is needed but no specific port number is given. Default: 443.
- **id** (String) The ID of this resource.


