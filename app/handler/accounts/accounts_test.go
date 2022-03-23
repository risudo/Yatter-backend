package accounts_test

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

func TestMain(m *testing.M) {
	c = setup()
	defer c.Close()
	code := m.Run()

	os.Exit(code)
}

//TODO: エンドポイントごとにテストする
func TestAccountRegistration(t *testing.T) {
	username := "testuser1"

	func() {
		resp, err := c.PostJSON("/v1/accounts", fmt.Sprintf(`{"username":"%s"}`, username))
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, http.StatusOK, resp.StatusCode)  {
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, username, j["username"])
		}
	}()

	func() {
		resp, err := c.Get(fmt.Sprintf("/v1/accounts/%s", username))
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, http.StatusOK, resp.StatusCode)  {
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, username, j["username"])
		}
	}()
}

// 2回account作成
func TestAccountRegistrationDupricate(t *testing.T) {
	username := "testuser2"

	func() {
		resp, err := c.PostJSON("/v1/accounts", fmt.Sprintf(`{"username":"%s"}`, username))
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, http.StatusOK, resp.StatusCode)  {
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, username, j["username"])
		}
	}()

	func() {
		resp, err := c.PostJSON("/v1/accounts", fmt.Sprintf(`{"username":"%s"}`, username))
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, http.StatusConflict, resp.StatusCode)  {
			return
		}
	}()
}

func TestGetNosuchAccount(t *testing.T) {
	username := "nosuchusername"

	resp, err := c.Get(fmt.Sprintf("/v1/accounts/%s", username))
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, http.StatusNotFound, resp.StatusCode)  {
		return
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
