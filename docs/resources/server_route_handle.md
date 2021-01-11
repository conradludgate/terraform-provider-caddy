---
page_title: "caddy_server_route_handle Resource - terraform-provider-caddy"
subcategory: ""
description: |-
  
---

# Resource `caddy_server_route_handle`





## Schema

### Required

- **handler** (String)
- **route** (String)

### Optional

- **body** (String)
- **headers** (Map of String)
- **id** (String) The ID of this resource.
- **upstream** (Block Set) (see [below for nested schema](#nestedblock--upstream))

<a id="nestedblock--upstream"></a>
### Nested Schema for `upstream`

Optional:

- **dial** (String)


