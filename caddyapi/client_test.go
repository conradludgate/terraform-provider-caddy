package caddyapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_EnforceExists(t *testing.T) {
	c := NewClient("http://localhost", nil)
	httpmock.ActivateNonDefault(c.client.GetClient())
	defer httpmock.DeactivateAndReset()

	mockRequest("GET", "http://localhost/config/apps/http", "", "path not found", http.StatusBadRequest)
	mockRequest("GET", "http://localhost/config/apps", "", "null", http.StatusOK)
	mockRequest("POST", "http://localhost/config/apps", `{"http":{}}`, "done", http.StatusOK)

	assert.NoError(t, c.EnforceExists("@config/apps/http"))
	callinfo := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, callinfo["GET http://localhost/config/apps/http"])
	assert.Equal(t, 1, callinfo["GET http://localhost/config/apps"])
	assert.Equal(t, 1, callinfo["POST http://localhost/config/apps"])
}

func mockRequest(method, u string, reqBody string, respBody string, code int) {
	httpmock.RegisterResponder(method, u, func(req *http.Request) (*http.Response, error) {
		if reqBody != "" {
			defer req.Body.Close()
			if req.Header.Get("Content-Type") != "application/json" {
				return nil, fmt.Errorf("unexpected content-type")
			}
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			if reqBody != string(b) {
				return nil, fmt.Errorf("unexpected body %s", string(b))
			}
		}

		return httpmock.NewStringResponse(code, respBody), nil
	})
}
