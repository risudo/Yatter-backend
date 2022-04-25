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
	// Implementation for repository.Account
	account struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

// FindByUsername : ユーザ名からユーザを取得
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	const query =`
	SELECT
		id,
		username,
		display_name,
		avatar,
		header,
		note,
		create_at,
		CASE
		WHEN
			(SELECT COUNT(*) FROM relation WHERE following_id = (SELECT id from account WHERE username = ?) GROUP BY following_id) IS NULL
		THEN 0
		ELSE
			(SELECT COUNT(*) FROM relation WHERE following_id = (SELECT id from account WHERE username = ?) GROUP BY following_id)
		END AS followingcount,
		CASE
		WHEN
			(SELECT COUNT(*) FROM relation WHERE follower_id = (SELECT id from account WHERE username = ?) GROUP BY follower_id) IS NULL
		THEN 0
		ELSE (SELECT COUNT(*) FROM relation WHERE follower_id = (SELECT id from account WHERE username = ?) GROUP BY follower_id)
		END AS followerscount
	FROM account
	WHERE username = ?
	`

	err := r.db.QueryRowxContext(ctx, query, username, username, username, username, username).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}

	return entity, nil
}

// アカウントを作成
func (r *account) Insert(ctx context.Context, a object.Account) (object.AccountID, error) {
	const query = "INSERT INTO account (username, password_hash) VALUES (?, ?)"

	row, err := r.db.ExecContext(ctx, query, a.Username, a.PasswordHash)
	if err != nil {
		return -1, fmt.Errorf("%w", err)
	}
	id, err := row.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *account) Update(ctx context.Context, a object.Account) error {
	const query = `
	UPDATE account
	SET
		display_name = ?,
		note = ?,
		avatar = ?,
		header = ?
	WHERE username = ?
	`
	_, err := r.db.ExecContext(ctx, query, a.DisplayName, a.Note, a.Avatar, a.Header, a.Username)
	if err != nil {
		return err
	}
	return nil
}
