package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Attachment interface {
	// Create Attachment
	Insert(ctx context.Context, a object.Attachment) (object.AttachmentID, error)

	// Fetch attachment which has specified statusID
	FindByStatusID(ctx context.Context, id object.StatusID) ([]object.Attachment, error)

	// Check if the attachment IDs exist
	HasAttachmentIDs(ctx context.Context, id []object.AttachmentID) (bool, error)
}
