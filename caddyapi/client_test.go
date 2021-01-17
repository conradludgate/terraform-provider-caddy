package caddyapi

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/conradludgate/terraform-provider-caddy/caddyapi/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_EnforceExists(t *testing.T) {
	transportMock := &mocks.RoundTripper{}
	c := NewClient(transportMock)

	mockRequest(transportMock, "GET", "http://localhost/config/apps/http", "", "path not found", http.StatusBadRequest).Once()
	mockRequest(transportMock, "GET", "http://localhost/config/apps", "", "null", http.StatusOK).Once()
	mockRequest(transportMock, "POST", "http://localhost/config/apps", `{"http":{}}`, "done", http.StatusOK).Once()

	assert.Nil(t, c.EnforceExists("@config/apps/http"))
}

func Test_EnforceExistsSlice(t *testing.T) {
	transportMock := &mocks.RoundTripper{}
	c := NewClient(transportMock)

	mockRequest(transportMock, "GET", "http://localhost/config/apps/http/servers/server/routes", "", "path not found", http.StatusBadRequest).Once()
	mockRequest(transportMock, "GET", "http://localhost/config/apps/http/servers/server", "", "path not found", http.StatusBadRequest).Once()
	mockRequest(transportMock, "GET", "http://localhost/config/apps/http/servers", "", "null", http.StatusOK).Once()
	mockRequest(transportMock, "POST", "http://localhost/config/apps/http/servers", `{"server":{}}`, "done", http.StatusOK).Once()
	mockRequest(transportMock, "POST", "http://localhost/config/apps/http/servers/server/routes", "[]", "done", http.StatusOK).Once()

	assert.Nil(t, c.EnforceExistsSlice("@config/apps/http/servers/server/routes"))
}

func mockRequest(transportMock *mocks.RoundTripper, method, u string, reqBody string, respBody string, code int) *mock.Call {
	reqMatcher := mock.MatchedBy(func(req *http.Request) bool {
		if req.URL.String() != u {
			return false
		}

		if reqBody == "" {
			if req.Body != nil {
				return false
			}
		} else {
			if req.Header.Get("Content-Type") != "application/json" {
				return false
			}

			if req.Body == nil {
				return false
			}
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return false
			}
			if string(b) != reqBody {
				return false
			}
		}
		return true
	})

	resp := &http.Response{
		StatusCode: code,
		Body:       ioutil.NopCloser(strings.NewReader(respBody)),
	}

	return transportMock.On("RoundTrip", reqMatcher).Return(resp, nil)
}
