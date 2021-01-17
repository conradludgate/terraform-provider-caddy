# terraform-provider-caddy

This is a terraform provider to manage the [caddy api](https://caddyserver.com/).

## Setup

First you will need caddy running with the admin api enabled. This provider supports two methods to accesing the API endpoint

### HTTP Endpoint

The simplest is to just use the default endpoint `http://localhost:2019`.
This is the default caddy uses and the default that this provider will use too.

However, this is not recommended as it is not secure.

### Unix Sockets

The recommended method is to use unix sockets. Modify your caddy config to use the admin endpoint `unix//path/to/admin.sock`.
Next, run caddy (preferrably with `-resume` and in the background. If you have systemd, check out the caddy-api service file).

Once caddy is running, it should create the unix socket at the path specified. Test it by running the following (making sure you have permission to access the socket)

```sh
curl -H "Host: " --unix-socket /path/to/admin.sock http://localhost/config/
```

If you don't get an error, all is good to go! Finally, set up the provider like so

```tf
provider "caddy" {
  host = "unix:///path/to/admin.sock
}
```

### SSH

In addition to using any of the two above methods to connect to the admin API, you can proxy the request through SSH to ensure authorized access over the internet.
Ensure you can SSH into the server where Caddy is running, and that user can access the admin endpoint, the add the following to your provider config

```tf
provider "caddy" {
  host = "unix:///path/to/admin.sock
  ssh = {
    host = "user@example.com:22" // port is required
    key_file = "~/.ssh/id_rsa" // or specify a password in the host field 'user:pass@example.com:22'
    host_key = "example.com ssh_rsa ...." // in 'known_hosts' format.
  }
}
```
