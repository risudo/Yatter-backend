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

type (
	// Implementation for repository.Status
	status struct {
		db *sqlx.DB
	}
)

// Create status repository
func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

// statusを投稿
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

// idからstatusを取得
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

// account idからaccountを取得
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

// idで指定したstatusを削除
func (r *status) Delete(ctx context.Context, id object.StatusID) error {
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

// timelineを取得
func (r *status) PublicTimeline(ctx context.Context) (object.Timelines, error) {
	var timeline object.Timelines
	const query = `
	SELECT
		s.id AS 'id',
		s.account_id AS 'account_id',
		a.username AS 'account.username',
		s.create_at AS 'create_at',
		s.content AS 'content',
		a.create_at AS 'account.create_at'
	FROM status AS s JOIN account AS a ON s.account_id = a.id;`

	err := r.db.SelectContext(ctx, &timeline, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}

	return timeline, nil
}
