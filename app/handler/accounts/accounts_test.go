package accounts_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler"

	"github.com/stretchr/testify/assert"
)

func TestAccountRegistration(t *testing.T) {
	c := setup(t)
	defer c.Close()
	username := "testuser1"

	func() {
		resp, err := c.PostJSON("/v1/accounts", fmt.Sprintf(`{"username":"%s"}`, username))
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, fmt.Sprintf("%s", username), j["username"])
		}
	}()

	func() {
		resp, err := c.Get(fmt.Sprintf("/v1/accounts/%s", username))
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, fmt.Sprintf("%s", username), j["username"])
		}
	}()
}

// 2回account作成
func TestAccountRegistrationDupricate(t *testing.T) {
	c := setup(t)
	defer c.Close()
	username := "testuser1"

	func() {
		resp, err := c.PostJSON("/v1/accounts", fmt.Sprintf(`{"username":"%s"}`, username))
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, fmt.Sprintf("%s", username), j["username"])
		}
	}()

	func() {
		resp, err := c.PostJSON("/v1/accounts", `{"username":"testuser2"}`)
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if err := json.Unmarshal(body, &j); err != nil {
			t.Fatal(err)
		}
		if j["username"] == fmt.Sprintf("%s", username) {
			t.Fatal(err)
		}
	}()
}

func setup(t *testing.T) *C {
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Dao.InitAll(); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(handler.NewRouter(app))

	return &C{
		App:    app,
		Server: server,
	}
}

type C struct {
	App    *app.App
	Server *httptest.Server
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
