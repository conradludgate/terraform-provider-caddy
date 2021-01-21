package caddyapi

import "fmt"

// Server represents the Caddy Server object
// https://caddyserver.com/docs/json/apps/http/servers/
type Server struct {
	ID     string         `json:"@id,omitempty"`
	Listen []string       `json:"listen"`
	Routes []Route        `json:"routes"`
	Errors *ServerErrors  `json:"errors"`
	Logs   *ServerLogging `json:"logs"`
}

type ServerErrors struct {
	Routes []Route `json:"routes,omitempty"`
}

type ServerLogging struct {
	DefaultLoggerName string            `json:"default_logger_name,omitempty"`
	LoggerNames       map[string]string `json:"logger_names,omitempty"`
	SkipHosts         []string          `json:"skip_hosts,omitempty"`
	SkipUnmappedHosts bool              `json:"skip_unmapped_hosts,omitempty"`
}

// CreateServer creates a http server
func (c *Client) CreateServer(name string, server Server) (string, error) {
	id := "@config/apps/http/servers/" + name

	if err := c.EnforceExists("@config/apps/http/servers"); err != nil {
		return "", fmt.Errorf("CreateServer [%s]: EnforceExists: %w", name, err)
	}

	resp, err := c.client.R().SetBody(server).Put(URLFromID(id))
	if err != nil {
		return "", fmt.Errorf("CreateServer: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("CreateServer: %w", StatusError{resp})
	}

	if server.ID == "" {
		return id, nil
	}
	return server.ID, nil
}

func (c *Client) UpdateServerListen(id string, listen []string) error {
	resp, err := c.client.R().SetBody(listen).Patch(URLFromID(id + "/listen"))
	if err != nil {
		return fmt.Errorf("UpdateServerListen: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("UpdateServerListen: %w", StatusError{resp})
	}
	return nil
}

func (c *Client) UpdateServerRoutes(id string, routes []Route) error {
	resp, err := c.client.R().SetBody(routes).Patch(URLFromID(id + "/routes"))
	if err != nil {
		return fmt.Errorf("UpdateServerRoutes: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("UpdateServerRoutes: %w", StatusError{resp})
	}
	return nil
}

func (c *Client) GetServer(id string) (*Server, error) {
	resp, err := c.client.R().SetResult(&Server{}).Get(URLFromID(id))
	if err != nil {
		return nil, fmt.Errorf("GetServer: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("GetServer: %w", StatusError{resp})
	}
	return resp.Result().(*Server), nil
}

func (c *Client) DeleteServer(id string) error {
	resp, err := c.client.R().Delete(URLFromID(id))
	if err != nil {
		return fmt.Errorf("DeleteServer: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("DeleteServer: %w", StatusError{resp})
	}
	return nil
}
