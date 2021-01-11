terraform {
  required_providers {
    caddy = {
      version = "~> 0.1.0"
      source  = "conradludgate/caddy"
    }
  }
}

provider "caddy" {
  host            = "ssh://terraform@ssh.conradludgate.com:22/tmp/caddy-admin.sock"
  ssh_key         = "/home/oon/.ssh/terraform"
  ignore_host_key = true
}

resource "caddy_http" "http" {
}

resource "caddy_server" "foo" {
  http = caddy_http.http.id

  name   = "foo"
  listen = [":443"]
}

resource "caddy_server_route" "route1" {
  server = caddy_server.foo.id
}

resource "caddy_server_route_match" "match1" {
  route = caddy_server_route.route1.id

  host = ["example1.conradludgate.com"]
}

resource "caddy_server_route_handle" "handler1" {
  route = caddy_server_route.route1.id

  handler = "reverse_proxy"

  upstream {
    dial = "localhost:8080"
  }
}
