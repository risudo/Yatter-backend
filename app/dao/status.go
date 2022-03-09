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
	const query = "INSERT INTO status (content, account_id) VALUES(?, ?)"

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
	const query = "SELECT * FROM status WHERE id = ?"

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
	const query = "SELECT * FROM account WHERE id = ?"

	err := r.db.QueryRowxContext(ctx, query, id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

func (r *status) Delete(ctx context.Context, id object.AccountID) error {
	const query = "DELETE FROM status WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (r *status) PublicTimeline(ctx context.Context) ([]object.Status, error) {
	var timeline []object.Status //TODO: 型作った方がよさそう
	var status object.Status
	const query = "SELECT * FROM status"

	rows, err := r.db.QueryxContext(ctx, query)
	// TODO:これ必要？
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	//TODO: アカウント入れる

	for rows.Next() {
		err := rows.StructScan(&status)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		timeline = append(timeline, status)
	}
	return timeline, nil
} 
