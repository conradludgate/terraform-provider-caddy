package caddyapi

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
