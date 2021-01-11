package caddyapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

// URLFromID Converts an ID into it's resource URL
// "@config/apps/http" -> "/config/apps/http"
// "my_server/routes/0" -> "/id/my_server/routes/0"
func URLFromID(id string) string {
	if strings.HasPrefix(id, "@config") {
		return "http://localhost/" + id[1:]
	}
	return "http://localhost/id/" + id
}

//go:generate mockery --name=HTTPClient
// HTTPClient is an interface over http.Client.Do
type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

// Client represents a Caddy API Client
type Client struct {
	HTTPClient HTTPClient
}

// Get performs a GET request
func (c *Client) Get(id string, respBody interface{}) error {
	return c.Request("GET", id, nil, respBody)
}

// Request performs a generic HTTP request
func (c *Client) Request(method, id string, body, respBody interface{}) error {
	log.Println("HTTPRequest:", method, id, body, respBody)

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("Caddy Request [%s %s]. Could not encode body: %w", method, id, err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, URLFromID(id), bodyReader)
	if err != nil {
		return fmt.Errorf("Caddy Request [%s %s]. Could not create request: %w", method, id, err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Caddy Request [%s %s]. Could not read response: %w", method, id, err)
	}

	if resp.StatusCode/100 != 2 {
		return StatusError{resp.StatusCode, b}
	}

	if respBody != nil {
		if err := json.Unmarshal(b, respBody); err != nil {
			return fmt.Errorf("Caddy Request [%s %s]. Could not decode response: %w", method, id, err)
		}
	}
	return nil
}

// EnforceExists ensures that the path exists. If it doesn't currently exist, it sets it to be empty
func (c *Client) EnforceExists(id string) error {
	log.Println("EnforceExists", id)

	split := strings.Split(id, "/")
	for i := len(split); i > 1; i-- {
		var data interface{}
		err := c.Get(path.Join(split[:i]...), &data)
		if err != nil {
			var statusError StatusError
			if errors.As(err, &statusError) {
				continue
			} else {
				return err
			}
		}

		if data == nil {
			data := makeEmptyObject(split[i:])
			return c.Request("POST", path.Join(split[:i]...), data, nil)
		}

		return nil
	}
	return nil
}

// EnforceExistsSlice ensures that the path exists. If it doesn't currently exist, it sets it to be an empty slice
func (c *Client) EnforceExistsSlice(id string) error {
	log.Println("EnforceExistsSlice", id)

	var data interface{}
	err := c.Get(id, &data)
	if err != nil {
		var statusError StatusError
		if errors.As(err, &statusError) {
			split := strings.Split(id, "/")
			if err := c.EnforceExists(path.Join(split[:len(split)-1]...)); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if data == nil {
		return c.Request("POST", id, []string{}, nil)
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
	StatusCode int
	Body       []byte
}

func (s StatusError) Error() string {
	if len(s.Body) == 0 {
		return fmt.Sprintf("%d: %s", s.StatusCode, http.StatusText(s.StatusCode))
	}
	return fmt.Sprintf("%d: %s - body: [%s]", s.StatusCode, http.StatusText(s.StatusCode), string(s.Body))
}
