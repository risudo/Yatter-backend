package dao

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type attachment struct {
	db *sqlx.DB
}

func NewAttachment(db *sqlx.DB) repository.Attachment {
	return &attachment{db: db}
}

func (r *attachment) Insert(ctx context.Context, a object.Attachment) error {
	return nil
}
