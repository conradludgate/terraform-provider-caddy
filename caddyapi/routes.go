package caddyapi

import (
	"log"

	"github.com/google/uuid"
)

// Route represents the Caddy Route object
// https://caddyserver.com/docs/json/apps/http/servers/routes/
type Route struct {
	ID       string                   `json:"@id,omitempty"`
	Group    string                   `json:"group,omitempty"`
	Matchers []Match                  `json:"match,omitempty"`
	Handlers []map[string]interface{} `json:"handle,omitempty"`
	Terminal bool                     `json:"terminal,omitempty"`
}

// CreateServerRoute creates a http server route
func (c *Client) CreateServerRoute(serverID string, route Route) (string, error) {
	log.Println("CreateServerRoute:", serverID, route)
	if route.ID == "" {
		// All routes should have an ID set in caddy because
		// deleting elements in an array breaks the structured IDs
		route.ID = uuid.New().String()
	}
	routesID := serverID + "/routes"

	if err := c.EnforceExistsSlice(routesID); err != nil {
		return "", err
	}

	var routes []Route
	err := c.Get(routesID, &routes)
	if err != nil {
		return "", err
	}

	if err := c.Request("POST", routesID+"/...", []Route{route}, nil); err != nil {
		return "", err
	}

	return route.ID, nil
}

func (c *Client) UpdateServerRouteGroup(id string, group string) error {
	return c.Request("PATCH", id+"/group", group, nil)
}
func (c *Client) UpdateServerRouteTerminal(id string, terminal bool) error {
	return c.Request("PATCH", id+"/terminal", terminal, nil)
}

func (c *Client) GetServerRoute(id string) (Route, error) {
	log.Println("GetServerRoute:", id)
	var route Route
	return route, c.Request("GET", id, nil, &route)
}

func (c *Client) DeleteServerRoute(id string) error {
	return c.Request("DELETE", id, nil, nil)
}
