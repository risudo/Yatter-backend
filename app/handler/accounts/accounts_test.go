package accounts_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"testing"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/handler_test_setup"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	m := handler_test_setup.MockSetup()
	defer m.Close()

	tests := []struct {
		name             string
		request          func(c *handler_test_setup.C) (*http.Response, error)
		expectStatusCode int
		expectAccount    *object.Account
	}{
		{
			name: "Create",
			request: func(m *handler_test_setup.C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"username":"%s"}`, handler_test_setup.CreateUser)))
				req, err := http.NewRequest("POST", m.AsURL("/v1/accounts"), body)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectAccount: &object.Account{
				Username: handler_test_setup.CreateUser,
			},
		},
		{
			name: "Fetch",
			request: func(m *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("GET", m.AsURL(fmt.Sprintf("/v1/accounts/%s", handler_test_setup.ExistingUsername1)), nil)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectAccount: &object.Account{
				Username: handler_test_setup.ExistingUsername1,
			},
		},
		{
			name: "Update",
			request: func(m *handler_test_setup.C) (*http.Response, error) {

				filepath := "../../../test/images/image.png"
				file, err := os.Open(filepath)
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()
				body := &bytes.Buffer{}
				mw := multipart.NewWriter(body)
				fw, err := mw.CreateFormFile("avatar", filepath)
				if err != nil {
					t.Fatal(err)
				}
				_, err = io.Copy(fw, file)
				if err != nil {
					t.Fatal(err)
				}
				contentType := mw.FormDataContentType()
				mw.Close()
				req, err := http.NewRequest("POST", m.AsURL("/v1/accounts/update_credentials"), body)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername1))
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectAccount: &object.Account{
				Username: handler_test_setup.ExistingUsername1,
			},
		},
		{
			name: "CreateDupricatedUsername",
			request: func(m *handler_test_setup.C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"username":"%s"}`, handler_test_setup.ExistingUsername1)))
				req, err := http.NewRequest("POST", m.AsURL("/v1/accounts"), body)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusConflict,
		},
		{
			name: "CreateEmptyUsername",
			request: func(m *handler_test_setup.C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`{"username":"%s"}`, "")))
				req, err := http.NewRequest("POST", m.AsURL("/v1/accounts"), body)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusBadRequest,
		},
		{
			name: "CreateFailUnmatshalJSON",
			request: func(m *handler_test_setup.C) (*http.Response, error) {
				body := bytes.NewReader([]byte(fmt.Sprintf(`"usernam":"%s"}`, "aaa")))
				req, err := http.NewRequest("POST", m.AsURL("/v1/accounts"), body)
				if err != nil {
					t.Fatal(err)
				}
				return m.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusBadRequest,
		},
		{
			name: "FetchNotExist",
			request: func(m *handler_test_setup.C) (*http.Response, error) {
				req, err := http.NewRequest("GET", m.AsURL(fmt.Sprintf("/v1/accounts/%s", handler_test_setup.NotExistingUser)), nil)
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
			defer resp.Body.Close()
			if !assert.Equal(t, tt.expectStatusCode, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				actual := new(object.Account)
				if err = json.Unmarshal(body, actual); err != nil {
					t.Fatal(err)
				}

				opt := cmpopts.IgnoreFields(object.Account{}, "DisplayName", "Note", "Avatar")
				if d := cmp.Diff(actual, tt.expectAccount, opt); len(d) != 0 {
					t.Fatalf("differs: (-got +want)\n%s", d)
				}
			}
		})
	}
	os.RemoveAll("attachments")
}

