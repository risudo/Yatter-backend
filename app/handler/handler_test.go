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
	"reflect"
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

const notExistingUser = "John"
const content = "hogehoge"

// 1 follows 2
var existing1 = &object.Account{
	ID:       1,
	Username: "john",
}

var existing2 = &object.Account{
	ID:       2,
	Username: "sum",
}

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
			name: "Create",
			request: func(m *C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"username":"%s"}`, notExistingUser)))
				req, err := http.NewRequest("POST", m.asURL("/v1/accounts"), body)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectUsername:   notExistingUser,
		},
		{
			name: "Fetch",
			request: func(m *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", m.asURL(fmt.Sprintf("/v1/accounts/%s", existing1.Username)), nil)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectUsername:   existing1.Username,
		},
		{
			name: "CreateDupricatedUsername",
			request: func(m *C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"username":"%s"}`, existing1.Username)))
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
			name: "FetchNotExist",
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
			name: "UnauthorizePost",
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
			name: "Post",
			request: func(c *C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"status":"%s"}`, content)))
				req, err := http.NewRequest("POST", c.asURL("/v1/statuses"), body)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing1.Username))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectContent:    content,
		},
		{
			name: "Fetch",
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
			name: "FetchNotExist",
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
			name: "Delete",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("DELETE", c.asURL("/v1/statuses/1"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing1.Username))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
		},
		{
			name: "DeleteNotExist",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("DELETE", c.asURL("/v1/statuses/10"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing1.Username))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusNotFound,
		},
		{
			name: "DeleteNotOwn",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("DELETE", c.asURL("/v1/statuses/1"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing2.Username))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusBadRequest,
		},
		{
			name: "UnauthorizeDelete",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("DELETE", c.asURL("/v1/statuses/1"), nil)
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

				var j map[string]interface{}
				if assert.NoError(t, json.Unmarshal(body, &j)) {
					assert.Equal(t, tt.expectContent, j["content"])
				}
			}
		})
	}
}

func TestTimeline(t *testing.T) {
	m := mockSetup()
	defer m.Close()

	tests := []struct {
		name             string
		request          func(c *C) (*http.Response, error)
		expectStatusCode int
		expectContent    string
	}{
		{
			name: "Public",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.asURL("/v1/timelines/public"), nil)
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
			name: "MoreThanMinLimitPublic",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.asURL("/v1/timelines/public"), nil)
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
			expectContent:    content,
		},
		{
			name: "LessThanMinLimitPublic",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.asURL("/v1/timelines/public"), nil)
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
			expectContent:    content,
		},
		{
			name: "Home",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.asURL("/v1/timelines/home"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing1.Username))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectContent:    content,
		},
		{
			name: "UnauthorizeHome",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.asURL("/v1/timelines/home"), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusUnauthorized,
			expectContent:    content,
		},
		{
			name: "MoreThanMinLimitHome",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.asURL("/v1/timelines/home"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("limit", "81")
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing1.Username))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusBadRequest,
			expectContent:    content,
		},
		{
			name: "LessThanMinLimitHome",
			request: func(c *C) (*http.Response, error) {
				req, err := http.NewRequest("GET", c.asURL("/v1/timelines/home"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("limit", "-1")
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing1.Username))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusBadRequest,
			expectContent:    content,
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

func TestFollowReturnRelation(t *testing.T) {
	m := mockSetup()
	defer m.Close()

	tests := []struct {
		name               string
		request            func(c *C) (*http.Response, error)
		expectStatusCode   int
		expectRelationWith *object.RelationWith
	}{
		{
			name: "UnauthorizeFollow",
			request: func(c *C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/follow", existing1.Username)
				req, err := http.NewRequest("POST", c.asURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusUnauthorized,
		},
		{
			name: "FollowNotExistAccount",
			request: func(c *C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/follow", notExistingUser)
				req, err := http.NewRequest("POST", c.asURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing1.Username))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusNotFound,
		},
		{
			name: "Follow",
			request: func(c *C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/follow", existing2.Username)
				req, err := http.NewRequest("POST", c.asURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing1.Username))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectRelationWith: &object.RelationWith{
				ID:         existing2.ID,
				Following:  true,
				FollowedBy: false,
			},
		},
		{
			name: "Unfolollow",
			request: func(c *C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/unfollow", existing1.Username)
				req, err := http.NewRequest("POST", c.asURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing2.Username))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectRelationWith: &object.RelationWith{
				ID:         existing1.ID,
				Following:  false,
				FollowedBy: true,
			},
		},
		{
			name: "Relationships",
			request: func(c *C) (*http.Response, error) {
				url := "/v1/accounts/relationships"
				req, err := http.NewRequest("GET", c.asURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("username", existing2.Username)
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing1.Username))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectRelationWith: &object.RelationWith{
				ID:         existing2.ID,
				Following:  true,
				FollowedBy: false,
			},
		},
		{
			name: "UnauthorizeRelationships",
			request: func(c *C) (*http.Response, error) {
				url := "/v1/accounts/relationships"
				req, err := http.NewRequest("GET", c.asURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("username", existing2.Username)
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusUnauthorized,
		},
		{
			name: "NotExistRelationships",
			request: func(c *C) (*http.Response, error) {
				url := "/v1/accounts/relationships"
				req, err := http.NewRequest("GET", c.asURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("username", notExistingUser)
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", existing1.Username))
				return c.Server.Client().Do(req)
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
				j := new(object.RelationWith)
				if assert.NoError(t, json.Unmarshal(body, j)) {
					if !reflect.DeepEqual(j, tt.expectRelationWith) {
						t.Fatal(fmt.Sprintf("mismatch RelationWith:\n\t expect:\t%v\n\t actual:\t%v", tt.expectRelationWith, j))
					}
				}
			}
		})
	}
}

func TestFollowReturnAccounts(t *testing.T) {
	m := mockSetup()
	defer m.Close()

	tests := []struct {
		name             string
		request          func(c *C) (*http.Response, error)
		expectStatusCode int
		expectAccounts   []object.Account
	}{
		{
			name: "following",
			request: func(c *C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/following", existing1.Username)
				req, err := http.NewRequest("GET", c.asURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectAccounts:   []object.Account{*existing2},
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
				var j []object.Account
				if assert.NoError(t, json.Unmarshal(body, &j)) {
					if len(j) > 0 && !reflect.DeepEqual(j[0].Username, tt.expectAccounts[0].Username) {
						t.Fatal(fmt.Sprintf("mismatch Account:\n\t expect:\t%v\n\t actual:\t%v", tt.expectAccounts[0], j[0]))
					}
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
		Username: notExistingUser,
	}, nil
}

func (m *mockdao) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	if username == existing1.Username {
		return &object.Account{
			ID:       1,
			Username: existing1.Username,
		}, nil
	}
	if username == existing2.Username {
		return &object.Account{
			ID:       2,
			Username: existing2.Username,
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
			Account: existing1,
		}, nil
	}
	return nil, nil
}

func (m *mockdao) Delete(ctx context.Context, id object.StatusID) error {
	return nil
}

func (m *mockdao) PublicTimeline(ctx context.Context, p *object.Parameters) (object.Timelines, error) {
	return object.Timelines{
		object.Status{Content: content},
	}, nil
}

func (m *mockdao) HomeTimeline(ctx context.Context, loginID object.AccountID, p *object.Parameters) (object.Timelines, error) {
	return object.Timelines{
		object.Status{Content: content},
	}, nil
}

func (m *mockdao) Follow(ctx context.Context, loginID object.AccountID, targetID object.AccountID) error {
	return nil
}

func (m *mockdao) IsFollowing(ctx context.Context, accountID object.AccountID, targetID object.AccountID) (bool, error) {
	if accountID == 1 && targetID == 2 {
		return true, nil
	}
	return false, nil
}

func (m *mockdao) Following(ctx context.Context, id object.AccountID) ([]object.Account, error) {
	if id == existing1.ID {
		return []object.Account{*existing2}, nil
	}
	return nil, nil
}

func (m *mockdao) Followers(ctx context.Context, id object.AccountID) ([]object.Account, error) {
	if id == existing2.ID {
		return []object.Account{*existing1}, nil
	}
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
