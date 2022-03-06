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
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	return entity, nil
}

func (r *account) CreateAccount(ctx context.Context, entity *object.Account) error {
	query := "insert into account (username, password_hash) values (?, ?)"
	_, err := r.db.ExecContext(ctx, query, entity.Username, entity.PasswordHash)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// 2回クエリ投げないといけないのが微妙かも
	err = r.db.QueryRowxContext(ctx, "SELECT * from account where username = ?", entity.Username).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("%w", err)
	}
	return nil
}