func TestFollowReturnRelation(t *testing.T) {
	m := handler_test_setup.MockSetup()
	defer m.Close()

	tests := []struct {
		name               string
		request            func(c *handler_test_setup.C) (*http.Response, error)
		expectStatusCode   int
		expectRelationWith *object.RelationShip
	}{
		{
			name: "UnauthorizeFollow",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/follow", handler_test_setup.ExistingUsername1)
				req, err := http.NewRequest("POST", c.AsURL(url), nil)
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
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/follow", handler_test_setup.CreateUser)
				req, err := http.NewRequest("POST", c.AsURL(url), nil)
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
			name: "Follow",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/follow", handler_test_setup.ExistingUsername2)
				req, err := http.NewRequest("POST", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername1))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectRelationWith: &object.RelationShip{
				ID:         handler_test_setup.ID2,
				Following:  true,
				FollowedBy: false,
			},
		},
		{
			name: "Unfolollow",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/unfollow", handler_test_setup.ExistingUsername1)
				req, err := http.NewRequest("POST", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername2))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectRelationWith: &object.RelationShip{
				ID:         handler_test_setup.ID1,
				Following:  false,
				FollowedBy: true,
			},
		},
		{
			name: "UnfolollowNotExistingUser",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/unfollow", handler_test_setup.NotExistingUser)
				req, err := http.NewRequest("POST", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername2))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusNotFound,
		},
		{
			name: "Relationships",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := "/v1/accounts/relationships"
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("username", handler_test_setup.ExistingUsername2)
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername1))
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectRelationWith: &object.RelationShip{
				ID:         handler_test_setup.ID2,
				Following:  true,
				FollowedBy: false,
			},
		},
		{
			name: "UnauthorizeRelationships",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := "/v1/accounts/relationships"
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusUnauthorized,
		},
		{
			name: "NotExistRelationships",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := "/v1/accounts/relationships"
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("username", handler_test_setup.NotExistingUser)
				req.URL.RawQuery = params.Encode()
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authentication", fmt.Sprintf("username %s", handler_test_setup.ExistingUsername1))
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
				j := new(object.RelationShip)
				if assert.NoError(t, json.Unmarshal(body, j)) {
					if !reflect.DeepEqual(j, tt.expectRelationWith) {
						t.Fatalf("mismatch RelationShip:\n\t expect:\t%v\n\t actual:\t%v", tt.expectRelationWith, j)
					}
				}
			}
		})
	}
}

func TestFollowReturnAccounts(t *testing.T) {
	m := handler_test_setup.MockSetup()
	defer m.Close()

	tests := []struct {
		name             string
		request          func(c *handler_test_setup.C) (*http.Response, error)
		expectStatusCode int
		expectAccounts   []object.Account
	}{
		{
			name: "Following",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/following", handler_test_setup.ExistingUsername1)
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectAccounts:   []object.Account{{Username: handler_test_setup.ExistingUsername2}},
		},
		{
			name: "EmptyFollowing",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/following", handler_test_setup.ExistingUsername2)
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectAccounts:   []object.Account{},
		},
		{
			name: "FollowingNotExistAccount",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/following", handler_test_setup.NotExistingUser)
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusNotFound,
		},
		{
			name: "MoreThanMaxLimitFollowing",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/following", handler_test_setup.ExistingUsername1)
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
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
		},
		{
			name: "LessThanMinLimitFollowing",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/following", handler_test_setup.ExistingUsername1)
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
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
		},
		{
			name: "Followers",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/followers", handler_test_setup.ExistingUsername2)
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectAccounts:   []object.Account{{Username: handler_test_setup.ExistingUsername1}},
		},
		{
			name: "EmptyFollowers",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/followers", handler_test_setup.ExistingUsername1)
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusOK,
			expectAccounts:   []object.Account{{Username: handler_test_setup.ExistingUsername2}},
		},
		{
			name: "FollowersNotExistAccount",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/followers", handler_test_setup.NotExistingUser)
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				return c.Server.Client().Do(req)
			},
			expectStatusCode: http.StatusNotFound,
		},
		{
			name: "MoreThanMaxLimitFollowers",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/followers", handler_test_setup.ExistingUsername1)
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
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
		},
		{
			name: "LessThanMinLimitFollowers",
			request: func(c *handler_test_setup.C) (*http.Response, error) {
				url := fmt.Sprintf("/v1/accounts/%s/followers", handler_test_setup.ExistingUsername1)
				req, err := http.NewRequest("GET", c.AsURL(url), nil)
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
						t.Fatalf("mismatch Account:\n\t expect:\t%v\n\t actual:\t%v", tt.expectAccounts[0], j[0])
					}
				}
			}
		})
	}
}
