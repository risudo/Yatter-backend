package dao

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

// Implementation for repository.Status
type status struct {
	db *sqlx.DB
}

func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

func (r *status) Post(ctx context.Context, status *object.Status) error {
	query := "INSERT INTO status (content, account_id) VALUES(?, ?)"
	_, err := r.db.DB.ExecContext(ctx, query, status.Content, status.Account.ID)
	if err != nil {
		return err
	}
	return nil
}
