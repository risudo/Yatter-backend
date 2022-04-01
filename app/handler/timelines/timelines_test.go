package timelines_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/handler_test_setup"

	"github.com/stretchr/testify/assert"
)

func TestTimeline(t *testing.T) {
	m := handler_test_setup.MockSetup()
	defer m.Close()

	tests := []struct {
		name             string
		request          func(c *handler_test_setup.C) (*http.Response, error)
		expectStatusCode int
		expectContent    string
	}{
		{
			name: "Public",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectContent:    handler_test_setup.Content,
		},
		{
			name: "MoreThanMinLimitPublic",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("limit", "81")
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusBadRequest,
			expectContent:    handler_test_setup.Content,
		},
		{
			name: "LessThanMinLimitPublic",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("limit", "-1")
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusBadRequest,
			expectContent:    handler_test_setup.Content,
		},
		{
			name: "Home",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.AsURL("/v1/timelines/home"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername2))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectContent:    handler_test_setup.Content,
		},
		{
			name: "UnauthorizeHome",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.AsURL("/v1/timelines/home"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusUnauthorized,
			expectContent:    handler_test_setup.Content,
		},
		{
			name: "MoreThanMaxLimitHome",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.AsURL("/v1/timelines/home"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("limit", "81")
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername1))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusBadRequest,
		},
		{
			name: "LessThanMinLimitHome",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.AsURL("/v1/timelines/home"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("limit", "-1")
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername1))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.request(m)
			if err != nil {
				t.Fatal(err)
			}
			if !assert.Equal(t, tt.expectStatusCode, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}

				var j object.Timelines
				if assert.NoError(t, json.Unmarshal(body, &j)) {
					if len(j) < 1 {
						t.Fatal("empty timeline")
					}
					assert.Equal(t, tt.expectContent, j[0].Content)
				}
			}
		})
	}
}
