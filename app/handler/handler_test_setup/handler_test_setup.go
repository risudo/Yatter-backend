package handler_test_setup

import (
	"context"
	"net/http/httptest"
	"net/url"
	"path"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/handler"
)

type (
	C struct {
		App    *app.App
		Server *httptest.Server
	}
	mockdao struct {
		accounts map[string]*object.Account
	}
)

const CreateUser = "smith"
const NotExistingUser = "fred"
const Content = "hello world"

const ID1 = 1
const ExistingUsername1 = "john"
const ID2 = 2
const ExistingUsername2 = "sum"

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

func (m *mockdao) InsertA(ctx context.Context, a object.Account) error {
	m.accounts[a.Username] = &object.Account{
		Username: a.Username,
	}
	return nil
}

func (m *mockdao) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	if account, ok := m.accounts[username]; ok {
		return account, nil
	}
	return nil, nil
}

func (m *mockdao) InsertS(ctx context.Context, status *object.Status) (object.StatusID, error) {
	return 1, nil
}

func (m *mockdao) FindByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
	if id == 1 {
		return &object.Status{
			Content: Content,
			Account: m.accounts[ExistingUsername1],
		}, nil
	}
	return nil, nil
}

func (m *mockdao) Delete(ctx context.Context, id object.StatusID) error {
	return nil
}

func (m *mockdao) PublicTimeline(ctx context.Context, p *object.Parameters) (object.Timelines, error) {
	return object.Timelines{
		object.Status{Content: Content},
	}, nil
}

func (m *mockdao) HomeTimeline(ctx context.Context, loginID object.AccountID, p *object.Parameters) (object.Timelines, error) {
	return object.Timelines{
		object.Status{Content: Content},
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

func (m *mockdao) Following(ctx context.Context, id object.AccountID, p object.Parameters) ([]object.Account, error) {
	if id == ID1 {
		return []object.Account{*m.accounts[ExistingUsername2]}, nil
	}
	return nil, nil
}

func (m *mockdao) Followers(ctx context.Context, id object.AccountID, p object.Parameters) ([]object.Account, error) {
	if id == ID2 {
		return []object.Account{*m.accounts[ExistingUsername1]}, nil
	}
	return nil, nil
}

func (m *mockdao) Unfollow(ctx context.Context, loginID object.AccountID, targetID object.AccountID) error {
	return nil
}

func MockSetup() *C {
	a1 := &object.Account{
		ID:       1,
		Username: ExistingUsername1,
	}
	a2 := &object.Account{
		ID:       2,
		Username: ExistingUsername2,
	}

	app := &app.App{Dao: &mockdao{accounts: map[string]*object.Account{
		a1.Username: a1,
		a2.Username: a2,
	}}}
	server := httptest.NewServer(handler.NewRouter(app))

	return &C{
		App:    app,
		Server: server,
	}
}

func (c *C) Close() {
	c.Server.Close()
}

func (c *C) AsURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
