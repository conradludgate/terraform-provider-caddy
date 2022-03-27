terraform {
  required_providers {
    caddy = {
      version = "~> 0.2.0"
      source  = "conradludgate/caddy"
    }
  }
}

provider "caddy" {
  host            = "unix:///tmp/caddy-admin.sock"
  ssh {
    host = "terraform@ssh.conradludgate.com:22"
    key_file = "/home/oon/.ssh/terraform"
  }
}

resource "caddy_server" "foo" {
  name   = "foo"
  listen = [":443"]

  routes = [
    data.caddy_server_route.route1.id,
  ]
}

data "caddy_server_route" "route1" {
  match {
    host = ["example1.conradludgate.com"]
  }

  handle {
    reverse_proxy {
      upstream {
        dial = "localhost:8080"
      }
    }
  }
}

resource "caddy_server_route" "route2" {
  server_name = "foo"
  route_id    = "route2"
  match {
    host = ["example2.conradludgate.com"]
  }

  handle {
    reverse_proxy {
      upstream {
        dial = "localhost:8082"
      }
    }
  }
  depends_on = [ caddy_server.foo ]
}
