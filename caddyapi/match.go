package caddyapi

import (
	"log"

	"github.com/google/uuid"
)

// Match represents the Caddy Match object
// https://caddyserver.com/docs/json/apps/http/servers/routes/match/
type Match struct {
	ID     string              `json:"@id,omitempty"`
	Host   []string            `json:"host,omitempty"`
	Path   []string            `json:"path,omitempty"`
	Method []string            `json:"method,omitempty"`
	Header map[string][]string `json:"header,omitempty"`
	Query  map[string][]string `json:"query,omitempty"`
	Not    []Match             `json:"not,omitempty"`
	// ...
}

// CreateServerRouteMatch creates a http server route matcher
func (c *Client) CreateServerRouteMatch(routeID string, match Match) (string, error) {
	log.Println("CreateServerRouteMatch:", routeID, match)
	if match.ID == "" {
		// All matchers should have an ID set in caddy because
		// deleting elements in an array breaks the structured IDs
		match.ID = uuid.New().String()
	}

	matchesID := routeID + "/match"

	if err := c.EnforceExistsSlice(matchesID); err != nil {
		return "", err
	}

	var matches []Match
	err := c.Get(matchesID, &matches)
	if err != nil {
		return "", err
	}

	if err := c.Request("POST", matchesID+"/...", []Match{match}, nil); err != nil {
		return "", err
	}

	return match.ID, nil
}

func (c *Client) UpdateServerRouteMatch(id string, match Match) error {
	return c.Request("PATCH", id, match, nil)
}

func (c *Client) GetServerRouteMatch(id string) (Match, error) {
	log.Println("GetServerRouteMatch:", id)
	var match Match
	return match, c.Request("GET", id, nil, &match)
}

func (c *Client) DeleteServerRouteMatch(id string) error {
	return c.Request("DELETE", id, nil, nil)
}
