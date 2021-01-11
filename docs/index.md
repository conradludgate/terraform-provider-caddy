---
page_title: "caddy Provider"
subcategory: ""
description: |-
  
---

# caddy Provider





## Schema

### Optional

- **host** (String) Caddy Admin API host. Must be a valid URL. If the scheme is `ssh`, then the path value is expected to be a unix socket
- **host_key** (String) SSH Host key file. Only needed if `host` points to an ssh server
- **ignore_host_key** (Boolean) Ignore SSH Host Key
- **ssh_key** (String) SSH Private key file. only needed if `host` points to an ssh server
