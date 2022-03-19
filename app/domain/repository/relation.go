package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Relation interface {
	Follow(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error
	IsFollowing(ctx context.Context, followingID object.AccountID, followerID object.AccountID) (bool, error)
	Following(ctx context.Context, followingID object.AccountID) ([]object.Account, error)
	Followers(ctx context.Context, followerID object.AccountID) ([]object.Account, error)
	Unfollow(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error
}
