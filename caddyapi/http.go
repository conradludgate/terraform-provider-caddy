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

	resp, err := c.client.R().SetBody(http).Put("/config/apps/http")
	if err != nil {
		return fmt.Errorf("CreateHTTP: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("CreateHTTP: %w", StatusError{resp})
	}
	return nil
}

func (c *Client) UpdateHTTPPort(httpPort int) error {
	resp, err := c.client.R().SetBody(httpPort).Patch("/config/apps/http/http_port")
	if err != nil {
		return fmt.Errorf("UpdateHTTPPort: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("UpdateHTTPPort: %w", StatusError{resp})
	}
	return nil
}

func (c *Client) UpdateHTTPSPort(httpsPort int) error {
	resp, err := c.client.R().SetBody(httpsPort).Patch("/config/apps/http/https_port")
	if err != nil {
		return fmt.Errorf("UpdateHTTPSPort: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("UpdateHTTPSPort: %w", StatusError{resp})
	}
	return nil
}

func (c *Client) UpdateHTTPGracePeriod(gracePeriod Duration) error {
	resp, err := c.client.R().SetBody(gracePeriod).Patch("/config/apps/http/grace_period")
	if err != nil {
		return fmt.Errorf("UpdateHTTPGracePeriod: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("UpdateHTTPGracePeriod: %w", StatusError{resp})
	}
	return nil
}

func (c *Client) GetHTTP() (*HTTP, error) {
	resp, err := c.client.R().SetResult(&HTTP{}).Get("/config/apps/http")
	if err != nil {
		return nil, fmt.Errorf("GetHTTP: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("GetHTTP: %w", StatusError{resp})
	}
	return resp.Result().(*HTTP), nil
}

func (c *Client) DeleteHTTP() error {
	resp, err := c.client.R().Delete("/config/apps/http")
	if err != nil {
		return fmt.Errorf("DeleteHTTP: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("DeleteHTTP: %w", StatusError{resp})
	}
	return nil
}
