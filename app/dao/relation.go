package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	relation struct {
		db *sqlx.DB
	}
	// TODO: 変数名きもいかも
	exist struct {
		Exist bool `db:"existing"`
	}
)

func NewRelation(db *sqlx.DB) repository.Relation {
	return &relation{db: db}
}

func (r *relation) Follow(ctx context.Context, followingID object.AccountID, followerID object.AccountID) error {
	const query = "INSERT INTO relation (following_id, follower_id) VALUES(?, ?)"

	_, err := r.db.ExecContext(ctx, query, followingID, followerID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (r *relation) IsFollowing(ctx context.Context, followingID object.AccountID, followerID object.AccountID) (bool, error) {
	const query = "SELECT EXISTS(SELECT * FROM relation WHERE following_id = ? AND follower_id = ?) AS existing"
	exists := new(exist)

	log.Println("followingID", followingID, "followerID", followerID)
	err := r.db.QueryRowxContext(ctx, query, followingID, followerID).StructScan(exists)
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}
	return exists.Exist, nil
}

func (r *relation) Following(ctx context.Context, followingID object.AccountID) ([]object.Account, error) {
	var entity []object.Account
	const query = `
	SELECT
		account.id,
		account.username,
		account.create_at
	FROM
		account
	JOIN
		relation ON account.id = relation.follower_id
	WHERE
		relation.following_id = ?`

	err := r.db.SelectContext(ctx, &entity, query, followingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

func (r *relation) Followers(ctx context.Context, followerID object.AccountID) ([]object.Account, error) {
	var entity []object.Account

	const query = `
	SELECT
		account.id,
		account.username,
		account.create_at
	FROM
		account
	JOIN
		relation ON account.id = relation.following_id
	WHERE
		relation.follower_id = ?`
	err := r.db.SelectContext(ctx, &entity, query, followerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}
