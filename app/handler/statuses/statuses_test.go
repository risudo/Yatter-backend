package statuses_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/handler_test_setup"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	m := handler_test_setup.MockSetup()
	defer m.Close()

	tests := []struct {
		name             string
		request          func(c *handler_test_setup.C) (*http.Response, error)
		expectStatusCode int
		expectContent    string
	}{
		{
			name: "UnauthorizePost",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"status":"%s"}`, handler_test_setup.Content)))
				req, err := http.NewRequest("POST", c.AsURL("/v1/statuses"), body)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Post",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"status":"%s"}`, handler_test_setup.Content)))
				req, err := http.NewRequest("POST", c.AsURL("/v1/statuses"), body)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername1))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectContent:    handler_test_setup.Content,
		},
		{
			name: "Fetch",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.AsURL("/v1/statuses/1"), nil)
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
			name: "FetchNotExist",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.AsURL("/v1/statuses/100"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusNotFound,
		},
		{
			name: "Delete",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("DELETE", c.AsURL("/v1/statuses/1"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername1))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectContent:    "",
		},
		{
			name: "DeleteNotExist",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("DELETE", c.AsURL("/v1/statuses/10"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername1))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusNotFound,
		},
		{
			name: "DeleteNotOwn",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("DELETE", c.AsURL("/v1/statuses/1"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername2))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusBadRequest,
		},
		{
			name: "UnauthorizeDelete",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("DELETE", c.AsURL("/v1/statuses/1"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusUnauthorized,
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

				var j object.Status
				if assert.NoError(t, json.Unmarshal(body, &j)) {
					assert.Equal(t, tt.expectContent, j.Content)
				}
			}
		})
	}
}
