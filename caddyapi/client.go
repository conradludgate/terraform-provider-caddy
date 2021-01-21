package caddyapi

import (
	"fmt"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/go-resty/resty/v2"
)

// URLFromID Converts an ID into it's resource URL
// "@config/apps/http" -> "/config/apps/http"
// "my_server/routes/0" -> "/id/my_server/routes/0"
func URLFromID(id string) string {
	if strings.HasPrefix(id, "@config") {
		return "/" + id[1:]
	}
	return "/id/" + id
}

// Client represents a Caddy API Client
type Client struct {
	client *resty.Client
}

type DialFunc = func(network, addr string) (net.Conn, error)

func NewClient(host string, dial DialFunc) *Client {
	return &Client{
		client: resty.New().SetTransport(&http.Transport{Dial: dial}).SetHostURL(host),
	}
}

func NewUnixClient(socket string, dial DialFunc) *Client {
	if dial == nil {
		dial = net.Dial
	}
	transport := http.Transport{
		Dial: func(_, _ string) (net.Conn, error) {
			return dial("unix", socket)
		},
	}
	return &Client{
		client: resty.New().SetTransport(&transport).SetScheme("http").SetHostURL(""),
	}
}

// EnforceExists ensures that the path exists. If it doesn't currently exist, it sets it to be empty
func (c *Client) EnforceExists(id string) error {
	split := strings.Split(id, "/")
	for i := len(split); i > 1; i-- {
		resp, err := c.client.R().Get(URLFromID(path.Join(split[:i]...)))
		if err != nil {
			return err
		}
		if resp.IsError() {
			continue
		}

		body := resp.Body()
		if string(body) == "null" {
			data := makeEmptyObject(split[i:])
			resp, err = c.client.R().SetBody(data).Post(URLFromID(path.Join(split[:i]...)))
			if err != nil {
				return err
			}
			if resp.IsError() {
				return StatusError{resp}
			}
		}

		return nil
	}
	return nil
}

func makeEmptyObject(split []string) map[string]interface{} {
	data := make(map[string]interface{})
	if len(split) != 0 {
		data[split[0]] = makeEmptyObject(split[1:])
	}
	return data
}

// StatusError is the error returned when Request responds with non-2XX
type StatusError struct {
	resp *resty.Response
}

func (s StatusError) Error() string {
	str := fmt.Sprintf("%s %s: %d %s", s.resp.Request.Method, s.resp.Request.URL, s.resp.StatusCode(), s.resp.Status())
	if body := s.resp.Body(); len(body) > 0 {
		return fmt.Sprintf("%s - body: [%s]", str, string(body))
	}
	return str
}
