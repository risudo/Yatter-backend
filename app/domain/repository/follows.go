package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Follows interface {
	Follow(ctx context.Context, id object.AccountID, follow *object.Account) error
}
