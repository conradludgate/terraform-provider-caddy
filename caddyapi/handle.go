package caddyapi

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

// Handle represents the Caddy Handle object
// https://caddyserver.com/docs/json/apps/http/servers/routes/handle/
type Handle struct {
	ID      string `json:"@id,omitempty"`
	Handler string `json:"handler"`

	// Static Response
	Body       *string      `json:"body,omitempty"`
	StatusCode *int         `json:"status_code,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`

	// Reverse Proxy
	Upstreams []Upstream `json:"upstreams,omitempty"`

	// ...
}

type Upstream struct {
	Dial string `json:"dial"`
}

// CreateServerRouteHandle creates a http server route handler
func (c *Client) CreateServerRouteHandle(routeID string, handle Handle) (string, error) {
	log.Println("CreateServerRouteHandle:", routeID, handle)
	if handle.ID == "" {
		// All handlers should have an ID set in caddy because
		// deleting elements in an array breaks the structured IDs
		handle.ID = uuid.New().String()
	}

	handlesID := routeID + "/handle"

	if err := c.EnforceExistsSlice(handlesID); err != nil {
		return "", err
	}

	var handles []Handle
	err := c.Get(handlesID, &handles)
	if err != nil {
		return "", err
	}

	if err := c.Request("POST", handlesID+"/...", []Handle{handle}, nil); err != nil {
		return "", err
	}

	return handle.ID, nil
}

func (c *Client) UpdateServerRouteHandle(id string, handle Handle) error {
	return c.Request("PATCH", id, handle, nil)
}

func (c *Client) GetServerRouteHandle(id string) (Handle, error) {
	log.Println("GetServerRouteHandle:", id)
	var handle Handle
	return handle, c.Request("GET", id, nil, &handle)
}

func (c *Client) DeleteServerRouteHandle(id string) error {
	return c.Request("DELETE", id, nil, nil)
}
