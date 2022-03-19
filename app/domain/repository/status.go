package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Post status
	Post(ctx context.Context, status *object.Status) error

	// Fetch status which has specified id
	FindById(ctx context.Context, id object.StatusID) (*object.Status, error)

	// Delete status
	Delete(ctx context.Context, id object.AccountID) error

	// Fetch Timelines
	PublicTimeline(ctx context.Context) (object.Timelines, error)

	HomeTimeline(ctx context.Context, loginID object.AccountID) (object.Timelines, error)
}
