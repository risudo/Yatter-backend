package dao

import (
	"context"
	"fmt"
	"log"
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
	row, err := r.db.ExecContext(ctx, query, status.Content, status.Account.ID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	id, err := row.LastInsertId()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	log.Println("before:", *status)
	err = r.db.QueryRowxContext(ctx, "SELECT * FROM status WHERE id = ?", id).StructScan(status)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("%w", err)
	}
	log.Println("after:", *status)
	return nil
}
