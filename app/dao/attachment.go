package dao

import (
	"context"
	"database/sql"
	"errors"
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

func (r *attachment) FindByStatusID(ctx context.Context, id object.StatusID) ([]object.Attachment, error) {
	var attachments []object.Attachment
	const query = `
	SELECT
		id,
		type,
		url,
		description
	FROM
		attachment A
	INNER JOIN status_contain_attachment S
	ON S.attachment_id = A.id
	WHERE status_id = ?`

	err := r.db.SelectContext(ctx, &attachments, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return attachments, nil
}
