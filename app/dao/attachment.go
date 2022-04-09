package dao

import (
	"context"
	"fmt"
	"log"
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

func (r *attachment) FindByIDs(ctx context.Context, ids []object.AttachmentID) ([]object.Attachment, error) {
	var attachments []object.Attachment
	query, args, err := sqlx.In("SELECT * FROM attachment WHERE id IN(?)", ids)
	if err != nil {
		return nil, err
	}
	log.Println("query: ", query)
	log.Println("args:", args)
	err = r.db.SelectContext(ctx, &attachments, r.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}
	return attachments, nil
}
