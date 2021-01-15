package caddy

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	// "github.com/kevinburke/ssh_config"
	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"golang.org/x/crypto/ssh"
)

func providerConfigurer(d *schema.ResourceData) (interface{}, error) {
	host, err := url.Parse(d.Get("host").(string))
	if err != nil {
		return nil, err
	}
	if host.Scheme == "ssh" {
		conn, err = newSSHConn(host, d.Get("ssh_key").(string), d.Get("host_key").(string))
		if err != nil {
			return nil, err
		}

		return caddyapi.NewClient(unixTransport{&http.Transport{Dial: conn.Dial}}), nil
	}
	return caddyapi.NewClient(caddyTransport{host}), nil
}

var conn *sshConn

type sshConn struct {
	sshClient *ssh.Client
	socket    string
}

func (c *sshConn) Close() error {
	return c.sshClient.Close()
}

// CloseConn closes the SSH connection
func CloseConn() error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func newSSHConn(host *url.URL, pkFile string, hkFile string) (*sshConn, error) {
	config := &ssh.ClientConfig{
		User: host.User.Username(),
	}

	if pw, ok := host.User.Password(); ok {
		config.Auth = append(config.Auth, ssh.Password(pw))
	}

	if pkFile != "" {
		// Read ssh private key
		b, err := ioutil.ReadFile(pkFile)
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(b)
		if err != nil {
			return nil, err
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	if hkFile == "" {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	} else {
		b, err := ioutil.ReadFile(hkFile)
		if err != nil {
			return nil, err
		}
		hostKey, err := ssh.ParsePublicKey(b)
		if err != nil {
			return nil, err
		}
		config.HostKeyCallback = ssh.FixedHostKey(hostKey)
	}

	// Connect to ssh server
	sshClient, err := ssh.Dial("tcp", host.Host, config)
	if err != nil {
		return nil, err
	}

	return &sshConn{sshClient, host.Path}, nil
}

func (c *sshConn) Dial(_, _ string) (net.Conn, error) {
	return c.sshClient.Dial("unix", c.socket)
}

type unixTransport struct {
	base http.RoundTripper
}

func (ut unixTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	// hack to have no host
	// https://github.com/caddyserver/caddy/blob/59071ea15d2aacb69fcfc088f4996717cd2bfc73/cmd/commandfuncs.go#L720-L735
	r.URL.Host = " "
	r.Host = ""
	return ut.base.RoundTrip(r)
}

type caddyTransport struct {
	host *url.URL
}

func (ct caddyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = ct.host.Scheme
	r.URL.User = ct.host.User
	r.URL.Host = ct.host.Host
	return http.DefaultTransport.RoundTrip(r)
}
