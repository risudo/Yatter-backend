package dao

import (
	"context"
	"fmt"
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

func (r *attachment) Insert(ctx context.Context, a object.Attachment) (object.AttachmentID, error) {
	const query = `INSERT INTO attachment (type, url, description) VALUES(?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, a.MediaType, a.URL, a.Description)
	if err != nil {
		return -1, fmt.Errorf("%w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%w", err)
	}
	return id, nil
}
