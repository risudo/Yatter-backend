package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Follows interface {
	Follow(ctx context.Context, account *object.Account, follow *object.Account) error
}
