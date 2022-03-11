package dao

import (
	"context"
	"testing"
	"yatter-backend-go/app/domain/object"
)

type accountRepositoryMock struct {
	findByUsernamefn func(ctx context.Context, username string) (*object.Account, error)
	createfn         func(ctx context.Context, entity *object.Account) error
}

func (s *accountRepositoryMock) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	return s.findByUsernamefn(ctx, username)
}

func (s *accountRepositoryMock) Create(ctx context.Context, entity *object.Account) error {
	return s.createfn(ctx, entity)
}

func TestAccountCreate(t *testing.T) {
	mock := &accountRepositoryMock{
		findByUsernamefn: func(ctx context.Context, username string) (*object.Account, error) {
			return nil, nil
		},
		createfn: func(ctx context.Context, entity *object.Account) error {
			return nil
		},
	}

	account := NewAccount(nil)
	if err := account.Create(context.Background(), &object.Account{
		Username: "testuser",
		PasswordHash: "testpass",
	}); err != nil {
		t.Fatal(err)
	}
}
