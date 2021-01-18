package caddy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"golang.org/x/crypto/ssh"
)

func providerConfigurer(d *schema.ResourceData) (interface{}, error) {
	dialer = &net.Dialer{}
	if sshSet := GetObjectSet(d, "ssh"); len(sshSet) == 1 {
		var err error
		dialer, err = parseSSHConfig(&sshSet[0])
		if err != nil {
			return nil, err
		}
	}

	host, err := url.Parse(GetString(d, "host"))
	if err != nil {
		return nil, err
	}

	if host.Scheme == "unix" {
		dialer = &unixDialer{dialer, host.Path}
	}

	return caddyapi.NewClient(caddyTransport{
		host,
		&http.Transport{Dial: dialer.Dial},
	}), nil
}

// Dialer is a net.Conn factory
type Dialer interface {
	Dial(network, addr string) (net.Conn, error)
}

var dialer Dialer

type unixDialer struct {
	base   Dialer
	socket string
}

func (unix *unixDialer) Dial(_, _ string) (net.Conn, error) {
	return unix.base.Dial("unix", unix.socket)
}

var conns []io.Closer

// CloseConns closes any remaining open connections
func CloseConns() error {
	for i := len(conns) - 1; i >= 0; i-- {
		if err := conns[i].Close(); err != nil {
			return err
		}
		conns = conns[:i]
	}
	return nil
}

func parseSSHHost(host string) (username, password, addr string) {
	if i := strings.Index(host, "@"); i != -1 {
		username = host[:i]
		addr = host[i+1:]
	}

	if i := strings.Index(username, ":"); i != -1 {
		password = username[:i]
		username = username[i+1:]
	}

	return
}

func parseSSHConfig(d *MapData) (*ssh.Client, error) {
	user, pass, addr := parseSSHHost(GetString(d, "host"))

	config := &ssh.ClientConfig{User: user}
	if pass != "" {
		config.Auth = append(config.Auth, ssh.Password(pass))
	} else {
		b, err := ioutil.ReadFile(GetString(d, "key_file"))
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(b)
		if err != nil {
			return nil, err
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	knownHost := []byte(GetString(d, "host_key"))
	_, _, hostKey, _, rest, err := ssh.ParseKnownHosts(knownHost)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("bytes leftover while parsing known_host: %s", string(rest))
	}
	config.HostKeyCallback = ssh.FixedHostKey(hostKey)

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	conns = append(conns, client)

	return client, nil
}

type caddyTransport struct {
	host *url.URL
	base http.RoundTripper
}

func (ct caddyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = ct.host.Scheme
	r.URL.User = ct.host.User
	r.URL.Host = ct.host.Host
	if r.URL.Host == "" {
		// https://github.com/caddyserver/caddy/blob/59071ea15d2aacb69fcfc088f4996717cd2bfc73/cmd/commandfuncs.go#L720-L735
		r.URL.Host = " "
	}
	r.Host = ct.host.Host
	return ct.base.RoundTrip(r)
}
