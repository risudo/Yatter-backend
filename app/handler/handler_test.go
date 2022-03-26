package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/handler"

	"github.com/stretchr/testify/assert"
)

type (
	C struct {
		App    *app.App
		Server *httptest.Server
	}
	mockdao struct{}
)

const createTestUser = "John"
const fetchTestUser = "joe"
const content = "hogehoge"

func TestAccount(t *testing.T) {
	m := mockSetup()
	defer m.Close()

	tests := []struct {
		name             string
		request          func(c *C) (*http.Response, error)
		expectStatusCode int
		expectUsername   string
	}{
		{
			name: "CreateAccount",
			request: func(m *C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"username":"%s"}`, createTestUser)))
				req, err := http.NewRequest("POST", m.asURL("/v1/accounts"), body)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectUsername:   createTestUser,
		},
		{
			name: "FetchAccount",
			request: func(m *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", m.asURL(fmt.Sprintf("/v1/accounts/%s", fetchTestUser)), nil)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectUsername:   fetchTestUser,
		},
		{
			name: "CreateDupricatedUsername",
			request: func(m *C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"username":"%s"}`, fetchTestUser)))
				req, err := http.NewRequest("POST", m.asURL("/v1/accounts"), body)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusConflict,
		},
		{
			name: "CreateEmptyUsername",
			request: func(m *C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"username":"%s"}`, "")))
				req, err := http.NewRequest("POST", m.asURL("/v1/accounts"), body)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusBadRequest,
		},
		{
			name: "FetchNotExistAccount",
			request: func(m *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", m.asURL(fmt.Sprintf("/v1/accounts/%s", "nosuchuser")), nil)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusNotFound,
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

				var j map[string]interface{}
				if assert.NoError(t, json.Unmarshal(body, &j)) {
					assert.Equal(t, tt.expectUsername, j["username"])
				}
			}
		})
	}
}

func TestStatus(t *testing.T) {
	m := mockSetup()
	defer m.Close()

	tests := []struct {
		name             string
		request          func(c *C) (*http.Response, error)
		expectStatusCode int
		expectContent    string
	}{
		{
			name: "UnauthorizePostStatus",
			request: func(c *C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"status":"%s"}`, content)))
				req, err := http.NewRequest("POST", c.asURL("/v1/statuses"), body)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusUnauthorized,
		},
		{
			name: "PostStatus",
			request: func(c *C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"status":"%s"}`, content)))
				req, err := http.NewRequest("POST", c.asURL("/v1/statuses"), body)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", fetchTestUser))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectContent:    content,
		},
		{
			name: "FetchStatus",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.asURL("/v1/statuses/1"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectContent:    content,
		},
		{
			name: "FetchNotExistStatus",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.asURL("/v1/statuses/100"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusNotFound,
		},
		{
			name: "DeleteStatus",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("DELETE", c.asURL("/v1/statuses/1"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", fetchTestUser))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
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

				var j map[string]interface{}
				if assert.NoError(t, json.Unmarshal(body, &j)) {
					assert.Equal(t, tt.expectContent, j["content"])
				}
			}
		})
	}
}

func (m *mockdao) Account() repository.Account {
	return m
}

func (m *mockdao) Status() repository.Status {
	return m
}

func (m *mockdao) Relation() repository.Relation {
	return m
}

func (m *mockdao) InitAll() error {
	return nil
}

func (m *mockdao) Create(ctx context.Context, a *object.Account) (*object.Account, error) {
	return &object.Account{
		Username: createTestUser,
	}, nil
}

func (m *mockdao) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	if username == fetchTestUser {
		return &object.Account{
			Username: fetchTestUser,
		}, nil
	}
	return nil, nil
}

func (m *mockdao) Post(ctx context.Context, status *object.Status) (*object.Status, error) {
	return &object.Status{
		Content: content,
	}, nil
}

func (m *mockdao) FindByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
	if id == 1 {
		return &object.Status{
			Content: content,
		}, nil
	}
	return nil, nil
}

func (m *mockdao) Delete(ctx context.Context, id object.StatusID) error {
	return nil
}

func (m *mockdao) PublicTimeline(ctx context.Context, p *object.Parameters) (object.Timelines, error) {
	return nil, nil
}

func (m *mockdao) HomeTimeline(ctx context.Context, loginID object.AccountID, p *object.Parameters) (object.Timelines, error) {
	return nil, nil
}

func (m *mockdao) Follow(ctx context.Context, loginID object.AccountID, targetID object.AccountID) error {
	return nil
}

func (m *mockdao) IsFollowing(ctx context.Context, accountID object.AccountID, targetID object.AccountID) (bool, error) {
	return false, nil
}

func (m *mockdao) Following(ctx context.Context, id object.AccountID) ([]object.Account, error) {
	return nil, nil
}

func (m *mockdao) Followers(ctx context.Context, id object.AccountID) ([]object.Account, error) {
	return nil, nil
}

func (m *mockdao) Unfollow(ctx context.Context, loginID object.AccountID, targetID object.AccountID) error {
	return nil
}

func mockSetup() *C {
	app := &app.App{Dao: &mockdao{}}
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
