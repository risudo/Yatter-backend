package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Relation interface {
	// Follow the account of followerID
	Follow(ctx context.Context, loginID object.AccountID, targetID object.AccountID) error

	// check if folloingID follows followeID
	IsFollowing(ctx context.Context, accountID object.AccountID, targetID object.AccountID) (bool, error)

	// Fetch accounts which the account of id follows
	Following(ctx context.Context, id object.AccountID, p object.Parameters) ([]object.Account, error)

	// Fetch accounts which follow the account of id
	Followers(ctx context.Context, id object.AccountID, p object.Parameters) ([]object.Account, error)

	// unfollow the account of followerID
	Unfollow(ctx context.Context, loginID object.AccountID, targetID object.AccountID) error
}
