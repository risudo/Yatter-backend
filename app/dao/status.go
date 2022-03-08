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
	err = r.db.QueryRowxContext(ctx, "SELECT * FROM status WHERE id = ?", id).StructScan(status)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (r *status) FindById(ctx context.Context, id object.StatusID) (*object.Status, error) {
	entity := new(object.Status)
	query := "SELECT * FROM status WHERE id = ?"
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	entity.Account, err = r.FindAccountById(ctx, entity.AccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

func (r *status) FindAccountById(ctx context.Context, id object.AccountID) (*object.Account, error) {
	entity := new(object.Account)
	query := "SELECT * FROM account WHERE id = ?"
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	return entity, nil
}
