package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler"

	"github.com/stretchr/testify/assert"
)

type C struct {
	App    *app.App
	Server *httptest.Server
}

var c *C

const preparedUsername1 = "preparedUser1"
const preparedUsername2 = "preparedUser2"
const preparedStatusContent = "prepare"

func TestMain(m *testing.M) {
	c = setup()
	defer c.Close()

	if _, err := c.PostJSON("/v1/accounts", fmt.Sprintf(`{"username":"%s"}`, preparedUsername1)); err != nil {
		panic(err)
	}

	if _, err := c.PostJSON("/v1/accounts", fmt.Sprintf(`{"username":"%s"}`, preparedUsername2)); err != nil {
		panic(err)
	}
	code := m.Run()

	os.Exit(code)
}

func TestAccount(t *testing.T) {
	username := "testuser1"

	tests := []struct {
		name             string
		request          func(c *C) (*http.Response, error)
		expextStatusCode int
	}{
		{
			name: "create account",
			request: func(c *C) (*http.Response, error) {
				return c.PostJSON("/v1/accounts", fmt.Sprintf(`{"username":"%s"}`, username))
			},
			expextStatusCode: http.StatusOK,
		},
		{
			name: "fetch account",
			request: func(c *C) (*http.Response, error) {
				return c.Get(fmt.Sprintf("/v1/accounts/%s", username))
			},
			expextStatusCode: http.StatusOK,
		},
		{
			name: "create dupricated account",
			request: func(c *C) (*http.Response, error) {
				return c.PostJSON("/v1/accounts", fmt.Sprintf(`{"username":"%s"}`, username))
			},
			expextStatusCode: http.StatusConflict,
		},
		{
			name: "create empty username",
			request: func(c *C) (*http.Response, error) {
				return c.PostJSON("/v1/accounts", fmt.Sprintf(`{"username":"%s"}`, ""))
			},
			expextStatusCode: http.StatusBadRequest,
		},
		{
			name: "fetch not exist username",
			request: func(c *C) (*http.Response, error) {
				return c.Get(fmt.Sprintf("/v1/accounts/%s", "NoSuchUser"))
			},
			expextStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.request(c)
			if err != nil {
				t.Fatal(err)
			}
			if !assert.Equal(t, tt.expextStatusCode, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}

				var j map[string]interface{}
				if assert.NoError(t, json.Unmarshal(body, &j)) {
					assert.Equal(t, username, j["username"])
				}
			}
		})
	}
}

func TestStatus(t *testing.T) {
	const content = "ピタ ゴラ スイッチ♪"

	tests := []struct {
		name             string
		request          func(c *C) (*http.Response, error)
		expextStatusCode int
		expextStatusID   int64
	}{
		{
			name: "post unauthorized status",
			request: func(c *C) (*http.Response, error) {
				return c.PostJSON("/v1/statuses", fmt.Sprintf(`{"status": "%s"}`, content))
			},
			expextStatusCode: http.StatusUnauthorized,
		},
		{
			name: "post status",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("POST", c.asURL("/v1/statuses"), bytes.NewReader([]byte(fmt.Sprintf(`{"status": "%s"}`, content))))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", preparedUsername1))
				return c.Server.Client().Do(req)
			},
			expextStatusCode: http.StatusOK,
			expextStatusID:   1,
		},
		{
			name: "fetch status",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.asURL("/v1/statuses/1"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expextStatusCode: http.StatusOK,
			expextStatusID:   1,
		},
		{
			name: "fetch not exist status",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.asURL("/v1/statuses/100"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expextStatusCode: http.StatusNotFound,
		},
		{
			name: "delete status",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("DELETE", c.asURL("/v1/statuses/1"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", preparedUsername1))
				return c.Server.Client().Do(req)
			},
			expextStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.request(c)
			if err != nil {
				t.Fatal(err)
			}
			if !assert.Equal(t, tt.expextStatusCode, resp.StatusCode) {
				return
			}
			if resp.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}

				var j map[string]interface{}
				if assert.NoError(t, json.Unmarshal(body, &j)) {
					assert.Equal(t, content, j["content"])
					assert.EqualValues(t, tt.expextStatusID, j["id"])
				}
			}
		})
	}
}

func setup() *C {
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Dao.InitAll(); err != nil {
		panic(err)
	}

	server := httptest.NewServer(handler.NewRouter(app))

	return &C{
		App:    app,
		Server: server,
	}
}

func (c *C) Close() {
	c.Server.Close()
}

func (c *C) PostJSON(apiPath string, payload string) (*http.Response, error) {
	return c.Server.Client().Post(c.asURL(apiPath), "application/json", bytes.NewReader([]byte(payload)))
}

func (c *C) Get(apiPath string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath))
}

func (c *C) asURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
