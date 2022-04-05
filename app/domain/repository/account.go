package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Account interface {
	// Fetch account which has specified username
	FindByUsername(ctx context.Context, username string) (*object.Account, error)

	// Create account
	Insert(ctx context.Context, account object.Account) error

	// Update account
	Update(ctx context.Context, account object.Account) error
}
