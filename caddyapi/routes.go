package caddyapi

// Route represents the Caddy Route object
// https://caddyserver.com/docs/json/apps/http/servers/routes/
type Route struct {
	ID       string                   `json:"@id,omitempty"`
	Group    string                   `json:"group,omitempty"`
	Matchers []Match                  `json:"match,omitempty"`
	Handlers []map[string]interface{} `json:"handle,omitempty"`
	Terminal bool                     `json:"terminal,omitempty"`
}
