package caddyapi_test

import (
	"encoding/json"
	"testing"

	"github.com/conradludgate/terraform-provider-caddy/caddyapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleMarshal(t *testing.T) {
	testCases := []struct {
		name     string
		handle   interface{}
		expected string
	}{
		{
			name: "ReverseProxy",
			handle: caddyapi.ReverseProxy{
				Upstreams: []caddyapi.Upstream{
					{
						Dial: "localhost:2020",
					},
				},
			},
			expected: `{
				"handler": "reverse_proxy",
				"upstreams": [
					{
						"dial": "localhost:2020"
					}
				]
			}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(caddyapi.HandleMarshal{tc.handle})
			require.NoError(t, err)

			assert.JSONEq(t, tc.expected, string(b))
		})
	}
}
func TestHandleUnmarshal(t *testing.T) {
	testCases := []struct {
		name     string
		data     string
		expected interface{}
	}{
		{
			name: "ReverseProxy",
			data: `{
				"handler": "reverse_proxy",
				"upstreams": [
					{
						"dial": "localhost:2020"
					}
				]
			}`,
			expected: caddyapi.ReverseProxy{
				Upstreams: []caddyapi.Upstream{
					{
						Dial: "localhost:2020",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handle := caddyapi.HandleMarshal{}
			require.NoError(t, json.Unmarshal([]byte(tc.data), &handle))

			assert.Equal(t, tc.expected, handle.Handle)
		})
	}
}
