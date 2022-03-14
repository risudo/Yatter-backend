package dao

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	follows struct {
		db *sqlx.DB
	}
)

func NewFollows(db *sqlx.DB) repository.Follows {
	return &follows{db: db}
}

func (r *follows) Follow(ctx context.Context, id object.AccountID, follow *object.Account) error {
	return nil
}
