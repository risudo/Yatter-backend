package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Relation interface {
	Follow(ctx context.Context, followeeID object.AccountID, followerID object.AccountID) error
}
