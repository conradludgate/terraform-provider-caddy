package caddyapi

import (
	"fmt"
	"time"
)

// Duration is a wrapper around time.Duration with appropriate text unmarshal func
type Duration time.Duration

// UnmarshalText converts `3.5s` and similar strings into the correct duration value
func (d *Duration) UnmarshalText(s string) error {
	if s == "" {
		return nil
	}
	td, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(td)
	return nil
}

func (d Duration) String() string {
	return time.Duration(d).String()
}

// HTTP is the representation of the Caddy HTTP App
type HTTP struct {
	HTTPPort    int      `json:"http_port"`
	HTTPSPort   int      `json:"https_port"`
	GracePeriod Duration `json:"grace_period"`
	// Servers     []Server `json:"-"`
}

// CreateHTTP creates a new http app
func (c *Client) CreateHTTP(http HTTP) error {
	if err := c.EnforceExists("@config/apps"); err != nil {
		return fmt.Errorf("CreateHTTP: %w", err)
	}

	if err := c.Request("PUT", "@config/apps/http", http, nil); err != nil {
		return fmt.Errorf("UpdateHTTPPort: %w", err)
	}
	return nil
}

func (c *Client) UpdateHTTPPort(httpPort int) error {
	if err := c.Request("PATCH", "@config/apps/http/http_port", httpPort, nil); err != nil {
		return fmt.Errorf("UpdateHTTPPort: %w", err)
	}
	return nil
}

func (c *Client) UpdateHTTPSPort(httpsPort int) error {
	if err := c.Request("PATCH", "@config/apps/http/https_port", httpsPort, nil); err != nil {
		return fmt.Errorf("UpdateHTTPSPort: %w", err)
	}
	return nil
}

func (c *Client) UpdateHTTPGracePeriod(gracePeriod Duration) error {
	if err := c.Request("PATCH", "@config/apps/http/grace_period", gracePeriod, nil); err != nil {
		return fmt.Errorf("UpdateHTTPGracePeriod: %w", err)
	}
	return nil
}

func (c *Client) GetHTTP() (HTTP, error) {
	var http HTTP
	if err := c.Request("GET", "@config/apps/http", nil, &http); err != nil {
		return http, fmt.Errorf("GetHTTP: %w", err)
	}
	return http, nil
}

func (c *Client) DeleteHTTP() error {
	if err := c.Request("DELETE", "@config/apps/http", nil, nil); err != nil {
		return fmt.Errorf("DeleteHTTP: %w", err)
	}
	return nil
}
