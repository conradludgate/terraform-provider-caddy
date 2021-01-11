package caddyapi

import "fmt"

// Server represents the Caddy Server object
// https://caddyserver.com/docs/json/apps/http/servers/
type Server struct {
	ID     string   `json:"@id,omitempty"`
	Listen []string `json:"listen,omitempty"`
	Routes []Route  `json:"routes,omitempty"`
}

// CreateServer creates a http server
func (c *Client) CreateServer(name string, server Server) (string, error) {
	id := "@config/apps/http/servers/" + name

	if err := c.EnforceExists("@config/apps/http/servers"); err != nil {
		return "", fmt.Errorf("CreateServer [%s]: EnforceExists: %w", name, err)
	}

	if err := c.Request("PUT", id, server, nil); err != nil {
		return "", fmt.Errorf("CreateServer [%s]: Request Put: %w", name, err)
	}

	if server.ID == "" {
		return id, nil
	}
	return server.ID, nil
}

func (c *Client) UpdateServerListen(id string, listen []string) error {
	if err := c.Request("PATCH", id+"/listen", listen, nil); err != nil {
		return fmt.Errorf("UpdateServerListen [%s]: %w", id, err)
	}
	return nil
}

func (c *Client) GetServer(id string) (Server, error) {
	var server Server
	if err := c.Request("GET", id, nil, &server); err != nil {
		return server, fmt.Errorf("GetServer [%s]: %w", id, err)
	}
	return server, nil
}

func (c *Client) DeleteServer(id string) error {
	if err := c.Request("DELETE", id, nil, nil); err != nil {
		return fmt.Errorf("DeleteServer [%s]: %w", id, err)
	}
	return nil
}
