package caddyapi

import (
	"encoding/json"
	"fmt"
)

type StaticResponse struct {
	StatusCode string              `json:"status_code,omitempty"`
	Headers    map[string][]string `json:"headers,omitempty"`
	Body       string              `json:"body,omitempty"`
	Close      bool                `json:"close,omitempty"`
}

type ReverseProxy struct {
	Upstreams []Upstream `json:"upstreams"`
}

type Upstream struct {
	Dial string `json:"dial"`
}

type RequestBody struct {
	MaxSize int `json:"max_size"`
}

type FileServer struct {
	Root          string   `json:"root,omitempty"`
	Hide          []string `json:"hide,omitempty"`
	IndexNames    []string `json:"index_names,omitempty"`
	CanonicalURIs bool     `json:"canonical_uris,omitempty"`
	PassThru      bool     `json:"pass_thru,omitempty"`
}

type Templates struct {
	FileRoot   string   `json:"file_root,omitempty"`
	MimeTypes  []string `json:"mime_types,omitempty"`
	Delimiters []string `json:"delimiters,omitempty"`
}

type HandleMarshal struct {
	Handle interface{}
}

func ToHandlerName(h interface{}) string {
	switch h.(type) {
	case StaticResponse:
		return "static_response"
	case ReverseProxy:
		return "reverse_proxy"
	case RequestBody:
		return "request_body"
	case FileServer:
		return "file_server"
	case Templates:
		return "templates"
	default:
		return ""
	}
}

func (s *HandleMarshal) FromHandlerName(h string, b []byte) error {
	switch h {
	case "static_response":
		handle := StaticResponse{}
		if err := json.Unmarshal(b, &handle); err != nil {
			return err
		}
		s.Handle = handle
	case "reverse_proxy":
		handle := ReverseProxy{}
		if err := json.Unmarshal(b, &handle); err != nil {
			return err
		}
		s.Handle = handle
	case "request_body":
		handle := RequestBody{}
		if err := json.Unmarshal(b, &handle); err != nil {
			return err
		}
		s.Handle = handle
	case "file_server":
		handle := FileServer{}
		if err := json.Unmarshal(b, &handle); err != nil {
			return err
		}
		s.Handle = handle
	case "templates":
		handle := Templates{}
		if err := json.Unmarshal(b, &handle); err != nil {
			return err
		}
		s.Handle = handle
	default:
		return fmt.Errorf("unsupported handler type %s", h)
	}
	return nil
}

func (s HandleMarshal) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(s.Handle)
	if err != nil {
		return nil, err
	}
	handle := make(map[string]interface{})
	if err := json.Unmarshal(b, &handle); err != nil {
		return nil, err
	}
	handle["handler"] = ToHandlerName(s.Handle)
	if handle["handler"] == "" {
		return nil, fmt.Errorf("unsupported handler %v", s.Handle)
	}
	return json.Marshal(handle)
}

func (s *HandleMarshal) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &object); err != nil {
		return err
	}
	var handler string
	h, ok := object["handler"]
	if !ok {
		return fmt.Errorf("missing required field 'handler'")
	}
	if err := json.Unmarshal([]byte(h), &handler); err != nil {
		return err
	}

	delete(object, "handler")
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}

	return s.FromHandlerName(handler, b)
}
