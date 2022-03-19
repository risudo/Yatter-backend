package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Relation interface {
	// Follow account
	Follow(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error

	// check if folloingID follows followeID
	IsFollowing(ctx context.Context, followingID object.AccountID, followerID object.AccountID) (bool, error)

	// Fetch accounts which the followingID follows
	Following(ctx context.Context, followingID object.AccountID) ([]object.Account, error)

	// Fetch accounts which follow followerID
	Followers(ctx context.Context, followerID object.AccountID) ([]object.Account, error)

	// followingID unfollow followerID
	Unfollow(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error
}
