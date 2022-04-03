package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Attachment interface {
	Insert(ctx context.Context, a object.Attachment) error
}
