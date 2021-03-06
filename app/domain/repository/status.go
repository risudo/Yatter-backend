package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Post status
	Insert(ctx context.Context, status object.Status, mediaIDs []object.AttachmentID) (object.StatusID, error)

	// Fetch status which has specified id
	FindByID(ctx context.Context, id object.StatusID) (*object.Status, error)

	// Delete status
	Delete(ctx context.Context, id object.StatusID) error

	// Fetch Public Timelines
	PublicTimeline(ctx context.Context, p object.Parameters) (object.Timelines, error)

	// Fetch Home Timelines
	HomeTimeline(ctx context.Context, loginID object.AccountID, p object.Parameters) (object.Timelines, error)
}
