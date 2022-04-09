package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Attachment interface {
	Insert(ctx context.Context, a object.Attachment) (object.AttachmentID, error)
	FindByIDs(ctx context.Context, ids []object.AttachmentID) ([]object.Attachment, error)
}
